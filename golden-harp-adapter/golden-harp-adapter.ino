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


// there are three distinct sets of "note" index values - I use the following terms to keep them from being confused:
//   HARDWARE_BYTE:  the bytes read from the serial connection that represent the raw hardware state of the keyboard
//   KEY_INDEX:      the index into the logic "strip" position or the musical keyboard
//   NOTE:           the scaled MIDI "note" value

// Key Indexes for the extent of the two strips and the musical keyboard:
#define MIN_R_STRIP 1           // "blue" 
#define MAX_R_STRIP 29          // "yellow"
#define MIN_L_STRIP 30          // labeled 1 on the keyboard
#define MAX_L_STRIP 56          // labeled 27 on the keyboard
#define MIN_MUSIC_KEYBOARD 60   // the leftmost C
#define MAX_MUSIC_KEYBOARD 97   // the rightmost C


#define MIDI_C4 60 // MIDI note value for middle-C

boolean keyState[64]; // we've sent an ON event for this key
boolean keyScan[64];  // this key was depressed on the current scan

int r_scale[30]; // scaling offsets for the right strip
int l_scale[30]; // scaling offsets for the left strip

int hardwareToKeyTable[64] = {
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
    keyState[i] = false;
  }

  // default to a major scale based at middle-C (and the left strip 2 octaves higher)
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

void addKey(int key) {
  keyScan[key] = true;
}

void convertHardwareByteToKey(byte hardwareByte, int index) {
  for (int count = 0; count < 8; count++) {
    if (hardwareByte & (1 << count)) {
      int lookupIndex = index * 8 + count;
      addKey(hardwareToKeyTable[lookupIndex]);
    }
  }
}

void getScannedKeys(volatile byte hardwareData[]) {
  convertHardwareByteToKey(hardwareData[0], 0);
  for (int i = 9; i < 16; i++) {
    convertHardwareByteToKey(hardwareData[i], i - 8);
  }
}

int scaleNote(int key) {
  if (key >= MIN_R_STRIP && key <= MAX_R_STRIP) {
    return r_scale[key - MIN_R_STRIP];
  } else {
    return l_scale[key - MIN_L_STRIP];
  }
}

void keyOut(int key) {
  if (keyScan[key] && keyState[key]) {
    // already sent - key still depressed
  } else if (keyScan[key]) {
    // detected "key down"
    keyState[key] = true;
    Serial.print("NOTEON ");
    Serial.print(scaleNote(key), DEC);
    Serial.println();
  } else if (keyState[key]) {
    // detected "key up"
    keyState[key] = false;
    Serial.print("NOTEOFF ");
    Serial.print(scaleNote(key), DEC);
    Serial.println();
  }
}
void loop()
{
  for (int i = 0; i < 64; i++) {
    keyScan[i] = false;
  }
  
  digitalWrite(LATCH_PIN, 1);
  digitalWrite(LATCH_PIN, 0);

  int hasData = 0;

  // Read 16 bytes off of the serial port.
  volatile byte hardwareBytes[16]; 
  for (int i = 0; i < 16; i++)
  {
    // Read 8 individual bits and pack them into a single byte.
    hardwareBytes[i] = 0;
    for (int j = 0; j < 8; j++)
    {
      hardwareBytes[i] <<= 1;
      hardwareBytes[i] |= PINB & 0x01;
      digitalWrite(CLOCK_PIN, 0);
      digitalWrite(CLOCK_PIN, 1);
    }
    hasData += hardwareBytes[i] != 0;
  }
  if (DEBUG_INPUT && hasData) {
    for (int i = 0; i < 16; i++)
    {
      Serial.print(hardwareBytes[i], DEC);
      Serial.print(':');
    }

    Serial.print('\n');
  }

  getScannedKeys(hardwareBytes);

  for (int i = 0; i < 64; i++) {
    keyOut(i);
  }
}
