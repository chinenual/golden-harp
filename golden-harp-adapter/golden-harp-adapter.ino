// pins for Harp Controller Keyboard DIN connector, left to right, back view, tab on top
// 1 = D2  (DIN pin 1)
// 2 = 5V  (DIN pin 4)
// 3 = GND (DIN pin 2)
// 4 = D8  (DIN pin 5)
// 5 = D3  (DIN pin 3)

// pins for the MIDI OUT connection, again, left to right, back view, tab on top
// 2 = 5V  (DIN pin 4) (via 220ohm resistor)
// 3 = GND (DIN pin 2)
// 4 = D11 (DIN pin 5)

#define DEBUG_INPUT 0
#define VERBOSE_MIDI 0

// We assume a simple Arduino (e.g. Uno) with only a single UART TX/RX pair.   We use that connection for the USB connection
// to the host computer.   We use software serial for both MIDI and the keyboard controller connnection.

// we use EEPROM as non-volatile ram to store preset configuration between reboots
// my Uno development board supports 1024 bytes of EEPROM

#define MIDI_BAUD 31250
#define USB_BAUD  9600 // use low baud for the USB port in attempt to reduce interference with real-time MIDI

#define MIDI_TX_PIN 11

#define KBD_LATCH_PIN 2
#define KBD_CLOCK_PIN 3
#define KBD_READ_PIN 8
#define KBD_READ_PIN_REGISTER PINB

// time in milliseconds for each scan of the controller; without this, we sometimes see both ON and OFF
// events within a millisecond of each other. Tune this so that the controller is responsive, but not
// spewing a lot of overlapping MIDI events
#define LOOP_TIME 20

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
#define MAX_MUSIC_KEYBOARD 96   // the rightmost C

#define MIDI_C4 60 // MIDI note value for middle-C
#define MIDI_VELOCITY 64

void setup()
{

  Serial.begin(USB_BAUD);
  Serial.println("# begin setup");

  midi_setup();
  harpin_setup();
  config_setup(); 
  
  Serial.println("# end setup");
}

void loop()
{
  unsigned long start = millis();

  harpin_scan();


  unsigned long elapsed = millis() - start;
  if (elapsed < LOOP_TIME) {
    delay(LOOP_TIME - elapsed);
  }
}
