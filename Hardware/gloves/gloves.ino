#include "Wire.h"
//#include <CircularBuffer.h>

//////////////////////// DEBUG /////////////////////////
//#define DEBUG
////////////////////////////////////////////////////////


///////////////// PRE-COMPILE DEFINES //////////////////
#define MPU6050_ADDR    0x68
#define UPDATE_RATE     50        // in Hz
#define LOOP_TIME       ((1.0f / UPDATE_RATE) * 1000000)  // in microseconds
////////////////////////////////////////////////////////

////////////// LOW PASS FILTER PARAMETERS //////////////
#define LOOP_TIME_S     (1.0f / UPDATE_RATE)              // in seconds
#define CUTOFF_FREQ     1                                 // in hz
#define RC_CONST        (1.0f / (6.283185307f * CUTOFF_FREQ))
#define ALPHA_COEFF     ((1.0f * RC_CONST) / (LOOP_TIME_S + RC_CONST))
////////////////////////////////////////////////////////

uint32_t prevTime;   // system uptime

/////////////////////  GYRO DATA ///////////////////////
int16_t gyroX, gyroY, gyroZ, temp, accX, accY, accZ;      // raw 16-bit gyro readings

#define CALIBRATION_SAMPLE_NUM    2000                    // take 2000 samples of gyro readings for calibraiton
int32_t gyroXCal, gyroYCal, gyroZCal;                     // calculated bias values of gyro readings
int currSample;                                           // variable to keep track the current number of samples


int16_t gyroXProcessed, gyroYProcessed, gyroZProcessed,     // processed gyro data
      accXProcessed, accYProcessed, accZProcessed;        // processed accel data
///////////////////////////////////////////////////////


///////////////// MOTION DETECTION ////////////////////
#define MOTION_ACC_MARGIN   3500                           // +- this for motion threshold for accelerometer
#define MOTION_GYRO_MARGIN  12000                            // +- this for motion threshold for gyroscope
float netAcceleration;                                    // pythagorian net acc vector magnitude
float avgNetAcceleration;                                 // average net acceleration upon calibration
///////////////////////////////////////////////////////


////////////// SAMPLE DATA COLLECTION /////////////////
#define MOTION_FIRST_DATA_FRAME_DURATION            1.5f     // in seconds
#define MOTION_SECOND_DATA_FRAME_DURATION           3        // in seconds
#define MOTION_FIRST_DATA_FRAME                     (MOTION_FIRST_DATA_FRAME_DURATION * UPDATE_RATE)
#define MOTION_SECOND_DATA_FRAME                    (MOTION_SECOND_DATA_FRAME_DURATION * UPDATE_RATE)
#define WINDOW_SIZE                                 75
int motionDataCounter = 0;
int windowDataCounter = 0;
bool isDataCollecting = false;
///////////////////////////////////////////////////////


///////////////////////////////////////////////////////
///////////////////// BLE Related /////////////////////
#define TIMEOUT_PARAM 2   // will re-send every TIMEOUT_PARAM * 0.02 seconds
#define MASK_BYTE 0xff

#define BUFFER_SIZE 100
#define PACKET_SIZE 16  // 16 bytes in a packet

const char ackPacket[] = {'A', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', 'A'};
uint8_t sendCounter = 0;

//IMU Data Structure
typedef struct IMUData {
  int16_t roll;
  int16_t pitch;
  int16_t yaw;
  int16_t x_val;
  int16_t y_val;
  int16_t z_val;
} IMUData;

//Buffer to store incoming IMU Readings
//CircularBuffer<IMUData, BUFFER_SIZE> buffer;

// real IMU data
IMUData imu_data;

bool has_handshake;
volatile bool has_ack;
uint8_t seq_no;

//Assembles the 16 byte packet and sends it out over serial
void assemble_and_send_data(IMUData data) {
  char packet[PACKET_SIZE];
  packet[0] = 'M';
  int ptr = 1;
  append_value(packet, motionDataCounter, &ptr);
  append_value(packet, data.roll, &ptr);
  append_value(packet, data.pitch, &ptr);
  append_value(packet, data.yaw, &ptr);
  append_value(packet, data.x_val, &ptr);
  append_value(packet, data.y_val, &ptr);
  append_value(packet, data.z_val, &ptr);

  char checksum = 0;
  
  for (int i = 0; i < PACKET_SIZE-1; i++) {
    checksum ^= packet[i];
  }

  packet[PACKET_SIZE-1] = checksum;

  for (int i = 0; i < PACKET_SIZE; i++) {
    serialSend(packet[i]);
  }
}

//Helper function to convert 16 bit ints to 2 bytes (represented as 2 string characters)
void append_value(char* packet, int16_t val, int *ptr) {
  packet[(*ptr)++] = (val >> 8) & MASK_BYTE;
  packet[(*ptr)++] = (val & MASK_BYTE);
}


//////////////////////////////////////////////////////
//////////////////////////////////////////////////////

void setup() {
  // initialise hardware interfaces
  pinMode(13, OUTPUT);
  //Serial.begin(115200);
  setupSerial(115200);
  Wire.begin();

  // setup the MPU6050
  setupMPU6050();

  // have some delay for the program to ready
  delay(1000);

  // calibrate MPU6050
  calibrateMPU6050();

  //Initial Handshaking

  
  while (!has_handshake) {
    if (serialRead() == 'H') {
      has_handshake = true;
      for (int i = 0; i < PACKET_SIZE; i++) {
        serialSend(ackPacket[i]);
      }
    } 
  }
  

  // enable UART interrupt
  enableSerialInterrupt();

  // initialise to false for the first transmission
  has_ack = false;  
  
  // update the current system time
  prevTime = micros();
}

void loop() {
  // retrieve and pre-process IMU data
  readMPU6050();
  processMPU6050Data();

  collectMotionData();

  
  // led indication that the program is still running
  digitalWrite(13, !digitalRead(13));

  

  if (micros() - prevTime > LOOP_TIME + 50) {
    while (1) {
      digitalWrite(13, !digitalRead(13));
      delay(500);
    }
  }
  
  // keep the update rate to 50Hz
  while (micros() - prevTime < LOOP_TIME); // if 50Hz period is not full

  // update previous time 
  prevTime = micros();
}

void setupMPU6050(void) {
  Wire.beginTransmission(MPU6050_ADDR);
  Wire.write(0x6B);                         // PWR_MGMT_1 register
  Wire.write(0x00);                         // reset the IMU unit
  Wire.endTransmission();

  Wire.beginTransmission(MPU6050_ADDR);
  Wire.write(0x1B);                         // GYRO_CONGIG register
  Wire.write(0x08);                         // +-500dps full scale
  Wire.endTransmission();

  Wire.beginTransmission(MPU6050_ADDR);
  Wire.write(0x1C);                         // ACCEL_CONFIG register
  Wire.write(0x10);                         // +-4g full scale
  Wire.endTransmission();

  Wire.beginTransmission(MPU6050_ADDR);
  Wire.write(0x1A);                         // CONFIG register
  Wire.write(0x03);                         // ~43Hz Low Pass Filter
  Wire.endTransmission();
}

void readMPU6050(void) {
  Wire.beginTransmission(MPU6050_ADDR);
  Wire.write(0x3B);                         // ACCEL_XOUT_H, start reading and increment the address from here
  Wire.endTransmission();
  Wire.requestFrom(MPU6050_ADDR, 14);

  accX = Wire.read() << 8 | Wire.read();
  accY = Wire.read() << 8 | Wire.read();
  accZ = Wire.read() << 8 | Wire.read();

  temp = Wire.read() << 8 | Wire.read();

  gyroX = Wire.read() << 8 | Wire.read();
  gyroY = Wire.read() << 8 | Wire.read();
  gyroZ = Wire.read() << 8 | Wire.read();
}

void calibrateMPU6050(void) {
  // initialise the current sample to 0  
  currSample = 0;

  // initialise the bias values to 0
  gyroXCal = 0; 
  gyroYCal = 0;
  gyroZCal = 0;

  // initialise average net acceleration to 0
  avgNetAcceleration = 0;
  
  // this loop will repeat CALIBRATION_SAMPLE_NUM times
  for (currSample = 0; currSample < CALIBRATION_SAMPLE_NUM; currSample++) {
    // some led flash to indicate we are still in calibration rountine
    if (currSample % 25 == 0) {
      digitalWrite(13, !digitalRead(13));
    }
    
    // retrieve raw IMU readings
    readMPU6050();
    processMPU6050Data();

    // accumulate each reading into the bias values
    gyroXCal += gyroX;
    gyroYCal += gyroY;
    gyroZCal += gyroZ;

    //accumulate each net acc calculation into average values
    avgNetAcceleration += sqrt(accXProcessed * accXProcessed + accYProcessed * accYProcessed + accZProcessed * accZProcessed);

    // have some delay to not make it read all the samples too fast
    delay(4);
  }

  // once the loop is completed, average the accumulated value by the sample number
  gyroXCal /= CALIBRATION_SAMPLE_NUM;   // bias for X axis gyro reading
  gyroYCal /= CALIBRATION_SAMPLE_NUM;   // bias for Y axis gyro reading
  gyroZCal /= CALIBRATION_SAMPLE_NUM;   // bias for Z aXIS gyro reading
  avgNetAcceleration /= CALIBRATION_SAMPLE_NUM;   // Average net acceleration

  
  // turn off LED
  digitalWrite(13, LOW);
}

void processMPU6050Data(void) {
  // if the bias values are calculated
  if (currSample == CALIBRATION_SAMPLE_NUM) {
    gyroXProcessed = gyroX - gyroXCal;            // subtract the raw reading by the bias value
    gyroYProcessed = gyroY - gyroYCal;            // subtract the raw reading by the bias value
    gyroZProcessed = gyroZ - gyroZCal;            // subtract the raw reading by the bias value
  }
  
  accXProcessed = (int16_t) (accXProcessed * 0.8f) + accX * 0.2f;     // Another LPF, 8192 LSB per g (in g)
  accYProcessed = (int16_t) (accYProcessed * 0.8f) + accY * 0.2f;     // Another LPF, 8192 LSB per g (in g)
  accZProcessed = (int16_t) (accZProcessed * 0.8f) + accZ * 0.2f;     // Another LPF, 8192 LSB per g (in g)
}

bool isMotionDetected(void) {
  // find pythagorian net acc vector
  netAcceleration = sqrt(accXProcessed * accXProcessed + accYProcessed * accYProcessed + accZProcessed * accZProcessed);

  // print out debug values
  #ifdef DEBUG
    //Serial.println(netAcceleration);
  #endif

  if (netAcceleration - (avgNetAcceleration + MOTION_ACC_MARGIN) > 1e-5 ||
      netAcceleration - (avgNetAcceleration - MOTION_ACC_MARGIN) < -1e-5 ||
      gyroXProcessed > MOTION_GYRO_MARGIN || gyroXProcessed < -1 * MOTION_GYRO_MARGIN ||
      gyroYProcessed > MOTION_GYRO_MARGIN || gyroYProcessed < -1 * MOTION_GYRO_MARGIN ||
      gyroZProcessed > MOTION_GYRO_MARGIN || gyroZProcessed < -1 * MOTION_GYRO_MARGIN) {
    return true;
  }
  
  return false;
}


void collectMotionData(void) {
  // if the motion is detected from rest (not collecting data) -> start data collecting flags
  if (isMotionDetected() && !isDataCollecting) {
    isDataCollecting = true;
    motionDataCounter = 0;
    windowDataCounter = WINDOW_SIZE;
  }

  // if we are collecting data
  if (isDataCollecting) {
    imu_data.roll = accXProcessed;
    imu_data.pitch = accYProcessed;
    imu_data.yaw = accZProcessed;
    imu_data.x_val = gyroXProcessed;
    imu_data.y_val = gyroYProcessed;
    imu_data.z_val = gyroZProcessed;

    
    assemble_and_send_data(imu_data);
    motionDataCounter++;
    windowDataCounter--;
    
    

    if (isMotionDetected()) {
      windowDataCounter = WINDOW_SIZE;
    }

    if (windowDataCounter == 0) {
      isDataCollecting = false;
    }
  }
}


void setupSerial(uint32_t baud) {
  uint32_t myUBRR = ( 16000000 / 16 / baud ) - 1;

  UCSR0A = 0;

  // set baud rate
  UBRR0H = (unsigned char) 0;
  UBRR0L = (unsigned char) 8; // 115200 baud rate

  // data format: 8N1
  UCSR0C = 0b00000110;

  // receive interrupt, enable receive and transmit
  UCSR0B = 0b00011000;
}

void enableSerialInterrupt(void) {
  UCSR0B |= (1 << RXCIE0);
}


char serialRead(void) {
  if((UCSR0A & 0b10000000) == 0)
    return 0;
  else
    return UDR0; 
}

void serialSend(char data) {
  while((UCSR0A & 0b00100000) == 0);
  UDR0 = data; 
}

ISR(USART_RX_vect)
{
  char hdr;
  hdr = UDR0;

  //In the event that connection loss occured and relay_node re-inits handshake
  if (hdr == 'H') {
    for (int i = 0; i < PACKET_SIZE; i++) {
      serialSend(ackPacket[i]);
    }
  }
}
