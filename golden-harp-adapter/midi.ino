// For the serial connection to MIDI:
#include <SendOnlySoftwareSerial.h>


SendOnlySoftwareSerial midi_serial_out(MIDI_TX_PIN);

void midi_setup() {
  pinMode(MIDI_TX_PIN, OUTPUT);
  midi_serial_out.begin(MIDI_BAUD);
}

void midi_note_on(int note, int channel) {
  int opcode = 0x90 | channel;
  midi_serial_out.write(opcode);
  midi_serial_out.write(note);
  midi_serial_out.write(MIDI_VELOCITY);

#if VERBOSE_MIDI
  Serial.print("# MIDI note on: ");
  Serial.print(opcode, HEX);
  Serial.print(" ");
  Serial.print(note, HEX);
  Serial.print(" ");
  Serial.print(MIDI_VELOCITY, HEX);
  Serial.println();
#endif
}

void midi_note_off(int note, int channel) {
  int opcode = 0x80 | channel;
  midi_serial_out.write(opcode);
  midi_serial_out.write(note);
  midi_serial_out.write((int)0);

#if VERBOSE_MIDI
  Serial.print("# MIDI note off: ");
  Serial.print(opcode, HEX);
  Serial.print(" ");
  Serial.print(note, HEX);
  Serial.print(" ");
  Serial.print(0, HEX);
  Serial.println();
#endif
}

