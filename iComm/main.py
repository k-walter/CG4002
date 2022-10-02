import logging
import grpc
import main_pb2
import main_pb2_grpc

HEADER_GUN = 71
HEADER_VEST = 86

def run(*args):
    with grpc.insecure_channel('localhost:8081') as channel:

        stub = main_pb2_grpc.RelayStub(channel)

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
            resp = stub.Gesture(msg)

        # shoot/shot
        else:
            if args[0] == HEADER_GUN:
                msg = main_pb2.Event(player=1, action=main_pb2.shoot)
                resp = stub.Shoot(msg)
            if args[0] == HEADER_VEST:
                msg = main_pb2.Event(player=1, action=main_pb2.shot)
                resp = stub.Shot(msg)

        print(f"Received resp {resp}")

if __name__ == "__main__":
    logging.basicConfig()
    run()

