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

  if (debug_midi_enabled) {
    debug_start();
    Serial.print(F("\"MIDI-ON\":\""));
    Serial.print(note, DEC);
    Serial.print(F(" "));
    Serial.print(MIDI_VELOCITY, DEC);
    Serial.print(F("\""));
    debug_end();
  }
}

void midi_note_off(int note, int channel) {
  int opcode = 0x80 | channel;
  midi_serial_out.write(opcode);
  midi_serial_out.write(note);
  midi_serial_out.write((int)0);

  if (debug_midi_enabled) {
    debug_start();
    Serial.print(F("\"MIDI-OFF\":\""));
    Serial.print(note, DEC);
    Serial.print(F(" "));
    Serial.print(0, DEC);
    Serial.print(F("\""));
    debug_end();
  }
}
