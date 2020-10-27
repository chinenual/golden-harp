// For the serial connection to the keyboard controller:
#include <avr/io.h>
#include <avr/interrupt.h>

#define NUM_KEYS 128

boolean keyState[NUM_KEYS]; // we've sent an ON event for this key
boolean keyScan[NUM_KEYS];  // this key was depressed on the current scan


int hardwareToKeyTable[NUM_KEYS] = {
  /// // key [byte-index, bit]

  -1, //  0 [0, 0] // keystrips
  24, //  1 [0, 1] //...
  25, //  2 [0, 2]
  14, //  3 [0, 3]
  23, //  4 [0, 4]
  3,  //  5 [0, 5]
  4,  //  6 [0, 6]
  13, //  7 [0, 7]

  -1, //  8 [1, 0] // musical keyboard
  91, //  9 [1, 1] //...
  75, // 10 [1, 2]
  83, // 11 [1, 3]
  -1, // 12 [1, 4]
  -1, // 13 [1, 5]
  -1, // 14 [1, 6]
  67, // 15 [1, 7]

  -1, // 16 [2, 0] // musical keyboard
  85, // 17 [2, 1] //...
  69, // 18 [2, 2]
  77, // 19 [2, 3]
  93, // 20 [2, 4]
  -1, // 21 [2, 5]
  -1, // 22 [2, 6]
  61, // 23 [2, 7]

  -1, // 24 [3, 0] // musical keyboard
  90, // 25 [3, 1] //...
  74, // 26 [3, 2]
  82, // 27 [3, 3]
  -1, // 28 [3, 4]
  -1, // 29 [3, 5]
  -1, // 30 [3, 6]
  66, // 31 [3, 7]

  -1, // 32 [4, 0] // musical keyboard
  89, // 33 [4, 1] //...
  73, // 34 [4, 2]
  81, // 35 [4, 3]
  -1, // 36 [4, 4]
  -1, // 37 [4, 5]
  -1, // 38 [4, 6]
  65, // 39 [4, 7]

  -1, // 40 [5, 0] // musical keyboard
  88, // 41 [5, 1] //...
  72, // 42 [5, 2]
  80, // 43 [5, 3]
  96, // 44 [5, 4]
  -1, // 45 [5, 5]
  -1, // 46 [5, 6]
  64, // 47 [5, 7]

  -1, // 48 [6, 0] // musical keyboard
  86, // 49 [6, 1] //...
  70, // 50 [6, 2]
  78, // 51 [6, 3]
  94, // 52 [6, 4]
  -1, // 53 [6, 5]
  -1, // 54 [6, 6]
  62, // 55 [6, 7]

  -1, // 56 [7, 0] // musical keyboard
  87, // 57 [7, 1] //...
  71, // 58 [7, 2]
  79, // 59 [7, 3]
  95, // 60 [7, 4]
  -1, // 61 [7, 5]
  -1, // 62 [7, 6]
  63, // 63 [7, 7]

  -1, // 64 [8, 0] // musical keyboard
  84, // 65 [8, 1] //...
  68, // 66 [8, 2]
  76, // 67 [8, 3]
  92, // 68 [8, 4]
  -1, // 69 [8, 5]
  -1, // 70 [8, 6]
  60, // 71 [8, 7]

  -1, // 72 [9, 0] // keystrips
  22, // 73 [9, 1] //..
  27, // 74 [9, 2]
  12, // 75 [9, 3]
  21, // 76 [9, 4]
  5,  // 77 [9, 5]
  6,  // 78 [9, 6]
  11, // 79 [9, 7]

  -1, // 80 [10, 0] // keystrips
  51, // 81 [10, 1] //..
  55, // 82 [10, 2]
  41, // 83 [10, 3]
  50, // 84 [10, 4]
  34, // 85 [10, 5]
  35, // 86 [10, 6]
  40, // 87 [10, 7]

  -1, // 88 [11, 0] // keystrips
  53, // 89 [11, 1] //..
  54, // 90 [11, 2]
  39, // 91 [11, 3]
  52, // 92 [11, 4]
  36, // 93 [11, 5]
  37, // 94 [11, 6]
  38, // 95 [11, 7]

  -1, // 96 [12, 0] // keystrips
  18, // 97 [12, 1] //..
  29, // 98 [12, 2]
  16, // 99 [12, 3]
  17, // 100 [12, 4]
  30, // 101 [12, 5]
  31, // 102 [12, 6]
  15, // 103 [12, 7]

  -1, // 104 [13, 0] // keystrips
  49, // 105 [13, 1] //..
  56, // 106 [13, 2]
  43, // 107 [13, 3]
  48, // 108 [13, 4]
  32, // 109 [13, 5]
  33, // 110 [13, 6]
  42, // 111 [13, 7]

  -1, // 112 [14, 0] // keystrips
  20, // 113 [14, 1] //..
  28, // 114 [14, 2]
  10, // 115 [14, 3]
  19, // 116 [14, 4]
  7,  // 117 [14, 5]
  8,  // 118 [14, 6]
  9,  // 119 [14, 7]

  -1, // 120 [15, 0] // keystrips
  47, // 121 [15, 1] //..
  26, // 122 [15, 2]
  45, // 123 [15, 3]
  46, // 124 [15, 4]
  1,  // 125 [15, 5]
  2,  // 126 [15, 6]
  44, // 127  [15, 7]
};

// Strips mapping - selected values to sanity check the original strip-only mapping table
//  30: LMIN:  0:0:0:0:0:0:0:0:0:0:0:0:32:0:0:0:  [12 - 32] -> 32+bits -> 32+5 = 37 (sparse variant) new variant: 101
//  56: LMAX: 0:0:0:0:0:0:0:0:0:0:0:0:0:4:0:0:    [13 - 4]  -> 40+bits -> 40+2 = 42
//  1:  RMIN: 0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:32:   [15 - 32] -> 56+bits -> 56+5 = 61
//  28: RMAX-1: 0:0:0:0:0:0:0:0:0:0:0:0:0:0:4:0:  [14 - 4]  -> 48+bits -> 48+2 = 50
//  29: RMAX:  0:0:0:0:0:0:0:0:0:0:0:0:4:0:0:0:   [12 - 4]  -> 32+bits -> 32+2 = 34

// Musical Keyboard mapping
// my test hardware has some missing keys - I am guessing values based on the patterns I see in the non-missing ones:
//                                            [index - value]
//   60: C0:  0:0:0:0:0:0:0:0:128:0:0:0:0:0:0:0:  [8 - 128] 8*8+7 =>
//   61: Db0: 0:0:128:0:0:0:0:0:0:0:0:0:0:0:0:0:  [2 - 128]
//   62: D0:  (missing key)                       [6 - 128]  : GUESS
//   63: Eb0: 0:0:0:0:0:0:0:128:0:0:0:0:0:0:0:0:  [7 - 128]
//   64: E0:  (missing key)                       [5 - 128]  : GUESS
//   65: F0:  0:0:0:0:128:0:0:0:0:0:0:0:0:0:0:0:  [4 - 128]
//   66: Gb0: 0:0:0:128:0:0:0:0:0:0:0:0:0:0:0:0:  [3 - 128]
//   67: G0:  0:128:0:0:0:0:0:0:0:0:0:0:0:0:0:0:  [1 - 128]

//   68: Ab0: 0:0:0:0:0:0:0:0:4:0:0:0:0:0:0:0:    [8 - 4] 8*8+2
//   69: A0:  0:0:4:0:0:0:0:0:0:0:0:0:0:0:0:0:    [2 - 4]
//   70: Bf0: 0:0:0:0:0:0:4:0:0:0:0:0:0:0:0:0:    [6 - 4]
//   71: B0:  0:0:0:0:0:0:0:4:0:0:0:0:0:0:0:0:    [7 - 4]
//   72: C1:  (missing key)                       [5 - 4]  : GUESS
//   73: Db1: (missing key)                       [4 - 4]  : GUESS
//   74: D1:  0:0:0:4:0:0:0:0:0:0:0:0:0:0:0:0:    [3 - 4]
//   75: Eb1: 0:4:0:0:0:0:0:0:0:0:0:0:0:0:0:0:    [1 - 4]

//   76: E1:  (missing key)                       [8 - 8]  : GUESS 8*8+3
//   77: F1:  0:0:8:0:0:0:0:0:0:0:0:0:0:0:0:0:    [2 - 8]
//   78: Gb1: 0:0:0:0:0:0:8:0:0:0:0:0:0:0:0:0:    [6 - 8]
//   79: G1:  0:0:0:0:0:0:0:8:0:0:0:0:0:0:0:0:    [7 - 8]
//   80: Ab1: 0:0:0:0:0:8:0:0:0:0:0:0:0:0:0:0:    [5 - 8]
//   81: A1:  0:0:0:0:8:0:0:0:0:0:0:0:0:0:0:0:    [4 - 8]
//   82: Bf1: 0:0:0:8:0:0:0:0:0:0:0:0:0:0:0:0:    [3 - 8]
//   83: B1:  (missing key)                       [1 - 8]  : GUESS

//   84: C2:  0:0:0:0:0:0:0:0:2:0:0:0:0:0:0:0:    [8 - 2] 8*8+1
//   85: Db2: 0:0:2:0:0:0:0:0:0:0:0:0:0:0:0:0:    [2 - 2]
//   86: D2:  0:0:0:0:0:0:2:0:0:0:0:0:0:0:0:0:    [6 - 2]
//   87: Eb2: 0:0:0:0:0:0:0:2:0:0:0:0:0:0:0:0:    [7 - 2]
//   88: E2:  (missing key)                       [5 - 2]  : GUESS
//   89: F2:  0:0:0:0:2:0:0:0:0:0:0:0:0:0:0:0:    [4 - 2]
//   90: Gb2: 0:0:0:2:0:0:0:0:0:0:0:0:0:0:0:0:    [3 - 2]
//   91: G2:  0:2:0:0:0:0:0:0:0:0:0:0:0:0:0:0:    [1 - 2]

//   92: Ab2: 0:0:0:0:0:0:0:0:16:0:0:0:0:0:0:0:   [8 - 16] 8*8+4
//   93: A2:  0:0:16:0:0:0:0:0:0:0:0:0:0:0:0:0:   [2 - 16]
//   94: Bf2: 0:0:0:0:0:0:16:0:0:0:0:0:0:0:0:0:   [6 - 16]
//   95: B2:  0:0:0:0:0:0:0:16:0:0:0:0:0:0:0:0:   [7 - 16]
//   96: C3:  0:0:0:0:0:16:0:0:0:0:0:0:0:0:0:0:   [5 - 16]



void harpin_setup() {
  pinMode(KBD_LATCH_PIN, OUTPUT);
  pinMode(KBD_CLOCK_PIN, OUTPUT);
  pinMode(KBD_READ_PIN, INPUT);

  digitalWrite(KBD_LATCH_PIN, 0);
  digitalWrite(KBD_CLOCK_PIN, 1);

  for (int i = 0; i < NUM_KEYS; i++) {
    keyState[i] = false;
  }
}


void addKey(int key) {
#if DEBUG_INPUT
  Serial.print("# saw "); Serial.print(key, DEC); Serial.println();
#endif
  keyScan[key] = true;
}

int keyToChannel(int key) {
  if (key >= MIN_R_STRIP && key <= MAX_R_STRIP) {
    return r_channel;
  } else {
    return l_channel;
  }
}

void keyOut(int key) {
  if (keyScan[key] && keyState[key]) {
    // already sent - key still depressed
  } else if (keyScan[key]) {
    // detected "key down"
    keyState[key] = true;
    if (key >= MIN_MUSIC_KEYBOARD && key <= MAX_MUSIC_KEYBOARD) {
      usePreset(key - MIN_MUSIC_KEYBOARD);
    } else {
      midiNoteOn(scaleNote(key), keyToChannel(key));
    }
  } else if (keyState[key]) {
    // no longer "down" - so detected "key up"
    keyState[key] = false;
    if (key >= MIN_MUSIC_KEYBOARD && key <= MAX_MUSIC_KEYBOARD) {
      // nop
    } else {
      midiNoteOff(scaleNote(key), keyToChannel(key));
    }
  }
}

int scaleNote(int key) {
  if (key >= MIN_R_STRIP && key <= MAX_R_STRIP) {
#if DEBUG_INPUT
    Serial.print("# saw R "); Serial.print(key, DEC); Serial.print(" -> "); Serial.println(r_scale[key - MIN_R_STRIP], DEC);
#endif
    return r_scale[key - MIN_R_STRIP];
  } else {
#if DEBUG_INPUT
    Serial.print("# saw L "); Serial.print(key, DEC); Serial.print(" -> "); Serial.println(l_scale[key - MIN_L_STRIP], DEC);
#endif
    return l_scale[key - MIN_L_STRIP];
  }
}


void convertHardwareByteToStripKey(byte hardwareByte, int index) {
  // for each non-zero but in the hardware byte add "key" value indexed by
  //     index * 8 + bit
  if (hardwareByte != 0) {
    // optimization: most bytes are zero - so avoid the inner loop most of the time
    for (int count = 0; count < 8; count++) {
      if (hardwareByte & (1 << count)) {
        int lookupIndex = index * 8 + count;
        addKey(hardwareToKeyTable[lookupIndex]);
      }
    }
  }
}

void getScannedKeys(volatile byte hardwareData[]) {
  // "strip" values are in the 0, 9..16 indexes:
  // "musical keyboard" values are in the 1..8 indexes
  for (int i = 0; i < 16; i++) {
    convertHardwareByteToStripKey(hardwareData[i], i);
  }
}

void harpin_scan() {
  for (int i = 0; i < NUM_KEYS; i++) {
    keyScan[i] = false;
  }

  digitalWrite(KBD_LATCH_PIN, 1);
  digitalWrite(KBD_LATCH_PIN, 0);

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
      hardwareBytes[i] |= KBD_READ_PIN_REGISTER & 0x01;
      digitalWrite(KBD_CLOCK_PIN, 0);
      digitalWrite(KBD_CLOCK_PIN, 1);
    }
#if DEBUG_INPUT
    hasData += hardwareBytes[i] != 0;
#endif
  }
#if DEBUG_INPUT
  if (hasData) {
    Serial.print("#");
    for (int i = 0; i < 16; i++)
    {
      Serial.print(hardwareBytes[i], DEC);
      Serial.print(':');
    }

    Serial.print('\n');
  }
#endif
  getScannedKeys(hardwareBytes);
  for (int i = 0; i < NUM_KEYS; i++) {
    keyOut(i);
  }
}
