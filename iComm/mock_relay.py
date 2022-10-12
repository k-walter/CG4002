import logging

import time
import random
import grpc
import main_pb2
import main_pb2_grpc


def run():
    with grpc.insecure_channel('localhost:8081') as channel:
        shootID = 1
        while True:
            axn = int(input("gesture(0), shoot(1)"))
            stub = main_pb2_grpc.RelayStub(channel)

            # grenade, reload, shield
            if axn == 0:
                msg = main_pb2.SensorData()
                data = msg.data.add()
                data.player = 1
                data.index = 1

                resp = stub.Gesture(msg)
                print("Sent gesture")

            # shoot
            else:
                def shoot():
                    msg = main_pb2.Event(player=1, shootID=shootID, action=main_pb2.shoot)
                    resp = stub.Shoot(msg)
                    print("sent shoot")
                def shot():
                    msg = main_pb2.Event(player=2, shootID=shootID, action=main_pb2.shot)
                    resp = stub.Shot(msg)
                    print("sent shot")

                # random order of shoot/shot
                axns = [0,1]
                random.shuffle(axns)
                for i in axns:
                    if i == 0:
                        shoot()
                    else:
                        shot()
                    time.sleep(.01)
                shootID += 1

if __name__ == "__main__":
    logging.basicConfig()
    run()
