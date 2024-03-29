void usbconfig_setup() {
  //  show_version();
}


void debug_start() {
  if (debug_enabled) {
    Serial.print(F("{\"DEBUG\": \"\", "));
  }
}
void debug_end() {
  if (debug_enabled) {
    Serial.print(F("}"));
    Serial.println();
  }
}

void show_version() {
  Serial.print(F("{\"status\": \"OK\", \"version\": \"" VERSION "\", \"timestamp\": \"" __DATE__ " " __TIME__ "\"}"));
}

void set_scale(byte total_n, byte scale_index, JsonArray data) {
  if (scale_index >= MAX_SCALES) {
    Serial.print(F("{\"status\": \"ERROR\", \"msg\": \"Invalid scale_index. MAX_SCALES = "));
    Serial.print(MAX_SCALES, DEC);
    Serial.print(F("\"}"));
    return;
  }

  int intervals[12];
  copyArray(data, intervals);

  pack_scale(data.size(), intervals, scale_index);
  config_write_byte(n_scales, total_n);
  Serial.print(F("{\"status\": \"OK\"}"));
}

void set_preset(byte total_n, byte preset_index, JsonObject cfg) {
  if (preset_index >= MAX_PRESETS) {
    Serial.print(F("{\"status\": \"ERROR\", \"msg\": \"Invalid preset_index. MAX_PRESETS = "));
    Serial.print(MAX_PRESETS, DEC);
    Serial.print(F("\"}"));
    return;
  }
  config_write_byte(presets[preset_index].key, cfg[F("key")]);
  config_write_byte(presets[preset_index].l_preset.base_note, cfg[F("l")][F("base")]);
  config_write_byte(presets[preset_index].l_preset.scale, cfg[F("l")][F("scale")]);
  config_write_byte(presets[preset_index].l_preset.midi_channel, cfg[F("l")][F("chan")]);
  config_write_byte(presets[preset_index].r_preset.base_note, cfg[F("r")][F("base")]);
  config_write_byte(presets[preset_index].r_preset.scale, cfg[F("r")][F("scale")]);
  config_write_byte(presets[preset_index].r_preset.midi_channel, cfg[F("r")][F("chan")]);

  config_write_byte(n_presets, total_n);
  Serial.print(F("{\"status\": \"OK\"}"));
}

void set_timing(unsigned short min_note_length, unsigned short max_note_length, byte loop_time) {
  config_write_uint16(min_note_length_ms, min_note_length);
  // we cache the value for quick retrieval in the loop() routine
  min_note_length_ms = min_note_length;

  config_write_uint16(max_note_length_ms, max_note_length);
  // we cache the value for quick retrieval in the loop() routine
  max_note_length_ms = max_note_length;

  config_write_byte(loop_time_ms, loop_time);
  // we cache the value for quick retrieval in the loop() routine
  loop_time_ms = loop_time_ms;
  Serial.print(F("{\"status\": \"OK\"}"));
}



void usbconfig_loop() {
  if (Serial.available()) {
    StaticJsonDocument<300> doc;

    // Read the JSON document from the "link" serial port
    DeserializationError err = deserializeJson(doc, Serial);

    if (err == DeserializationError::Ok) {
      if (doc[F("cmd")] == F("version")) {
        show_version();

      } else if (doc[F("cmd")] == F("setdebug")) {
        // example:  {"cmd": "setdebug", "gen": true, "midi": true, "hw": true}
        debug_enabled = doc[F("gen")].as<bool>();
        debug_midi_enabled = doc[F("midi")].as<bool>();
        debug_hw_enabled = doc[F("hw")].as<bool>();
        Serial.print(F("{\"status\": \"OK\"}"));

      } else if (doc[F("cmd")] == F("getconfig")) {
        config_print();

      } else if (doc[F("cmd")] == F("setpreset")) {
        set_preset(doc[F("total_n")].as<int>(), doc[F("n")].as<int>(), doc[F("preset")]);

      } else if (doc[F("cmd")] == F("setscale")) {
        set_scale(doc[F("total_n")].as<int>(), doc[F("n")].as<int>(), doc[F("i")]);

      } else if (doc[F("cmd")] == F("settiming")) {
        // example: {"cmd": "settiming", "minnotelen": 180, "maxnotelen": 2200, "looptime": 15}
        set_timing(doc[F("minnotelen")].as<short>(), doc[F("maxnotelen")].as<short>(), doc[F("looptime")].as<int>());

      } else {
        Serial.print(F("{status: \"ERROR\", msg: \"Invalid cmd\"}"));
      }
    } else {
      Serial.print(F("{status: \"ERROR\", msg: \"JSON parse error: "));
      Serial.print(err.c_str());
      Serial.print(F("\"}"));
      // flush the input
      while (Serial.available()) {
        Serial.read();
      }
    }
    // response must be terminated with a newline:
    Serial.println();
  }
}
