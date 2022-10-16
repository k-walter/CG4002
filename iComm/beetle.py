from bluepy import btle
from delegate import Delegate

SERIAL_UUID = "0000dfb1-0000-1000-8000-00805f9b34fb"

class Beetle():
    def __init__(self, address, name, header):
        self.peripheral = btle.Peripheral()
        self.address = address
        self.name = name
        self.header = header
        self.is_connected = False
        self.has_handshake = False
        self.serial_char = None
        self.delegate = None

    # only call this when we are sure we are not connected
    def connect_with_retries(self, retries):
        self.is_connected = False

        while not self.is_connected and retries > 0:
            try:
                print(f"{retries} Attempts Left:")
                self.__connect()
            except btle.BTLEException as e:
                print(e)
            retries -= 1
        if not self.is_connected:
            print(f"{self.name} could not connect!")
            print("Exiting Program...")
            exit()

    # only call when we encounter BTLEDisconnectError
    def set_disconnected(self):
        print(f"{self.name} disconnected. Attempting Reconnection...")
        self.is_connected = False

    def init_peripheral(self):
        self.__set_serial_char()
        self.__set_delegate()
        self.__attach_delegate()

    def init_handshake(self):
        try:
            self.__try_init_handshake()

        # Catches disconnect exception raised by __try_init_handshake.
        # Callers of init_handshake will handle reconnection process
        except btle.BTLEDisconnectError:
            self.set_disconnected()


    #Private Methods
    def __connect(self):
        print(f"Connecting to {self.address}...")
        self.peripheral.connect(self.address)
        self.is_connected = True
        print("Connected.")

    def __try_init_handshake(self):
        self.has_handshake = False
        while not self.has_handshake:
            self.__send_handshake()
            if self.peripheral.waitForNotifications(5.0):
                self.delegate.handle_data()
                self.__receive_ack()

    # Sets serial characteristic in order to write to beetle
    def __set_serial_char(self):
        print(f"Setting serial characteristic for {self.name}...")
        chars = self.peripheral.getCharacteristics()
        serial_char = [c for c in chars if c.uuid == SERIAL_UUID][0]
        self.serial_char = serial_char
        print("Serial characteristic set.")

    # Creates delegate object to receive notifications
    def __set_delegate(self):
        print(f"Setting Delegate for {self.name}...")
        self.delegate = Delegate(self.serial_char, self.header)
        print("Delegate set.")

    # Attaches delegate object to peripheral
    def __attach_delegate(self):
        print(f"Attaching {self.name} delegate to peripheral...")
        self.peripheral.withDelegate(self.delegate)
        print("Done.")

    # Sends Handshake Packet
    def __send_handshake(self):
        print("Handshake in Progress...")
        self.serial_char.write(b'H')

    def __receive_ack(self):
        if self.delegate.hand_ack:
            print(f"Handshake ACK received from {self.name}")
            self.has_handshake = True
            self.serial_char.write(b'A')