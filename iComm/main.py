import logging

import grpc
import main_pb2
import main_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:8081') as channel:
        while True:
            input("Enter to send dummy packet...")
            stub = main_pb2_grpc.RelayToUltraStub(channel)
            msg = main_pb2.google_dot_protobuf_dot_empty__pb2.Empty()
            resp = stub.EmitFoo(msg)
            print(f"Received resp {resp}")

if __name__ == "__main__":
    logging.basicConfig()
    run()
