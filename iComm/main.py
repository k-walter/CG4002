import logging

import random
import grpc
import main_pb2
import main_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:8081') as channel:
        while True:
            input("Enter to send dummy packet...")
            stub = main_pb2_grpc.RelayStub(channel)
            msg = main_pb2.SensorData(player=1)
            if random.randint(1, 4) < 4:
                resp = stub.Gesture(msg) # grenade, reload, shield
            else:
                resp = stub.Shoot(msg)
            print(f"Received resp {resp}")

if __name__ == "__main__":
    logging.basicConfig()
    run()
