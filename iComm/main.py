import logging
import grpc
import main_pb2
import main_pb2_grpc

HEADER_GUN = 70
HEADER_VEST = 86
HEADER_GLOVE = 77

def pass_params(header, data_obj):
    with grpc.insecure_channel('localhost:8081') as channel:
        stub = main_pb2_grpc.RelayStub(channel)

        # grenade, reload, shield
        if header == HEADER_GLOVE:
            msg = main_pb2.SensorData(
                player=1,
                index=data_obj["index"],
                roll=data_obj["roll"],
                pitch=data_obj["pitch"],
                yaw=data_obj["yaw"],
                x=data_obj["x"],
                y=data_obj["y"],
                z=data_obj["z"],
            )
            resp = stub.Gesture(msg)
        # shoot/shot
        # TODO Add Packet IDs for Vest/Gun
        elif header == HEADER_GUN:
            msg = main_pb2.Event(player=1, shootID=data_obj, action=main_pb2.shoot)
            resp = stub.Shoot(msg)
        elif header == HEADER_VEST:
            msg = main_pb2.Event(player=1, shootID=data_obj, action=main_pb2.shot)
            resp = stub.Shot(msg)

        print(f"Received resp {resp}")

if __name__ == "__main__":
    logging.basicConfig()
    pass_params()

