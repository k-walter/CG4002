from concurrent import futures
import logging

from google.protobuf import empty_pb2
import grpc
import main_pb2
import main_pb2_grpc

import random


predToAction = [ # [0, 4]
    main_pb2.shield,
    main_pb2.reload,
    main_pb2.grenade,
    main_pb2.logout,

    main_pb2.none, # -1
]

class Pynq(main_pb2_grpc.PynqServicer):
    def __init__(self):
        pass

    def Emit(self, req: main_pb2.Data, context):
        # Resp
        axn = int(input("shield(0), reload(1), grenade(2), logout(3)"))
        return main_pb2.Event(
            player=req.player,
            time=req.time,
            action=predToAction[axn]
        )

    # NOT IN USE
    def Poll(self, request, context) -> main_pb2.Event:
        return main_pb2.Event(player=1, action=main_pb2.none)

def run():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    main_pb2_grpc.add_PynqServicer_to_server(Pynq(), server)
    server.add_insecure_port('localhost:8082')
    server.start()

    # Blocking
    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    run()
