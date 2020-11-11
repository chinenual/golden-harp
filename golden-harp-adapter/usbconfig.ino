void usbconfig_setup() {
  //  show_version();
}

int show_version() {
  Serial.print(F("{\"status\": \"OK\", \"version\": \"" __DATE__ " " __TIME__ "\"}"));
}

int set_scale(int scale_index, JsonArray data) {
  if (scale_index >= MAX_SCALES) {
    Serial.print(F("{\status\": \"ERROR\", \"msg\": \"Invalid scale_index. MAX_SCALES = "));
    Serial.print(MAX_SCALES, DEC);
    Serial.print(F("\"}"));
    return;
  }

  int intervals[12];
  copyArray(data, intervals);

  pack_scale(data.size(), intervals, scale_index);
  config_read_byte(byte n_scales,n_scales);
  if (scale_index >= n_scales) {
    config_write_byte(n_scales, scale_index + 1);
  }
  Serial.print(F("{\"status\": \"OK\"}"));
}

int set_preset(byte preset_index, JsonObject cfg) {
  if (preset_index >= MAX_PRESETS) {
    Serial.print(F("{\"status\": \"ERROR\", \"msg\": \"Invalid preset_index. MAX_PRESETS = "));
    Serial.print(MAX_PRESETS, DEC);
    Serial.print(F("\"}"));
    return;
  }
  config_write_byte(presets[preset_index].key, cfg[F("key")]);
  config_write_byte(presets[preset_index].l_preset.base_note   , cfg[F("l")][F("base")]);
  config_write_byte(presets[preset_index].l_preset.scale       , cfg[F("l")][F("scale")]);
  config_write_byte(presets[preset_index].l_preset.midi_channel, cfg[F("l")][F("chan")]);
  config_write_byte(presets[preset_index].r_preset.base_note   , cfg[F("r")][F("base")]);
  config_write_byte(presets[preset_index].r_preset.scale       , cfg[F("r")][F("scale")]);
  config_write_byte(presets[preset_index].r_preset.midi_channel, cfg[F("r")][F("chan")]);

  if (preset_index >= config.n_presets) {
    config_write_byte(n_presets, preset_index + 1);
  }
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

      } else if (doc[F("cmd")] == F("getconfig")) {
        config_print();

      } else if (doc[F("cmd")] == F("setpreset")) {
        set_preset(doc[F("n")].as<int>(), doc[F("preset")]);

      } else if (doc[F("cmd")] == F("setscale")) {
        set_scale(doc[F("n")].as<int>(), doc[F("i")]);

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
