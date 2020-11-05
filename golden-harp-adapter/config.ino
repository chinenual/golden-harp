int r_channel; // MIDI channel for right strip
int l_channel; // MIDI channel for right strip
int r_scale[MAX_R_STRIP - MIN_R_STRIP + 1]; // scaling offsets for the right strip
int l_scale[MAX_L_STRIP - MIN_L_STRIP + 1]; // scaling offsets for the left strip

struct stripPreset {
  byte baseNote;
  byte scale;
  byte midiChannel;
};

// SIZEOF(preset) == 6 bytes
struct preset {
  stripPreset l_preset;
  stripPreset r_preset;
};

// SIZEOF(packedScaleDefinition) == 6 bytes
typedef byte packedScaleDefinition[6];

// so if we map all keys to trigger presets tthere are 37 presets. If each has two distict scales, we need 2*37 scales
// so 666 bytes.  Wont fit in 512byte EEPROM.  So we'll go with max of 48 scales to fit.
//    6 * 37 = 222 bytes for presets
//    6 * 48 = 288 bytes for packed scale definitions
//             510 bytes total
#define CONFIG_IN_EEPROM 0
#if CONFIG_IN_EEPROM
#define MAX_PRESETS 37
#define MAX_SCALES  48
#else
// when developing, don't use EEPROM (since it has limited number of write cycles) but reduce number of bytes needed in order to have enough free RAM
#define MAX_PRESETS 10
#define MAX_SCALES  10
#endif

struct {
  byte n_scales;
  byte n_presets;
  preset presets[MAX_PRESETS];
  packedScaleDefinition packedScaleDefs[MAX_SCALES];
} config;

void config_setup() {
  //default to two presets - 1: major triad, 2: minor triad

  config.n_scales = 2;
  config.n_presets = 2;

  // chromatic
  int intervals0[] = {0, 1, 2, 3, 4, 5, 6, 7};
  packScale(8, intervals0, config.packedScaleDefs[0]);
  // "pipes of pan - IV"
  int intervals1[] = {0, 2, 4, 7, 10, 11};
  packScale(6, intervals1, config.packedScaleDefs[1]);

  // chromatic scale based at 2 octaves below middle-C (and the left strip 1 octaves higher)
  config.presets[0].l_preset.baseNote = MIDI_C4 - 24 + 12;
  config.presets[0].l_preset.scale = 0;
  config.presets[0].l_preset.midiChannel = 0; // "All"
  config.presets[0].r_preset.baseNote = MIDI_C4 - 24;
  config.presets[0].r_preset.scale = 0;
  config.presets[0].r_preset.midiChannel = 0; // "All"

  // "Pipes of Pan - IV" scale based at 2 octaves below middle-C (and the left strip 1 octaves higher)
  config.presets[1].l_preset.baseNote = MIDI_C4 - 24 + 12;
  config.presets[1].l_preset.scale = 1;
  config.presets[1].l_preset.midiChannel = 0; // "All"
  config.presets[1].r_preset.baseNote = MIDI_C4 - 24;
  config.presets[1].r_preset.scale = 1;
  config.presets[1].r_preset.midiChannel = 0; // "All"

  usePreset(0);
}

void packScale(byte numNotes, int intervals[], byte packed[]) {
  // first nibble of the packed array is a note count
  // each following nibble is an interval.
  packed[0] = numNotes;
  int j = 0;
  for (int i = 0; i < numNotes; i++) {
    if (i % 2) {
      // even
      packed[j] = intervals[i];
    } else {
      // odds go into the upper nibble
      packed[j] |= (intervals[i] << 4);
      // next nibble will be in the next packed byte
      j++;
    }
  }
}

int unpackScale(byte packed[], int intervals[]) {
  int numNotes = packed[0] & 0x0f;
  int j = 0;
  //Serial.print("UNPACKED ");
  for (int i = 0; i < numNotes; i++) {
    if (i % 2) {
      // even
      intervals[i] = packed[j] & 0x0f;
    } else {
      // odds go into the upper nibble
      intervals[i] = (packed[j] >> 4) & 0x0f;
      // next nibble will be in the next packed byte
      j++;
    }
    //Serial.print(" ");Serial.print(intervals[i],DEC);
  }
  //Serial.println();
  return numNotes;
}

void scaleInit(byte packed[], int baseNote, int numValues, int scale[]) {
  int scaleOctave[12];
  int scaleLength = unpackScale(packed, scaleOctave);

  int j = 0;
  for (int i = 0; i < numValues; i++) {
    if (j >= scaleLength) {
      // next octave
      j = 0;
      baseNote += 12;
    }
    scale[i] = baseNote + scaleOctave[j];
    j++;
    //Serial.print(" ");Serial.print(scale[i],DEC);
  }
  //Serial.println();
}

void usePreset(int num) {
//  Serial.print("# PRESET ");
//  Serial.print(num, DEC);
//  Serial.println();

  if (num > config.n_presets) {
//    Serial.println("#  -> no such preset defined. Ignored.");
    return;
  }

  scaleInit(config.packedScaleDefs[config.presets[num].r_preset.scale],
            config.presets[num].r_preset.baseNote,
            MAX_R_STRIP - MIN_R_STRIP,
            r_scale);
  r_channel = config.presets[num].r_preset.midiChannel;

  scaleInit(config.packedScaleDefs[config.presets[num].l_preset.scale],
            config.presets[num].l_preset.baseNote,
            MAX_L_STRIP - MIN_L_STRIP,
            l_scale);
  r_channel = config.presets[num].l_preset.midiChannel;
}

void config_print() {
  Serial.print("{");
  config_printScales();
  Serial.print(", ");
  config_printPresets();
  Serial.print("}");
}

void config_printScales() {
  Serial.print("scales: [");
  for (int i = 0; i < config.n_scales; i++) {
    if (i != 0) {
      Serial.print(",");
    }
    config_printScale(config.packedScaleDefs[i]);
  }
  Serial.print("]");
}

void config_printScale(byte packed[]) {
  int intervals[12];
  int n = unpackScale(packed, intervals);
  Serial.print("[");
  for (int i = 0; i < n; i++) {
    if (i != 0) {
      Serial.print(", ");
    }
    Serial.print(intervals[i]);
  }
  Serial.print("]");
}

void config_printPresets() {
  Serial.print("\"presets\": [");
  for (int i = 0; i < config.n_presets; i++) {
    if (i != 0) {
      Serial.print(", ");
    }
    config_printPreset(config.presets[i]);
  }
  Serial.print("]");
}

void config_printPreset(struct preset p) {
  Serial.print("{\"l\": {");
  Serial.print("\"base\": ");
  Serial.print(p.l_preset.baseNote);
  Serial.print(", \"scale\": ");
  Serial.print(p.l_preset.scale);
  Serial.print(", \"chan\": ");
  Serial.print(p.l_preset.midiChannel);
  Serial.print("} \"r\": {");
  Serial.print("\"base\": ");
  Serial.print(p.r_preset.baseNote);
  Serial.print(", \"scale\": ");
  Serial.print(p.r_preset.scale);
  Serial.print(", \"chan\": ");
  Serial.print(p.r_preset.midiChannel);
  Serial.print("}}");
}
