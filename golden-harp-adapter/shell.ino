int showVersion(int argc=0, char**argv=NULL)
{
  Serial.println(F( "Harp Controler Version " __DATE__ " " __TIME__ ));
};

int getConfig(int argc, char **argv) {
  config_print();
}

int setScale(int argc, char **argv) {
  int scaleNum = atoi(argv[1]);
  if (scaleNum >= MAX_SCALES) {
    Serial.print("ERROR: MAX_SCALES = ");
    Serial.println(MAX_SCALES,DEC);   
    return; 
  }
  int len = argc - 2;
  int intervals[12];
  for (int i = 0; i < len; i++) {
    intervals[i] = atoi(argv[i+2]);
  }
  packScale(len, intervals, config.packedScaleDefs[scaleNum]);
  if (scaleNum >= config.n_scales) {
    config.n_scales = scaleNum+1;
  }
  Serial.print("OK: set scale ");
  Serial.print(scaleNum,DEC);
  Serial.println("");
}

int setPreset(int argc, char **argv) {
  int presetNum = atoi(argv[1]);
  if (presetNum >= MAX_PRESETS) {
    Serial.print("ERROR: MAX_PRESETS = ");
    Serial.println(MAX_PRESETS,DEC);   
    return; 
  }
  config.presets[presetNum].l_preset.baseNote    = atoi(argv[2]);
  config.presets[presetNum].l_preset.scale       = atoi(argv[3]);
  config.presets[presetNum].l_preset.midiChannel = atoi(argv[4]);
  config.presets[presetNum].r_preset.baseNote    = atoi(argv[5]);
  config.presets[presetNum].r_preset.scale       = atoi(argv[6]);
  config.presets[presetNum].r_preset.midiChannel = atoi(argv[7]);
  
  if (presetNum >= config.n_presets) {
    config.n_presets = presetNum+1;
  }
  Serial.print("OK: set preset ");
  Serial.print(presetNum,DEC);
  Serial.println("");
}

void shell_setup() {
  shell.attach(Serial);
  shell.addCommand(F("version"), showVersion);
  shell.addCommand(F("getconfig"), getConfig);
  shell.addCommand(F("setscale"), setScale);
  shell.addCommand(F("setpreset"), setPreset);
  showVersion();
}
