from bluepy import btle
from beetle import Beetle

import threading
import time
import grpc
import main_pb2
import main_pb2_grpc

ADDRESS_GUN_2 = "d0:39:72:bf:c6:51"

ADDRESS_GLOVE_1 = "d0:39:72:bf:c6:47"
ADDRESS_GLOVE_2 = "d0:39:72:bf:c8:87"

ADDRESS_VEST_1 = "c4:be:84:20:1b:5e"
ADDRESS_BACKUP = "50:F1:4A:DA:CC:EB"

PACKET_SIZE = 15

HEADER_ACK = 65
HEADER_GUN = 70
HEADER_GLOVE = 77
HEADER_VEST = 86

MAX_16_BIT_SIGNED = 32768
MAX_16_BIT_UNSIGNED = 65535

next_send = 0
send_buf = main_pb2.SensorData()

MY_PORT = 'localhost:8080'

# address_list = [ADDRESS_VEST, ADDRESS_GLOVE, ADDRESS_GUN]
# name_list = ["Vest1", "Glove1", "Gun1"]
# header_list = [HEADER_VEST, HEADER_GLOVE, HEADER_GUN]

# address_list = [ADDRESS_GLOVE_1]
# name_list = ["Beetle Glove"]
# header_list = [HEADER_GLOVE]

address_list = [ADDRESS_GLOVE_2, ADDRESS_GUN_2, ADDRESS_VEST_1]
name_list = ["GLOVE", "Beetle gun", "vest"]
header_list = [HEADER_GLOVE, HEADER_GUN, HEADER_VEST]

RETRY_COUNT = 8
ultra96_mutex = threading.Lock()

# Helper function for initialisation
def initialise_beetle_list():
    init_params_beetle_list()
    init_connect_beetle_list()
    init_peripheral_beetle_list()
    init_handshake_beetle_list()


# Creates Beetle objects with their corresponding address and name,
# and appends it to the beetle list.
def init_params_beetle_list():
    for addr, name, header in zip(address_list, name_list, header_list):
        beetle_list.append(Beetle(addr, name, header))


def init_connect_beetle_list():
    for beetle in beetle_list:
        beetle.connect_with_retries(RETRY_COUNT)


def init_peripheral_beetle_list():
    for beetle in beetle_list:
        beetle.init_peripheral()


def init_handshake_beetle_list():
    for beetle in beetle_list:
        while not beetle.has_handshake and beetle.is_connected:
            beetle.init_handshake()

            # If beetle disconnects during the handshake, will attempt
            # reconnection and handshaking
            if not beetle.is_connected:
                beetle.connect_with_retries(RETRY_COUNT)

# Helper function to unpack IMU sensor data
def unpack_glove_data_into_dict(glove_data):
    glove_dict = {
        "index": glove_data[0],
        "roll": bytes_to_uint16_t(glove_data[1:3]),
        "pitch": bytes_to_uint16_t(glove_data[3:5]),
        "yaw": bytes_to_uint16_t(glove_data[5:7]),
        "x": bytes_to_uint16_t(glove_data[7:9]),
        "y": bytes_to_uint16_t(glove_data[9:11]),
        "z": bytes_to_uint16_t(glove_data[11:13]),
    }
    return glove_dict

def bytes_to_uint16_t(bytes):
    val = (bytes[0] << 8) + bytes[1]
    if val > MAX_16_BIT_SIGNED:
        val = val - MAX_16_BIT_UNSIGNED
    return val

def pass_params(header, data_obj):
    global next_send
    global send_buf
    now = time.monotonic_ns()

    # grenade, reload, shield
    if header == HEADER_GLOVE:
        data = send_buf.data.add()
        data.player=1
        data.index=data_obj["index"]
        data.roll=data_obj["roll"]
        data.pitch=data_obj["pitch"]
        data.yaw=data_obj["yaw"]
        data.x=data_obj["x"]
        data.y=data_obj["y"]
        data.z=data_obj["z"]

        if now < next_send:
            return
        
        next_send = now + int(20e6)
        resp = stub.Gesture(send_buf)
        send_buf = main_pb2.SensorData()

    # shoot/shot
    # TODO Add Packet IDs for Vest/Gun
    elif header == HEADER_GUN:
        msg = main_pb2.Event(player=1, shootID=data_obj, action=main_pb2.shoot)
        resp = stub.Shoot(msg)
    elif header == HEADER_VEST:
        msg = main_pb2.Event(player=2, shootID=data_obj, action=main_pb2.shot)
        resp = stub.Shot(msg)

def beetle_receiver(beetle):
    while True:
        try:
            beetle.peripheral.waitForNotifications(0.01)
        # Handles disconnect by attempting reconnection/rehandshake
        except btle.BTLEDisconnectError:
            beetle.set_disconnected()
            while not beetle.is_connected:
                beetle.connect_with_retries(RETRY_COUNT)
                print(f"{beetle.name} reconnected. Reinitialising handshake...")
                beetle.init_handshake()

# Thread Worker that relays valid data to Ultra96
def serial_handler(beetle):
    global packet_count
    global expected_index
    
    while True:
        try:
            beetle.peripheral.waitForNotifications(0.01)
            while len(beetle.delegate.data_buffer) >= PACKET_SIZE:
                beetle.delegate.handle_data()
                if beetle.delegate.is_valid_data:
                    if beetle.header == HEADER_GLOVE:
                        glove_data = beetle.delegate.packet[1:14]
                        data_obj = unpack_glove_data_into_dict(glove_data)
                        if (ultra96_mutex.acquire()):
                            pass_params(beetle.header, data_obj)
                            ultra96_mutex.release()
                            
                    else:
                        if ultra96_mutex.acquire():
                            data_obj = beetle.delegate.packet[2]
                            pass_params(beetle.header, data_obj)
                            ultra96_mutex.release()
                
        except btle.BTLEDisconnectError:
            beetle.set_disconnected()
            while not beetle.is_connected:
                beetle.connect_with_retries(RETRY_COUNT)
                print(f"{beetle.name} reconnected. Reinitialising handshake...")
                beetle.init_handshake()

       
      

if __name__ == "__main__":

    beetle_list = []
    initialise_beetle_list()

    # Delay to read the initial print statements
    time.sleep(2)

    start_time = time.time()
    packet_count = 0
    expected_index = 0

    channel = grpc.insecure_channel('localhost:8081')
    stub = main_pb2_grpc.RelayStub(channel)

    # Starts threads
    for beetle in beetle_list:
        handler = threading.Thread(
            target=serial_handler, args=(beetle,))
        # receiver = threading.Thread(
        #     target=beetle_receiver, args=(beetle,))
        handler.start()
        # receiver.start()

    while True:
        pass

