// pins for Harp Controller Keyboard DIN connector, left to right, back view, tab on bottom
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
#define DEBUG_INPUT 0

// raw mapped note values for each position on the Left and Right strips:
#define MIN_R_STRIP 1
#define MAX_R_STRIP 29
#define MIN_L_STRIP 30
#define MAX_L_STRIP 56

#define MIDI_C4 60 // MIDI note value for middle-C

volatile byte bufferA[NUM_OF_BYTES * 2]; // Hopefully enough?

boolean noteState[64]; // we've sent an OFF or ON event for this note?
boolean noteScan[64];  // this note was depressed on the current scan
int r_scale[30]; // scaling offsets for the right strip
int l_scale[30]; // scaling offsets for the left strip

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

// Musical Keyboard mapping
//   C0:  0:0:0:0:0:0:0:0:128:0:0:0:0:0:0:0:
//   Db0: 0:0:128:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   D0:  (missing key)
//   Eb0: 0:0:0:0:0:0:0:128:0:0:0:0:0:0:0:0:
//   E0:  (missing key)
//   F0:  0:0:0:0:128:0:0:0:0:0:0:0:0:0:0:0:
//   Gb0: 0:0:0:128:0:0:0:0:0:0:0:0:0:0:0:0:
//   G0:  0:128:0:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   Ab0: 0:0:0:0:0:0:0:0:4:0:0:0:0:0:0:0:
//   A0:  0:0:4:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   Bf0: 0:0:0:0:0:0:4:0:0:0:0:0:0:0:0:0:
//   B0:  0:0:0:0:0:0:0:4:0:0:0:0:0:0:0:0:
//   C1:  (missing key)
//   Db1: (missing key)
//   D1:  0:0:0:4:0:0:0:0:0:0:0:0:0:0:0:0:
//   Eb1: 0:4:0:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   E1:  (missing key)
//   F1:  0:0:8:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   Gb1: 0:0:0:0:0:0:8:0:0:0:0:0:0:0:0:0:
//   G1:  0:0:0:0:0:0:0:8:0:0:0:0:0:0:0:0:
//   Ab1: 0:0:0:0:0:8:0:0:0:0:0:0:0:0:0:0:
//   A1:  0:0:0:0:8:0:0:0:0:0:0:0:0:0:0:0:
//   Bf1: 0:0:0:8:0:0:0:0:0:0:0:0:0:0:0:0:
//   B1:  (missing key)
//   C2:  0:0:0:0:0:0:0:0:2:0:0:0:0:0:0:0:
//   Db2: 0:0:2:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   D2:  0:0:0:0:0:0:2:0:0:0:0:0:0:0:0:0:
//   Eb2: 0:0:0:0:0:0:0:2:0:0:0:0:0:0:0:0:
//   E2:  (missing key)
//   F2:  0:0:0:0:2:0:0:0:0:0:0:0:0:0:0:0:
//   Gb2: 0:0:0:2:0:0:0:0:0:0:0:0:0:0:0:0:
//   G2:  0:2:0:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   Ab2: 0:0:0:0:0:0:0:0:16:0:0:0:0:0:0:0:
//   A2:  0:0:16:0:0:0:0:0:0:0:0:0:0:0:0:0:
//   Bf2: 0:0:0:0:0:0:16:0:0:0:0:0:0:0:0:0:
//   B2:  0:0:0:0:0:0:0:16:0:0:0:0:0:0:0:0:
//   C3:  0:0:0:0:0:16:0:0:0:0:0:0:0:0:0:0:

void setup()
{
  Serial.begin(9600); //115200);
  pinMode(2, OUTPUT); // latch
  pinMode(3, OUTPUT); // clock
  pinMode(8, INPUT); // read
  digitalWrite(LATCH_PIN, 0);
  digitalWrite(CLOCK_PIN, 1);
  for (int i = 0; i < 64; i++) {
    noteState[i] = false;
  }

  // default to a major scale based at middle-C (and the left strip 2 ocraves higher)
  int scaleDefinition[] = { 0, 2, 4, 5, 7, 9, 11, -1 };
  scaleInit(r_scale, MAX_R_STRIP - MIN_R_STRIP, scaleDefinition, MIDI_C4);
  scaleInit(l_scale, MAX_L_STRIP - MIN_L_STRIP, scaleDefinition, MIDI_C4+24);
}

void scaleInit(int scale[], int numValues, int scaleDefinition[], int baseNote) {
  int j = 0;
  for (int i = 0; i < MAX_R_STRIP - MIN_R_STRIP; i++) {
    if (scaleDefinition[j] < 0) {
      // next octave
      j = 0;
      baseNote += 12;
    }
    scale[i] = baseNote + scaleDefinition[j];
    j++;
  }
}

void addNote(int noteValue, int notes[]) {
  notes[0] += 1;
  int insertIndex = notes[0];
  notes[insertIndex] = noteValue;

  noteScan[noteValue] = true;
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

int scaleNote(int index) {
  if (index >= MIN_R_STRIP && index <= MAX_R_STRIP) {
    return r_scale[index - MIN_R_STRIP];
  } else {
    return l_scale[index - MIN_L_STRIP];
  }
}
void noteOut(int index) {
  if (noteScan[index] && noteState[index]) {
    // already sent - key still depressed
  } else if (noteScan[index]) {
    noteState[index] = true;
    Serial.print("NOTEON ");
    Serial.print(scaleNote(index), DEC);
    Serial.println();
  } else if (noteState[index]) {
    noteState[index] = false;
    Serial.print("NOTEOFF ");
    Serial.print(scaleNote(index), DEC);
    Serial.println();
  }
}
void loop()
{
  for (int i = 0; i < 64; i++) {
    noteScan[i] = false;
  }

  char text[16];
  digitalWrite(LATCH_PIN, 1);
  digitalWrite(LATCH_PIN, 0);

  // Read 16 bytes off of the serial port.
  int hasData = 0;
  for (int i = 0; i < 16; i++)
  {
    // Read 8 individual bits and pack them into a single byte.
    bufferA[i] = 0;
    for (int j = 0; j < 8; j++)
    {
      bufferA[i] <<= 1;
      bufferA[i] |= PINB & 0x01;
      digitalWrite(CLOCK_PIN, 0);
      digitalWrite(CLOCK_PIN, 1);
    }
    hasData += bufferA[i] != 0;
  }
  if (DEBUG_INPUT && hasData) {
    for (int i = 0; i < 16; i++)
    {
      Serial.print(bufferA[i], DEC);
      Serial.print(':');
    }

    Serial.print('\n');
  }

  int notes[60];
  getNoteList(bufferA, notes);

  for (int i = 0; i < 64; i++) {
    noteOut(i);
  }
}
