from threading import Thread
from constants import PACKET_SIZE, RETRY_COUNT
from helper import bytes_to_uint16_t, unpack_glove_data_into_dict
from bluepy.btle import BTLEDisconnectError
import time
import main_pb2
import logging

class SerialHandler(Thread):
    def __init__(self, beetle, lock, stub):
        Thread.__init__(self)
        self.beetle = beetle
        self.lock = lock
        self.stub = stub
        self.player_no = beetle.player_no

    def run(self):
        while True:
            try:
                self.beetle.peripheral.waitForNotifications(0.01)
                while len(self.beetle.delegate.data_buffer) >= PACKET_SIZE:
                    self.beetle.delegate.handle_data()
                    if not self.beetle.delegate.is_valid_data:
                        continue
                    with self.lock:
                        self.pass_params(self.beetle.delegate.packet)
                    
            except BTLEDisconnectError:
                self.beetle.set_disconnected()
                while not self.beetle.is_connected:
                    self.beetle.connect_with_retries(RETRY_COUNT)
                    print(f"{self.beetle.name} reconnected. \
                        Reinitialising handshake...")
                    self.beetle.init_handshake()

    def pass_params(data_obj):
        pass
     

class GloveHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)
        self.next_send = 0
        self.rnd = 1
        self.send_buf = main_pb2.SensorData()

    def pass_params(self, packet):
        glove_data = packet[3:-1]
        data_obj = unpack_glove_data_into_dict(glove_data)
        index = bytes_to_uint16_t(packet[1:3])
        now = time.monotonic_ns()

        data = self.send_buf.data.add()
        data.player = self.player_no
        data.rnd = self.rnd
        data.index = index
        data.roll = data_obj["roll"]
        data.pitch = data_obj["pitch"]
        data.yaw = data_obj["yaw"]
        data.x = data_obj["x"]
        data.y = data_obj["y"]
        data.z = data_obj["z"]

        if now < self.next_send:
            return
        
        self.next_send = now + int(20e6)
        self.rnd = self.stub.Gesture(self.send_buf)
        self.send_buf = main_pb2.SensorData()

    
class VestHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)

    def pass_params(self, packet):
        shoot_id = packet[2]
        msg = main_pb2.Event(player=self.player_no, shootID=shoot_id, action=main_pb2.shot)
        self.stub.Shot(msg)


class GunHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)

    def pass_params(self, packet):
        shoot_id = packet[2]
        msg = main_pb2.Event(player=self.player_no, shootID=shoot_id, action=main_pb2.shoot)
        self.stub.Shoot(msg)


