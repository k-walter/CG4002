from bluepy import btle
from beetle import Beetle

import threading
import time
import grpc
import main_pb2
import main_pb2_grpc

ADDRESS_GUN_1 = "d0:39:72:bf:c8:87"
ADDRESS_GUN_2 = "d0:39:72:bf:c6:51"

ADDRESS_GLOVE_1 = "d0:39:72:bf:c6:47"

ADDRESS_VEST_1 = "c4:be:84:20:1b:5e"
ADDRESS_BACKUP = "50:F1:4A:DA:CC:EB"

PACKET_SIZE = 15

HEADER_ACK = 65
HEADER_GUN = 70
HEADER_GLOVE = 77
HEADER_VEST = 86

MAX_16_BIT_SIGNED = 32768
MAX_16_BIT_UNSIGNED = 65535

MY_PORT = 'localhost:8080'

count = 0

# address_list = [ADDRESS_VEST, ADDRESS_GLOVE, ADDRESS_GUN]
# name_list = ["Vest1", "Glove1", "Gun1"]
# header_list = [HEADER_VEST, HEADER_GLOVE, HEADER_GUN]

address_list = [ADDRESS_GLOVE_1]
name_list = ["Beetle Glove"]
header_list = [HEADER_GLOVE]

# address_list = [ADDRESS_GUN_2, ADDRESS_VEST_1]
# name_list = ["Beetle gun", "vest"]
# header_list = [HEADER_GUN, HEADER_VEST]

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

def dummy_dict(index):
    dummy_dict = {
        "index": index,
        "roll": 0,
        "pitch": 0,
        "yaw": 0,
        "x": 0,
        "y": 0,
        "z": 0,
    }
    return dummy_dict

def bytes_to_uint16_t(bytes):
    val = (bytes[0] << 8) + bytes[1]
    if val > MAX_16_BIT_SIGNED:
        val = val - MAX_16_BIT_UNSIGNED
    return val

def pass_params(header, data_obj):

    st = time.monotonic_ns()

    # grenade, reload, shield
    if header == HEADER_GLOVE:
        msg = main_pb2.SensorData(
            player=1,
            index=data_obj["index"],
            roll=data_obj["roll"],
            pitch=data_obj["pitch"],
            yaw=data_obj["yaw"],
            x=data_obj["x"],
            y=data_obj["y"],
            z=data_obj["z"],
        )
        resp = stub.Gesture(msg)
    # shoot/shot
    # TODO Add Packet IDs for Vest/Gun
    elif header == HEADER_GUN:
        msg = main_pb2.Event(player=1, shootID=data_obj, action=main_pb2.shoot)
        resp = stub.Shoot(msg)
    elif header == HEADER_VEST:
        msg = main_pb2.Event(player=2, shootID=data_obj, action=main_pb2.shot)
        resp = stub.Shot(msg)

    # print(f"Received resp {resp}")

    print(f"rtt {(time.monotonic_ns() - st)/1e6}ms")

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
        # print("Buffer Len: ", len(beetle.delegate.data_buffer))
        try:
            #print("waiting")
            beetle.peripheral.waitForNotifications(0.01)
            # beetle.peripheral.waitForNotifications(0.0001)

            while len(beetle.delegate.data_buffer) >= PACKET_SIZE:
                #print(f"buffer size {len(beetle.delegate.data_buffer)}")
                beetle.delegate.handle_data()
                if beetle.delegate.is_valid_data:
                    #packet_count += 1
                    #print("Time Elapsed: ", time.time()-start_time)
                    if beetle.header == HEADER_GLOVE:
                        glove_data = beetle.delegate.packet[1:14]
                        curr_index = beetle.delegate.packet[1]
                        #print("Index: ", curr_index)
                        #print("Valid data! Packet Index: ", curr_index)
                        data_obj = unpack_glove_data_into_dict(glove_data)
                        # if curr_index == 0:
                        #     expected_index = 0
                        # else:
                        #     expected_index += 1

                        st = time.monotonic_ns()
                        if (ultra96_mutex.acquire()):
                            # while curr_index > expected_index:
                            #     pass_params(beetle.header, dummy_dict(expected_index))
                            #     expected_index += 1
                            pass_params(beetle.header, data_obj)
                            ultra96_mutex.release()
                        print(f"sending took {(time.monotonic_ns() - st)/1e6}ms")
                            
                    else:

                        st = time.monotonic_ns()
                        if ultra96_mutex.acquire():
                            data_obj = beetle.delegate.packet[2]
                            pass_params(beetle.header, data_obj)
                            ultra96_mutex.release()
                        
                        print(f"sending took {(time.monotonic_ns() - st)/1e6}ms")
                    #print(data_obj)
                
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

