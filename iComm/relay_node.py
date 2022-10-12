from bluepy import btle
from beetle import Beetle
from main import run
import threading
import time

ADDRESS_GUN = "d0:39:72:bf:c8:87"
ADDRESS_VEST = "d0:39:72:bf:c6:51"
ADDRESS_GLOVE = "d0:39:72:bf:c6:47"
ADDRESS_BACKUP = "50:F1:4A:DA:CC:EB"

HEADER_ACK = 65
HEADER_GUN = 71
HEADER_GLOVE = 77
HEADER_VEST = 86

PACKET_SIZE = 15

MAX_16_BIT_UNSIGNED = 65535

count = 0

# address_list = [ADDRESS_VEST, ADDRESS_GLOVE, ADDRESS_GUN]
# name_list = ["Vest1", "Glove1", "Gun1"]
# header_list = [HEADER_VEST, HEADER_GLOVE, HEADER_GUN]

# address_list = [ADDRESS_GLOVE, ADDRESS_VEST]
# name_list = ["Beetle Glove", "Beetle Vest"]
# header_list = [HEADER_GLOVE, HEADER_VEST]

address_list = [ADDRESS_BACKUP]
name_list = ["Beetle Vest"]
header_list = [HEADER_GLOVE]

RETRY_COUNT = 8

data_count = 0
frag_packet_count = 0
drop_packet_count = 0

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
        "roll": bytes_to_uint16_t(glove_data[0:2]),
        "pitch": bytes_to_uint16_t(glove_data[2:4]),
        "yaw": bytes_to_uint16_t(glove_data[4:6]),
        "x": bytes_to_uint16_t(glove_data[6:8]),
        "y": bytes_to_uint16_t(glove_data[8:10]),
        "z": bytes_to_uint16_t(glove_data[10:12]),
    }
    return glove_dict

def bytes_to_uint16_t(bytes):
    val = (bytes[0] << 8) + bytes[1]
    if val > MAX_16_BIT_UNSIGNED:
        val = MAX_16_BIT_UNSIGNED - val
    return val

# Thread Worker that relays valid data to Ultra96
def serial_handler(beetle):
    while True:
        try:
            if beetle.peripheral.waitForNotifications(5):         
                if beetle.delegate.is_valid_data:
                    print("Valid data! Relaying to Ultra96...")
                    print("Data: ", beetle.delegate.packet)
                    if beetle.header == HEADER_GLOVE:
                        glove_data = beetle.delegate.packet[2:14]
                        data_obj = unpack_glove_data_into_dict(glove_data)
                        # if (ultra96_mutex.acquire()):
                        #     run(beetle.header, data_obj)
                        #     ultra96_mutex.release()
                        print(data_obj)
                    beetle.serial_char.write(b'A')           
                else: 
                    if beetle.delegate.is_duplicate_pkt:         
                        beetle.serial_char.write(b'A')            
                                   
                
        # Handles disconnect by attempting reconnection/rehandshake
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

    # Starts threads
    for beetle in beetle_list:
        handler = threading.Thread(
            target=serial_handler, args=(beetle,))
        handler.start()

    while True:
        pass



    # elapsed_time = time.time() - start_time
    #                 print(f"Packet Statuses for {beetle.name}:")
    #                 print("Time Elapsed: ", elapsed_time, end=" ")
    #                 print("Good Packets: ", data_count, end=" ")
    #                 print("Dropped Packets: ", drop_packet_count, end=" ")
    #                 print("Fragmented Packets: ", frag_packet_count, end=" ")
    #                 data_rate = ((data_count*15)/1000) / elapsed_time
    #                 print(f"Data Rate: {data_rate:.3f}kB/s")
