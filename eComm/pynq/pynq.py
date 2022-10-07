from concurrent import futures
import logging

import time
import random
from google.protobuf import empty_pb2
import grpc
import main_pb2
import main_pb2_grpc

class Pynq(main_pb2_grpc.PynqServicer):
    def __init__(self):
        self.action = main_pb2.none
        self.actions = (
            main_pb2.grenade,
            main_pb2.reload,
            main_pb2.shield
        )

    def Emit(self, request: main_pb2.SensorData, context):
        # TODO format & forward to fpga
        print(f"Received {request}")
        self.action = random.choice(self.actions)
        return empty_pb2.Empty()

    def Poll(self, request, context) -> main_pb2.Event:
        # TODO poll fpga and reset
        axn = self.action
        if self.action != main_pb2.none:
            print(f"Predicted {axn}")

        # reset detection
        self.action = main_pb2.none

        return main_pb2.Event(player=1, action=axn)

def run():
    with grpc.insecure_channel('localhost:8083') as channel:
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        main_pb2_grpc.add_PynqServicer_to_server(Pynq(), server)
        server.add_insecure_port('localhost:8082')
        server.start()
        server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    run()
