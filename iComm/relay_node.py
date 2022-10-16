from beetle_manager import BeetleManager
from serial_handler_factory import SerialHandlerFactory
from constants import MY_PORT

import threading
import time
import grpc
import main_pb2_grpc

if __name__ == "__main__":

    beetle_manager = BeetleManager()
    beetle_manager.initialise_beetle_list()

    # Delay to read the initial print statements
    time.sleep(2)

    channel = grpc.insecure_channel(MY_PORT)
    stub = main_pb2_grpc.RelayStub(channel)
    lock = threading.Lock()

    # Starts threads
    for beetle in beetle_manager.beetle_list:
        handler = SerialHandlerFactory.get_serial_handler(beetle, lock, stub)
        handler.start()

    while True:
        time.sleep(int(1e9))

