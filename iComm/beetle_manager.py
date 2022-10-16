from constants import address_list, name_list, header_list, RETRY_COUNT
from beetle import Beetle

class BeetleManager():
    def __init__(self):
        self.beetle_list = []

    def initialise_beetle_list(self):
        self.init_params_beetle_list()
        self.init_connect_beetle_list()
        self.init_peripheral_beetle_list()
        self.init_handshake_beetle_list()

    # Creates Beetle objects with their corresponding address and name,
    # and appends it to the beetle list.
    def init_params_beetle_list(self):
        for addr, name, header in zip(address_list, name_list, header_list):
            self.beetle_list.append(Beetle(addr, name, header))


    def init_connect_beetle_list(self):
        for beetle in self.beetle_list:
            beetle.connect_with_retries(RETRY_COUNT)


    def init_peripheral_beetle_list(self):
        for beetle in self.beetle_list:
            beetle.init_peripheral()


    def init_handshake_beetle_list(self):
        for beetle in self.beetle_list:
            while not beetle.has_handshake and beetle.is_connected:
                beetle.init_handshake()

                # If beetle disconnects during the handshake, will attempt
                # reconnection and handshaking
                if not beetle.is_connected:
                    beetle.connect_with_retries(RETRY_COUNT)