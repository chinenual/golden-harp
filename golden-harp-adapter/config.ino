#define MAX_PRESETS 10
#define MAX_SCALES  10

int r_channel; // MIDI channel for right strip
int l_channel; // MIDI channel for right strip
int r_scale[MAX_R_STRIP - MIN_R_STRIP + 1]; // scaling offsets for the right strip
int l_scale[MAX_L_STRIP - MIN_L_STRIP + 1]; // scaling offsets for the left strip

typedef struct {
  int baseNote;
  int scale;
  int midiChannel;
} stripPreset;

typedef struct {
  stripPreset l_preset;
  stripPreset r_preset;
} preset;

typedef int scaleDefinition[12];

struct {
  int n_scales;
  int n_presets;
  preset presets[MAX_PRESETS];
  scaleDefinition scales[MAX_SCALES];
} config;

void config_setup() {
  //default to two presets - 1: major triad, 2: minor triad

  config.n_scales = 2;
  config.n_presets = 2;

  config.scales[0][0] = 0;
  config.scales[0][1] = 4;
  config.scales[0][2] = 7;
  config.scales[0][3] = -1;

  config.scales[1][0] = 0;
  config.scales[1][1] = 3;
  config.scales[1][2] = 7;
  config.scales[1][3] = -1;

  // major scale based at 2 octaves below middle-C (and the left strip 1 octaves higher)
  config.presets[0].l_preset.baseNote = MIDI_C4 - 24 + 12;
  config.presets[0].l_preset.scale = 0;
  config.presets[0].l_preset.midiChannel = 0; // "All"
  config.presets[0].r_preset.baseNote = MIDI_C4 - 24;
  config.presets[0].r_preset.scale = 0;
  config.presets[0].r_preset.midiChannel = 0; // "All"

  // minor scale based at 2 octaves below middle-C (and the left strip 1 octaves higher)
  config.presets[1].l_preset.baseNote = MIDI_C4 - 24 + 12;
  config.presets[1].l_preset.scale = 1;
  config.presets[1].l_preset.midiChannel = 0; // "All"
  config.presets[1].r_preset.baseNote = MIDI_C4 - 24;
  config.presets[1].r_preset.scale = 1;
  config.presets[1].r_preset.midiChannel = 0; // "All"

  usePreset(0);
}

void scaleInit(int scale[], int numValues, int scaleDefinition[], int baseNote) {
  int j = 0;
  for (int i = 0; i < numValues; i++) {
    if (scaleDefinition[j] < 0) {
      // next octave
      j = 0;
      baseNote += 12;
    }
    scale[i] = baseNote + scaleDefinition[j];
    j++;
  }
}

void usePreset(int num) {
  Serial.print("# PRESET ");
  Serial.print(num, DEC);
  Serial.println();

  if (num > config.n_presets) {
    Serial.println("#  -> no such preset defined. Ignored.");
    return;
  }

  scaleInit(r_scale, MAX_R_STRIP - MIN_R_STRIP,
            config.scales[config.presets[num].r_preset.scale],
            config.presets[num].r_preset.baseNote);
  r_channel = config.presets[num].r_preset.midiChannel;

  scaleInit(l_scale, MAX_L_STRIP - MIN_L_STRIP,
            config.scales[config.presets[num].l_preset.scale],
            config.presets[num].l_preset.baseNote);
  r_channel = config.presets[num].l_preset.midiChannel;
}

