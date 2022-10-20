#include <IRremote.hpp>
#include "pitches.h"

///////////////////////// MELODY FOR GUN TRIGGER //////////////////////////
// notes in the melody:
int melody[] = {
  NOTE_C4, NOTE_G3, NOTE_G3, NOTE_A3, NOTE_G3, 0, NOTE_B3, NOTE_C4
};

// note durations: 4 = quarter note, 8 = eighth note, etc.:
int noteDurations[] = {
  4, 8, 8, 4, 4, 4, 4, 4
};
///////////////////////////////////////////////////////////////////////////

////////////////////////// HARDWARE OUTPUT PIN ////////////////////////////
#define TRIGGER     5 // actually an input trigger
#define LED_PIN     13
#define IR_PIN      3
#define BUZZER_PIN  2
///////////////////////////////////////////////////////////////////////////

///////////////////////// PLAYER IDENTIFICATION ///////////////////////////
// Player 1 
const uint16_t PLAYER_1_ADDRESS = 0x0102;

// Player 2 -> this sketch
const uint16_t PLAYER_2_ADDRESS = 0x0105;
///////////////////////////////////////////////////////////////////////////


//////////////////////// BUTTON DEBOUNCE RELATED //////////////////////////
int buttonState;
int lastButtonState = LOW;

// Time variables to keep track of voltage level change
uint32_t lastDebounceTime = 0;
uint32_t debounceDelay = 50;
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
// limit each shot to only 1 trigger press
bool isTriggered = false;

// Play the shooting tune once
bool shootTune = false;
///////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////
/////////////////////////////// BLE RELATED ///////////////////////////////
#define TIMEOUT 200
#define PACKET_SIZE 15
uint8_t shootID;
bool has_handshake;
bool has_ack;
char seq_no;
const char ackPacket[] = {'A', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', 'q'};

//Assembles the 15 byte packet and sends it out over serial
void assemble_and_send_data(uint8_t shootID) {
    char packet[PACKET_SIZE];

    packet[0] = 'F';
    packet[1] = seq_no;
    packet[2] = shootID;

    for (int i = 3; i < PACKET_SIZE - 1; i++) {
      packet[i] = '1'; // padding with ASCII '1'
    }

    char checksum = 0;
    for (int i = 0; i < PACKET_SIZE - 1; i++) {
      checksum ^= packet[i];
    }

    packet[PACKET_SIZE - 1] = checksum;
    
    Serial.write(packet, PACKET_SIZE);
}
///////////////////////////////////////////////////////////////////////////


void setup() {
  pinMode(TRIGGER, INPUT);
  pinMode(LED_PIN, OUTPUT);
  pinMode(BUZZER_PIN, OUTPUT);

  shootID = 0; // set to 0 at the beginning 

  Serial.begin(115200);

  //Initialises flags
  has_handshake = false;
  has_ack = false;
  seq_no = '0';

  //Initial Handshaking
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
  
  
  /////////////////////////
  has_ack = true; // initially set to true to detect first shot
  /////////////////////////
  
  IrSender.begin(IR_PIN);
}

void loop() {
//////////////////////// BUTTON DEBOUNCE ROUTINE //////////////////////////
  int reading = digitalRead(TRIGGER);

  // if voltage level is changed (due to noise or bouncing or real press)
  if (reading != lastButtonState) {
    // reset debounce timer
    lastDebounceTime = millis();
  }

  if ((millis() - lastDebounceTime) > debounceDelay) {
    // if the level is long enough
    if (reading != buttonState) {
      buttonState = reading;
    }
  }
///////////////////////////////////////////////////////////////////////////
//TODO: use buttonState for remote IR control
  if (buttonState == HIGH && !isTriggered && has_ack) {
    // increment shootID
    shootID += 1;
    
    IrSender.sendNEC(PLAYER_2_ADDRESS, shootID, 0);
    isTriggered = true;

    // set ack to false to trigger sending of shoot ID
    has_ack = false;

    // play funny tone on the playewr's side
    shootTune = true;
  }

  // re-send every loop if theres no ack yet
  if (!has_ack) {
    assemble_and_send_data(shootID); // to put the shoot id in here
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

  // play the shooting tune
  if (shootTune == true) {
    playTriggerTone();
    shootTune = false;
  }

  if (buttonState == LOW) {
    isTriggered = false;
  }
  // Save the current button reading for next iteration
  lastButtonState = reading;
}


void playTriggerTone(void) {
  for (int thisNote = 0; thisNote < 8; thisNote++) {
    // to calculate the note duration, take one second divided by the note type.
    //e.g. quarter note = 1000 / 4, eighth note = 1000/8, etc.
    int noteDuration = 1000 / (2 * noteDurations[thisNote]);
    tone(BUZZER_PIN, melody[thisNote], noteDuration);

    // to distinguish the notes, set a minimum time between them.
    // the note's duration + 30% seems to work well:
    int pauseBetweenNotes = noteDuration * 1.30;
    delay(pauseBetweenNotes);
    // stop the tone playing:
    noTone(BUZZER_PIN);
  }
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
