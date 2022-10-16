from constants import HEADER_GLOVE, HEADER_GUN, HEADER_VEST
from serial_handler import GloveHandler, GunHandler, VestHandler

class SerialHandlerFactory():
    def get_serial_handler(beetle, lock, stub):
        if beetle.header == HEADER_GLOVE:
            return GloveHandler(beetle, lock, stub)
        elif beetle.header == HEADER_GUN:
            return GunHandler(beetle, lock, stub)
        elif beetle.header == HEADER_VEST:
            return VestHandler(beetle, lock, stub)


