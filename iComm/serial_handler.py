from threading import Thread
from constants import PACKET_SIZE, RETRY_COUNT
from helper import unpack_glove_data_into_dict
from bluepy.btle import BTLEDisconnectError
import time
import main_pb2
import os

count = 0
data_str = ''

class SerialHandler(Thread):
    def __init__(self, beetle, lock, stub):
        Thread.__init__(self)
        self.beetle = beetle
        self.lock = lock
        self.stub = stub

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
                    print(f"{self.beetle.name} reconnected. Reinitialising handshake...")
                    self.beetle.init_handshake()

    def pass_params(data_obj):
        pass
     

class GloveHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)
        self.next_send = 0
        self.send_buf = main_pb2.SensorData()

    def pass_params(self, packet):
        global count
        global data_str
        glove_data = packet[1:14]
        data_obj = unpack_glove_data_into_dict(glove_data)
        # now = time.monotonic_ns()
        if data_obj["index"] == 0:
            if count == 50:
                print("\n\n\n\n\nCollected 50! Exiting...")
                exit()
            if len(data_str) != 0:
                with open(f'final{count}.csv', 'w') as f:
                    f.write(data_str)
                    count += 1
                    data_str = ''
                    print(f'wrote to final{count-1}')

            
        # data = self.send_buf.data.add()
        # data.player=1
        # data.index=data_obj["index"]
        # data.roll=data_obj["roll"]
        # data.pitch=data_obj["pitch"]
        # data.yaw=data_obj["yaw"]
        # data.x=data_obj["x"]
        # data.y=data_obj["y"]
        # data.z=data_obj["z"]

        data_str += f'{data_obj["roll"]}, {data_obj["pitch"]}, {data_obj["yaw"]}, {data_obj["x"]}, {data_obj["y"]}, {data_obj["z"]}\n'

        # if now < self.next_send:
        #     return
        
        # self.next_send = now + int(20e6)
        # self.stub.Gesture(self.send_buf)
        # self.send_buf = main_pb2.SensorData()

    
class VestHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)

    def pass_params(self, packet):
        shoot_id = packet[2]
        msg = main_pb2.Event(player=2, shootID=shoot_id, action=main_pb2.shot)
        self.stub.Shot(msg)


class GunHandler(SerialHandler):
    def __init__(self, beetle, lock, stub):
        super().__init__(beetle, lock, stub)

    def pass_params(self, packet):
        shoot_id = packet[2]
        msg = main_pb2.Event(player=1, shootID=shoot_id, action=main_pb2.shoot)
        self.stub.Shoot(msg)


