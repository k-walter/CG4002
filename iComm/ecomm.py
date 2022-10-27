#!/usr/bin/env python3

from threading import Thread
import asyncio
import grpc
import main_pb2
import main_pb2_grpc

import queue
from typing import Union, Iterable

class EComm(Thread):
    def __init__(self, port: str):
        # External: Event loop to read from icomm
        self._loop = asyncio.new_event_loop()
        Thread.__init__(self, target=self._loop.run_forever)

        # Internal: Setup connection, streams
        self._channel = grpc.insecure_channel(port)
        self._stub = main_pb2_grpc.RelayStub(self._channel)
        self._q = asyncio.Queue()
        self._loop.create_task(self.sendGesture())
        self._loop.create_task(self.getRound())

        # Internal: Bookeeping
        self._rnd = 1

        print("init ecomm")

    """ THREAD SAFE """

    def gesture(self, data: main_pb2.Data) -> None:
        async def f(data: main_pb2.Data):
            data.rnd = self._rnd
            await self._q.put(data)
        self._loop.create_task(f(data))

    def shoot(self, data: main_pb2.Event) -> None:
        async def f(data: main_pb2.Event):
            data.rnd = self._rnd
            await self._stub.Shoot(data) # OPTIMIZE stream
        self._loop.create_task(f(data))

    def shot(self, data: main_pb2.Event) -> None:
        async def f(data: main_pb2.Event):
            data.rnd = self._rnd
            await self._stub.Shot(data)
        self._loop.create_task(f(data))

    """ ASYNCIO SAFE """
    async def getRound(self):
        for resp in self._stub.GetRound():
            self._rnd = resp.rnd

    async def sendGesture(self):
        await self._stub.Gesture(iter(self._q.get, None))
