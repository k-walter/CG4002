from bluepy import btle

PACKET_SIZE = 16

class Delegate(btle.DefaultDelegate):
    def __init__(self, serial_char, header):
        btle.DefaultDelegate.__init__(self)
        self.serial_char = serial_char
        self.header = header
        self.data_buffer = b""
        self.packet = b""
        self.prev_seq_no = None
        self.hand_ack = False       
        self.is_valid_data = False
        self.is_duplicate_pkt = False
        self.corrupt_pkt_count = 0

    # Triggers whenever data comes in to the characteristic
    def handleNotification(self, cHandle, data):
        self.receive_data(data)

    def receive_data(self, data):
        if (len(self.data_buffer) > 0 
        or (data[0] == self.header or data[0] == 65)):
            self.data_buffer += data

    def handle_data(self):
        if len(self.data_buffer) >= PACKET_SIZE:
            # Assemble packet (To be sent or dropped)
            self.packet = self.data_buffer[:PACKET_SIZE]
            self.data_buffer = self.data_buffer[PACKET_SIZE:]

            if not self.is_ack_pkt() and self.header == 77:
                self.__handle_without_ack()
            else:
                self.__handle_with_ack()
        
        # Packet not assembled yet, do not send
        else:
            self.is_valid_data = False
            print("ASSEMBLING PACKET")

    def __handle_without_ack(self):
        if self.__is_valid_checksum():
            self.is_valid_data = True
            self.corrupt_pkt_count = 0
        else:
            self.is_valid_data = False
            self.corrupt_pkt_count += 1
            print("CORRUPTED")
            if self.corrupt_pkt_count >= 10:
                print("Flushing buffer...")
                self.data_buffer = b""

    def __handle_with_ack(self):
        if self.__is_valid_checksum():
            if self.is_ack_pkt():
                if not self.hand_ack:
                    self.hand_ack = True
                ## Need to handle case for when relay node send to beetle
            else:
                if not self.__is_duplicate():
                    self.is_duplicate_pkt = False
                    self.is_valid_data = True
                    self.prev_seq_no = self.packet[1]
                    self.serial_char.write(b'A')
                else:
                    self.is_duplicate_pkt = True
                    self.is_valid_data = False
                    self.serial_char.write(b'A')
                    print("DUP PACKET")    

        # Invalid data            
        else:
            self.is_valid_data = False
            self.is_duplicate_pkt = False
            print("CORRUPTED PACKET")

    def __is_duplicate(self):
        return self.prev_seq_no == self.packet[1]

    #Helper function to calculate and compare checksum
    def __is_valid_checksum(self):
        checksum = 0
        for i in range(PACKET_SIZE-1):
            checksum ^= self.packet[i]
        return checksum == self.packet[-1]

    def is_ack_pkt(self):
        return self.packet[0] == 65


        

       




