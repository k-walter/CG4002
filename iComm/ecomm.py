#!/usr/bin/env python3
from threading import Thread
import asyncio
import grpc
import main_pb2
import main_pb2_grpc
from google.protobuf.empty_pb2 import Empty


class EComm(Thread):
    def __init__(self, port: str):
        # External: Setup thread, coroutines
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)  # for grpc to use event loop
        Thread.__init__(self, target=self._loop.run_forever)

        # Internal: Init connection
        self._channel = grpc.aio.insecure_channel(port)  # asyncio interfaces
        self._stub = main_pb2_grpc.RelayStub(self._channel)

        # Internal: Init streams
        self._gestureStub = self._stub.Gesture()
        self._shootStub = self._stub.Shoot()
        self._shotStub = self._stub.Shot()
        self._loop.create_task(self._getRound())

        # Internal: Bookeeping
        self._rnd = 1

        print("init ecomm")

    """ THREAD SAFE """

    def gesture(self, data: main_pb2.Data) -> None:
        async def f(data: main_pb2.Data):
            data.rnd = self._rnd
            await self._gestureStub.write(data)
        asyncio.run_coroutine_threadsafe(f(data), self._loop)

    def shoot(self, data: main_pb2.Event) -> None:
        async def f(data: main_pb2.Event):
            data.rnd = self._rnd
            await self._shootStub.write(data)
            print(f"sent shoot ID={data.shootID}")
        asyncio.run_coroutine_threadsafe(f(data), self._loop)

    def shot(self, data: main_pb2.Event) -> None:
        async def f(data: main_pb2.Event):
            data.rnd = self._rnd
            await self._shotStub.write(data)
            print(f"sent shot ID={data.shootID}")
        asyncio.run_coroutine_threadsafe(f(data), self._loop)

    """ ASYNCIO SAFE """

    async def _getRound(self):
        async for resp in self._stub.GetRound(Empty()):
            self._rnd = resp.rnd  # assume monotonic incr


def benchmark():
    ec = EComm("localhost:8081")
    ec.start()
    N = 100

    for i in range(1,N+1):
        ec.shoot(main_pb2.Event(player=1, shootID=i, action=main_pb2.shoot))
    for i in range(1,N+1):
        ec.gesture(main_pb2.Data(player=2, index=i,
                   roll=1, pitch=1, yaw=1, x=1, y=1, z=1))


if __name__ == "__main__":
    benchmark()
