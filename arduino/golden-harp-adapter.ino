// pins for MIDI connector, left to right, back view, tab on bottom
// 1 = D3
// 2 = D8
// 3 = GND
// 4 = 5V
// 5 = D2 

#include <avr/io.h>
#include <avr/interrupt.h>

#define NUM_OF_BYTES 128
#define LATCH_PIN 2
#define CLOCK_PIN 3

volatile byte bitIndex, bitIndexB;
volatile byte bufferA[NUM_OF_BYTES*2];  // Hopefully enough?
volatile byte lastState, currentState;
volatile boolean bitLatched;
volatile byte bufferSelect;
volatile byte switchBuffer;
volatile byte latchDecimator;
volatile boolean okToClock;

void setup() 
{                
  Serial.begin(115200);
  pinMode(2, OUTPUT); // latch
  pinMode(3, OUTPUT); // clock
  pinMode(8, INPUT); // read
  latchDecimator = 0;
  digitalWrite(LATCH_PIN, 0);
  digitalWrite(CLOCK_PIN, 1);
}

int noteTable[64] = {
  -1, //0
  24, //1
  25, //2
  14, //3
  23, //4
   3, //5
   4, //6
  13, //7
  -1, //8
  22, //9
  27, //10
  12, //11
  21, //12
   5, //13
   6, //14
  11, //15
  -1, //16
  51, //17
  55, //18
  41, //19
  50, //20
  34, //21
  35, //22
  40, //23
  -1, //24
  53, //25
  54, //26
  39, //27
  52, //28
  36, //29
  37, //30
  38, //31
  -1, //32
  18, //33
  29, //34
  16, //35
  17, //36
  30, //37
  31, //38
  15, //39
  -1, //40
  49, //41
  56, //42
  43, //43
  48, //44
  32, //45
  33, //46
  42, //47
  -1, //48
  20, //49
  28, //50
  10, //51
  19, //52
   7, //53
   8, //54
   9, //55
  -1, //56
  47, //57
  26, //58
  45, //59
  46, //60
   1, //61
   2, //62
  44, //63
};

void addNote(int noteValue, int notes[]) {
  notes[0] += 1;
  int insertIndex = notes[0];
  notes[insertIndex] = noteValue; 
}

void convertByteToNote(byte hardwareByte, int index, int notes[]) {
   for (int count = 0; count < 8; count++) {
     if (hardwareByte & (1 << count)) {
       int lookupIndex = index * 8 + count;
       addNote(noteTable[lookupIndex], notes);
     } 
   }
}

void getNoteList(volatile byte hardwareData[], int notes[]) {
  notes[0] = 0;
  convertByteToNote(hardwareData[0], 0, notes);
  for (int i = 9; i < 16; i++) {
    convertByteToNote(hardwareData[i], i - 8, notes);
  }
}

void loop() 
{ 
  char text[16];
  digitalWrite(LATCH_PIN, 1);
  digitalWrite(LATCH_PIN, 0);

  // Read 16 bytes off of the serial port.
  for(int i=0; i < 16; i++)
  {
    // Read 8 individual bits and pack them into a single byte.
    //bufferA[i] = 0;
    for(int j=0; j < 8; j++)
    {
      bufferA[i] <<= 1;
      bufferA[i] |= PINB & 0x01;
      digitalWrite(CLOCK_PIN, 0);
      digitalWrite(CLOCK_PIN, 1);
    }
  }
  
  int notes[60];
  getNoteList(bufferA, notes);
  
  int count = notes[0];
  for (int i = 0; i < count; i++) {
    Serial.print(notes[i + 1],DEC); // sends values
    //Serial.print(notes[i],DEC); // sends values, doesnt add 1
    Serial.write(","); // same as .print(",",BYTE); it seems
  }
  Serial.write("*");  // same as .print("*",BYTE); it seems
}

