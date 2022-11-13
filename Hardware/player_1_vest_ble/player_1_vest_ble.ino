#include <IRremote.hpp>

// same protocol as the sender
#define DECODE_NEC

////////////////////// HARDWARE INPUT PIN ///////////////////
#define IR_RECEIVE_PIN  2
//#define VIBRATION_PIN   4
const int LED_BAR[] = {3, 4, 5};
/////////////////////////////////////////////////////////////

///////////////////////// PLAYER IDENTIFICATION ///////////////////////////
// Player 1 -> this sketch
const uint16_t PLAYER_1_ADDRESS = 0x0102;

// Player 2 
const uint16_t PLAYER_2_ADDRESS = 0x0105;
///////////////////////////////////////////////////////////////////////////

/////////////////// VIBRATION MOTOR /////////////////////////
#define VIBRATION_DURATION    200   // in milliseconds
uint32_t prevShotTime;
bool isShot = false;
/////////////////////////////////////////////////////////////

//////////////////////// BLE RELATED ////////////////////////
#define PACKET_SIZE 16
#define TIMEOUT 300

const char ackPacket[] = {'A', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', 'A'};

uint8_t shotID;
bool has_handshake;
bool has_ack;
char seq_no;

void assemble_and_send_data(uint8_t shotID) {
  char packet[PACKET_SIZE];

  packet[0] = 'V';
  packet[1] = seq_no;
  packet[2] = shotID;

  for (int i = 3; i < PACKET_SIZE - 1; i++) {
    packet[i] = '1';  // padding with ASCII '1'
  }

  char checksum = 0;

  for (int i = 0; i < PACKET_SIZE - 1; i++) {
    checksum ^= packet[i];
  }

  packet[PACKET_SIZE - 1] = checksum;

  Serial.write(packet, PACKET_SIZE);
}

/////////////////////////////////////////////////////////////

void setup() {
  // turn on serial to debug
  Serial.begin(115200);

  for (int i = 0; i < 3; i++) {
    pinMode(LED_BAR[i], OUTPUT);
  }

  for (int i = 0; i < 3; i++) {
    digitalWrite(LED_BAR[i], LOW);
  }
  
  //Initialises flags
  has_handshake = false;
  has_ack = true; // set to true first
  seq_no = '0';
  //shotID = 0;

  //Initial Handshaking
  digitalWrite(LED_BAR[2], HIGH);
  while (!has_handshake) {
    if (Serial.available() && Serial.read() == 'H') {
      has_handshake = true;
      Serial.write(ackPacket, PACKET_SIZE);
    } 
  }
  
  // Setup and start periodic timer to check for handshake
  cli();
  setupTimer1();
  startTimer1();
  sei();

  // play led Squence to signal start
  playLEDSequence();
  
  // start the IR receiver
  IrReceiver.begin (IR_RECEIVE_PIN, ENABLE_LED_FEEDBACK);
}

void loop() {
  if (IrReceiver.decode()) {
    
    // enable receiving of next IR signal
    IrReceiver.resume();

    if (IrReceiver.decodedIRData.address == PLAYER_2_ADDRESS && has_ack) {
      // set isSHot to true to trigger vibration routine
      isShot = true;
      shotID = IrReceiver.decodedIRData.command;
      has_ack = false;   
    }
  }

  if (!has_ack) {
    assemble_and_send_data(shotID); // to put the shoot id in here
  }

  //Initialises timer to check for packet timeout
  unsigned long curr_time = millis();

  while (!has_ack && ((millis() - curr_time) < TIMEOUT)) {
    if (Serial.available()) {
      
      //Should read either 'H' for handshake or 'A' for normal ACK
      char hdr = Serial.read();

      //Successfully received ACK, and flips seq_no for next packet
      if (hdr == 'A') {
        has_ack = true;
        seq_no = (seq_no == '0') ? '1' : '0';
      }

      //In the event that connection loss occured and relay_node re-inits handshake
      if (hdr == 'H') {
        Serial.write(ackPacket, PACKET_SIZE);
      }
    }
  }

  // vibration motor routine
  if (isShot == true) {
    playLEDSequence();
    isShot = false;
  }

  if (has_ack) {
    for (int i = 0; i < 3; i++) {
      digitalWrite(LED_BAR[i], LOW);
    }
    digitalWrite(LED_BAR[0], HIGH);
  }
}


void playLEDSequence() {
  for (int i = 0; i < 3; i++) {
    digitalWrite(LED_BAR[i], LOW);
  }

  delay(300);
  
  for (int i = 0; i < 3; i++) {
    digitalWrite(LED_BAR[i], HIGH);
    delay(100);
  }  
  delay(150);
}


//////////////////////////////////////////////////////////////
// Timer to trigger interrupt once every 1 second to check for 
// re-handshake
void setupTimer1(void) {
  TCCR1A = 0; // OC1A and OC1B disconnected

  // 1s per interrupt
  OCR1A = 62500;

  // trigger output compare 1 match interrupt
  TIMSK1 |= 0b10;

  // initial count value to 0
  TCNT1 = 0;
}

void startTimer1(void) {
  // prescaler = 256
  TCCR1B = 0b00001100;
}

ISR(TIMER1_COMPA_vect) {
  if (Serial.available()) {
    //Should read either 'H' for handshake or 'A' for normal ACK
    char hdr = Serial.read();
    if (hdr == 'H') {
      Serial.write(ackPacket, PACKET_SIZE);
    }
  }
}
