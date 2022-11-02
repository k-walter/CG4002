from concurrent import futures
import threading
import logging

from google.protobuf import empty_pb2
import grpc
import main_pb2
import main_pb2_grpc

import driver
import random
import numpy as np


predToAction = [ # [0, 4]
    main_pb2.shield,
    main_pb2.reload,
    main_pb2.grenade,
    main_pb2.logout,

    main_pb2.none, # -1
]

class Pynq(main_pb2_grpc.PynqServicer):
    def __init__(self):
        self.myip = driver.Model("bitstream/cnn.bit")
        self.myip.setCNNWeights(np.load("CNN_weights.npy"))
        self.myip.setCNNBias(np.load("CNN_bias.npy"))
        self.myip.setDenseWeights(np.load("dense_weights.npy"))
        self.myip.setDenseBias(np.load("dense_bias.npy"))

        self.myip.debug = False
        self.debounce = [False, False]
        self.mu = threading.Lock()

    def Emit(self, req: main_pb2.Data, context):
        self.mu.acquire()
        p = req.player - 1
        # Reset actions
        if req.index == 0:
            self.debounce[p] = False

        # Debounce inference until next glove action
        if self.debounce[p]:
            self.mu.release()
            return main_pb2.Event(
                player=req.player,
                time=req.time,
                action=predToAction[-1]
            )

        # Blocking inference
        axn = self.myip.inference(data = [
            req.roll,
            req.pitch,
            req.yaw,
            req.x,
            req.y,
            req.z
        ], user_number = p)

        # Inferred something?
        if axn != -1:
            self.myip.resetBuffer(p)
            self.debounce[p] = True

        # Resp
        self.mu.release()
        return main_pb2.Event(
            player=req.player,
            time=req.time,
            rnd=req.rnd,
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

    print("running")

    # Blocking
    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    run()
