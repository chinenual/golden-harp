// pins for Harp Controller Keyboard DIN connector, left to right, back view, tab on top
// Nano/Stripboard version:
// 1 = D8 (DIN pin 1) LATCH
#define KBD_LATCH_PIN 8
// 2 = 5V   (DIN pin 4) 5V
// 3 = GND  (DIN pin 2) GND
// 4 = D3 (DIN pin 5) READ
#define KBD_READ_PIN 3
// 5 = D2 (DIN pin 3) CLOCK
#define KBD_CLOCK_PIN 2

// pins for the MIDI OUT connection, again, left to right, back view, tab on top
// Nano/Stripboard version:
// 2 = 5V  (DIN pin 4)  5V (via 220ohm resistor)
// 3 = GND (DIN pin 2)  GND 
// 4 = D4 (DIN pin 5)   DATA
#define MIDI_TX_PIN 4

#define CONFIG_IN_EEPROM 1
#define DEBUG_INPUT 0
#define VERBOSE_MIDI 0

#include "version.h"

#include <stddef.h> // for offsetof
#include <ArduinoJson.h>

#include <EEPROM.h>


// We assume a simple Arduino (e.g. Uno) with only a single UART TX/RX pair.   We use that connection for the USB connection
// to the host computer.   We use software serial for both MIDI and the keyboard controller connnection.

// we use EEPROM as non-volatile ram to store preset configuration between reboots
// my Uno development and the Nano "production" board supports 1024 bytes of EEPROM

#define MIDI_BAUD 31250
#define USB_BAUD  9600 // use low baud for the USB port in attempt to reduce interference with real-time MIDI


#define NOTE_ON_LED_PIN 13

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

#define MIDI_MIDDLE_C 60 // MIDI note value for middle-C
#define MIDI_VELOCITY 64

// cached loop_time value so we're not reading EEPROM every time around the loop
unsigned long loop_time_ms;
unsigned short max_note_length_ms;


void setup()
{

  Serial.begin(USB_BAUD);
//  Serial.println("# begin setup");

  usbconfig_setup();
  midi_setup();
  harpin_setup();
  config_setup(); 
  
//  Serial.println("# end setup");
}

void loop()
{
  unsigned long start = millis();

  usbconfig_loop();
  harpin_loop();

  unsigned long elapsed = millis() - start;
  if (elapsed < loop_time_ms) {
    delay(loop_time_ms - elapsed);
  }
}
