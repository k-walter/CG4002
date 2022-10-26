from constants import HEADER_GLOVE, HEADER_GUN, HEADER_VEST
from serial_handler import GloveHandler, GunHandler, VestHandler

class SerialHandlerFactory():
    def get_serial_handler(beetle, ec: 'EComm'):
        if beetle.header == HEADER_GLOVE:
            return GloveHandler(beetle, ec.gesture)
        elif beetle.header == HEADER_GUN:
            return GunHandler(beetle, ec.shoot)
        elif beetle.header == HEADER_VEST:
            return VestHandler(beetle, ec.shot)
