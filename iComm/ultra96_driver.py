from contextlib import contextmanager
import logging
import grpc
import main_pb2
import main_pb2_grpc

HEADER_GUN = 70
HEADER_VEST = 86

class Ultra96Driver():
    def __init__(self, port):
        self.port = port

    def __enter__(self):
        self.channel = grpc.insecure_channel(self.port)
        return self.channel

    def __exit__(self, excep_type, excep_val, traceback):
        self.channel.close()

    def pass_params(self, *args):
        self.stub = main_pb2_grpc.RelayStub(self.channel)
        # grenade, reload, shield
        if len(args) == 2:
            data_obj = args[1]
            msg = main_pb2.SensorData(
                player=1,
                roll=data_obj["roll"],
                pitch=data_obj["pitch"],
                yaw=data_obj["yaw"],
                x=data_obj["x"],
                y=data_obj["y"],
                z=data_obj["z"],
            )
            resp = self.stub.Gesture(msg)

        # shoot/shot
        else:
            if args[0] == HEADER_GUN:
                msg = main_pb2.Event(player=1, action=main_pb2.shoot)
                resp = self.stub.Shoot(msg)
            if args[0] == HEADER_VEST:
                msg = main_pb2.Event(player=1, action=main_pb2.shot)
                resp = self.stub.Shot(msg)

        print(f"Received resp {resp}")

