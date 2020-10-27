// For the serial connection to MIDI:
#include <SendOnlySoftwareSerial.h>


SendOnlySoftwareSerial midiSerialOut(MIDI_TX_PIN);

void midi_setup() {
  pinMode(MIDI_TX_PIN, OUTPUT);
  midiSerialOut.begin(MIDI_BAUD);
}

void midiNoteOn(int note, int channel) {
  int opcode = 0x90 | channel;
  midiSerialOut.write(opcode);
  midiSerialOut.write(note);
  midiSerialOut.write(MIDI_VELOCITY);

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

void midiNoteOff(int note, int channel) {
  int opcode = 0x80 | channel;
  midiSerialOut.write(opcode);
  midiSerialOut.write(note);
  midiSerialOut.write((int)0);

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

