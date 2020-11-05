void usbconfig_setup() {
//  showVersion();
}

int showVersion() {
  Serial.println(F("{status: \"OK\", version: \"" __DATE__ " " __TIME__ "\"}"));
}

int setScale(int scaleNum, JsonArray data) {
  if (scaleNum >= MAX_SCALES) {
    Serial.print(F("{status: \"ERROR\", msg: \"Invalid scaleNum. MAX_SCALES = "));
    Serial.print(MAX_SCALES, DEC);
    Serial.println(F("\"}"));
    return;
  }

  int intervals[12];
  copyArray(data, intervals);

  packScale(data.size(), intervals, config.packedScaleDefs[scaleNum]);
  if (scaleNum >= config.n_scales) {
    config.n_scales = scaleNum + 1;
  }
  Serial.println(F("{status: \"OK\"}"));
}

int setPreset(int presetNum, JsonObject cfg) {
  if (presetNum >= MAX_PRESETS) {
    Serial.print(F("{status: \"ERROR\", msg: \"Invalid presetNum. MAX_PRESETS = "));
    Serial.print(MAX_PRESETS, DEC);
    Serial.println(F("\"}"));
    return;
  }
  config.presets[presetNum].l_preset.baseNote    = cfg[F("l")][F("base")];
  config.presets[presetNum].l_preset.scale       = cfg[F("l")][F("scale")];
  config.presets[presetNum].l_preset.midiChannel = cfg[F("l")][F("chan")];
  config.presets[presetNum].r_preset.baseNote    = cfg[F("r")][F("base")];
  config.presets[presetNum].r_preset.scale       = cfg[F("r")][F("scale")];
  config.presets[presetNum].r_preset.midiChannel = cfg[F("r")][F("chan")];

  if (presetNum >= config.n_presets) {
    config.n_presets = presetNum + 1;
  }
  Serial.println(F("{status: \"OK\"}"));
}

void usbconfig_loop() {
  if (Serial.available()) {
    StaticJsonDocument<300> doc;

    // Read the JSON document from the "link" serial port
    DeserializationError err = deserializeJson(doc, Serial);

    if (err == DeserializationError::Ok) {
      if (doc[F("cmd")] == "version") {
        showVersion();

      } else if (doc[F("cmd")] == F("getconfig")) {
        config_print();

      } else if (doc[F("cmd")] == "setpreset") {
        setPreset(doc["presetNum"].as<int>(), doc[F("preset")]);

      } else if (doc[F("cmd")] == "setscale") {
        setScale(doc[F("scaleNum")].as<int>(), doc[F("intervals")]);

      } else {
        Serial.print(F("{status: \"ERROR\", msg: \"Invalid cmd\"}"));

      }
    } else {
      Serial.print(F("{status: \"ERROR\", msg: \"JSON parse error: "));
      Serial.print(err.c_str());
      Serial.println(F("\"}"));
      // flush the input
      while (Serial.available()) {
        Serial.read();
      }
    }
  }
}
