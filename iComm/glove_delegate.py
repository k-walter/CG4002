from delegate import Delegate

PACKET_SIZE = 15
WIN_SIZE = 15

class GloveDelegate(Delegate):
    def __init__(self, serial_char, header):
        super().__init__(serial_char, header)

    def handleNotification(self, cHandle, data):

        if (len(self.data_buffer) > 0
         or (data[0] == self.header) or data[0] == 65):
            self.data_buffer += data
        
        packet_size = PACKET_SIZE if self.hand_ack else PACKET_SIZE*WIN_SIZE

        if len(self.data_buffer) >= packet_size:
            # Assemble packet (To be sent or dropped)
            self.packet = self.data_buffer[:packet_size]
            self.data_buffer = self.data_buffer[packet_size:]

            print("Assembled Packet: ", self.packet)
            
            if self.__is_window_valid():
                if self.is_ack_pkt():
                    if not self.hand_ack:
                        self.hand_ack = True
                    ## Need to handle case for when relay node send to beetle
                else:
                    if not self.__exists_duplicates():
                        self.is_duplicate_pkt = False
                        self.is_valid_data = True
                        print("DATA OK TO SEND")
                    else:
                        self.is_duplicate_pkt = True
                        self.is_valid_data = False
                        print("DUP PACKET")
                        

            # Invalid data            
            else:
                self.is_valid_data = False
                self.is_duplicate_pkt = False
                print("CORRUPTED PACKET")
        
        # Packet not assembled yet, do not send
        else:
            self.is_valid_data = False
            print("ASSEMBLING PACKET")
        
    def __is_window_valid(self, packet_size):
        for i in range(0, (packet_size), PACKET_SIZE):
            j = i + PACKET_SIZE
            if not is_valid_checksum(self.packet[i:j]):
                return False
        return True

    def __exists_duplicates(self, packet_size):
        for i in range(0, (packet_size), PACKET_SIZE):
            expected_seq_no = i // PACKET_SIZE
            j = i + PACKET_SIZE
            if not is_expected_packet(expected_seq_no, self.packet[i:j]):
                return True
        return False

def is_valid_checksum(packet):
    checksum = 0
    for i in range(PACKET_SIZE-1):
        checksum ^= packet[i]
    return checksum == packet[-1]

def is_expected_packet(seq_no, packet):
    return seq_no == packet[1]

    