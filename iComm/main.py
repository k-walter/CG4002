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

            # grenade, reload, shield
            if random.randint(1, 4) != 1:
                msg = main_pb2.SensorData(
                    player=1,
                    roll=random.randint(0, (1<<16) - 1),
                    pitch=random.randint(0, (1<<16) - 1),
                    roll=random.randint(0, (1<<16) - 1),
                    x=random.randint(0, (1<<16) - 1),
                    y=random.randint(0, (1<<16) - 1),
                    z=random.randint(0, (1<<16) - 1),
                )
                resp = stub.Gesture(msg)

            # shoot
            else:
                msg = main_pb2.Event(player=1, action=main_pb2.shoot)
                resp = stub.Shoot(msg)

            print(f"Received resp {resp}")

if __name__ == "__main__":
    logging.basicConfig()
    run()
