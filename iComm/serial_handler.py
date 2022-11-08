from threading import Thread
from constants import PACKET_SIZE, RETRY_COUNT
from helper import unpack_glove_data_into_dict
from bluepy.btle import BTLEDisconnectError
import time
import main_pb2
import os

action_classes = ["final", "grenade", "reload", "shield", "idle"]
data_per_class = 25


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

        # get user number
        self.user = input("Enter user number: ")
        # create directory for user
        current_directory = os.getcwd()
        self.data_directory = os.path.join(current_directory, f'data{self.user}')
        if not os.path.exists(self.data_directory):
            os.makedirs(self.data_directory)

        # initialize a sequence_of_action with data_per_class for each action class
        self.sequence_of_action = []     
        for action_class_index in range(len(action_classes)):
            self.sequence_of_action.extend([action_class_index] * data_per_class)
        self.action_class_count = [0, 0, 0, 0, 0]

        # shuffle the sequence_of_action
        import random
        random.shuffle(self.sequence_of_action)  
        self.current_action = self.sequence_of_action.pop(0)
        print(f"Current action: {action_classes[self.current_action]}")

        self.data_str = ''

    def pass_params(self, packet):
        glove_data = packet[1:15]
        data_obj = unpack_glove_data_into_dict(glove_data)
        # print(data_obj)
        if data_obj["index"] == 0:

            if len(self.data_str) != 0:
                with open(f'data{self.user}/{action_classes[self.current_action]}{self.action_class_count[self.current_action]}.csv', 'w') as f:
                    # print(f'wrote to {action_classes[self.current_action]}{self.action_class_count[self.current_action]}')
                    f.write(self.data_str)
                    self.action_class_count[self.current_action] += 1
                    self.data_str = ''

                if len(self.sequence_of_action) == 0:
                    print("All data collected! Exiting...")
                    exit()
                self.current_action = self.sequence_of_action.pop(0)
            if len(self.sequence_of_action) != 0:
                print(f"Current action: {action_classes[self.sequence_of_action[0]]}")
            else:
                print("ended")

        self.data_str += f'{data_obj["roll"]}, {data_obj["pitch"]}, {data_obj["yaw"]}, {data_obj["x"]}, {data_obj["y"]}, {data_obj["z"]}\n'

    
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


