
int r_channel; // MIDI channel for right strip
int l_channel; // MIDI channel for right strip
int r_scale[MAX_R_STRIP - MIN_R_STRIP + 1]; // scaling offsets for the right strip
int l_scale[MAX_L_STRIP - MIN_L_STRIP + 1]; // scaling offsets for the left strip

typedef struct strip_preset_s {
  byte base_note;
  byte scale;
  byte midi_channel;
} strip_preset_t;

// SIZEOF(preset) == 7 bytes
typedef struct preset_s {
  byte key;
  strip_preset_t l_preset;
  strip_preset_t r_preset;
} preset_t;

// Encode  the scale via bits (one bit per interval)
// SIZEOF(packed_scale_definition_u) == 2 bytes
typedef union packed_scale_definition_u {
  unsigned int word;
  struct {
    byte upper;
    byte lower;
  };
};

#if CONFIG_IN_EEPROM
// Each arduino has a different amount of EEPROM.
//    The Duemilanove based prototype Don sent Iasos has 512 or 1024 bytes, depending on which variant.
//    My Uno dev system has 1024 bytes
//    The Nano that we'll be using for the production build has 1024 bytes
//
//    7 * 37 = 259 bytes for presets (one per key on the musical keyboard)
//    2 * 74 = 148 bytes for packed scale definitions (worst case a distinct scale for each strip in each preset)
//             407 bytes total
// plenty of the EEPROM free for other config if we need it.

#define MAX_PRESETS 37
#define MAX_SCALES  74

// Use EEPROM.update to try to minimize the absolute number of writes to the EEPROM (which is limited to 100,000 cycles)
#  define config_write_byte(offset_expression, val) EEPROM.update(&(config_for_offset.offset_expression)-(byte*)&config_for_offset,val)
#  define config_read_byte(tgt, offset_expression) EEPROM.get(&(config_for_offset.offset_expression)-(byte*)&config_for_offset,tgt)

#  define config_write_uint16(offset_expression,val)\
{\
  union { byte b[2]; unsigned short v;} u;  \
  u.v = val; \
  EEPROM.update((byte*)&(config_for_offset.offset_expression)-(byte*)&config_for_offset,u.b[0]);\
  EEPROM.update((byte*)&(config_for_offset.offset_expression)-(byte*)&config_for_offset + 1,u.b[1]);\
}
#  define config_read_uint16(tgt,offset_expression)\
{\
  union { byte b[2]; unsigned short v;} u;  \
  EEPROM.get((byte*)&(config_for_offset.offset_expression)-(byte*)&config_for_offset,u.b[0]);\
  EEPROM.get((byte*)&(config_for_offset.offset_expression)-(byte*)&config_for_offset + 1,u.b[1]);\
  tgt = u.v;\
}

#else
// when developing need a smaller number of presets and scales since everything has to fit in memory
#define MAX_PRESETS 4
#define MAX_SCALES  8

// when developing, don't use EEPROM (since it has limited number of write cycles)

#  define config_write_byte(offset_expression, val) config.offset_expression = (val)
#  define config_read_byte(tgt, offset_expression) tgt = config.offset_expression
#  define config_write_uint16(offset_expression, val) config.offset_expression = (val)
#  define config_read_uint16(tgt, offset_expression) tgt = config.offset_expression

#endif

// we look for CONFIG_TELLTALE in the initialized telltale to detect a brand new flashed EEPROM without any config in it

typedef struct config_s {
  byte initialized;
  byte n_scales;
  byte n_presets;
  preset_t presets[MAX_PRESETS];
  packed_scale_definition_u packed_scale_defs[MAX_SCALES];

  // time in milliseconds for each scan of the controller; without this, we sometimes see both ON and OFF
  // events within a millisecond of each other. Tune this so that the controller is responsive, but not
  // spewing a lot of overlapping MIDI events
  byte loop_time_ms;
  
  // time in milliseconds for max amount of time a note can sound
  unsigned short max_note_length_ms;

  // add new config at the end and change the CONFIG_TELLTALE telltale byte
#define CONFIG_TELLTALE_1_0 0xaa
#define CONFIG_TELLTALE_1_1 0xab
#define CONFIG_TELLTALE_CURRENT CONFIG_TELLTALE_1_1
} config_t;

#if CONFIG_IN_EEPROM
config_t config_for_offset; // unused except as a convenient way to compute the offset to a given byte in the EEPROM
#else
config_t config;
#endif

void config_setup() {

  byte telltale;
  config_read_byte(telltale, initialized);
  if ((!CONFIG_IN_EEPROM) || (telltale != CONFIG_TELLTALE_CURRENT)) {
    if (telltale != CONFIG_TELLTALE_1_0 && telltale != CONFIG_TELLTALE_1_1) {

      // If config in EEPROM, don't configure defaults -- trust what's in the EEPROM. (NOTE the very first time the arduino runs with
      // EEPROM enabled, the config will be uninitialized and may contain garbage values.

      config_write_byte(n_scales,  1);
      config_write_byte(n_presets, 1);

      // chromatic
      int intervals0[] = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11};
      pack_scale(8, intervals0, 0);

      // chromatic scale based at 2 octaves below middle-C (and the left strip 1 octaves higher)
      config_write_byte(presets[0].key, 0);
      config_write_byte(presets[0].l_preset.base_note, MIDI_MIDDLE_C - 24 + 12);
      config_write_byte(presets[0].l_preset.scale, 0);
      config_write_byte(presets[0].l_preset.midi_channel, 0); // "All"
      config_write_byte(presets[0].r_preset.base_note, MIDI_MIDDLE_C - 24);
      config_write_byte(presets[0].r_preset.scale, 0);
      config_write_byte(presets[0].r_preset.midi_channel, 0); // "All"

    }

    if (telltale != CONFIG_TELLTALE_1_1) {
      // new config since prior version
      config_write_byte(loop_time_ms, 50);
      config_write_uint16(max_note_length_ms, 2500);
    }
  }

  config_write_byte(initialized, CONFIG_TELLTALE_CURRENT);

  // set the cached values for timing related config:
  byte v;
  config_read_byte(v, loop_time_ms);
  loop_time_ms = v;
  config_read_uint16(max_note_length_ms, max_note_length_ms);

  use_preset(0);
}

void pack_scale(byte num_notes, int intervals[], byte scale_num) {
  packed_scale_definition_u packed;
  packed.word = 0;
  for (int i = 0; i < num_notes; i++) {
    packed.word |= (0x1 << intervals[i]);
  }
  config_write_byte(packed_scale_defs[scale_num].upper, packed.upper);
  config_write_byte(packed_scale_defs[scale_num].lower, packed.lower);
}

byte unpack_scale(byte scale_num, int intervals[]) {
  packed_scale_definition_u packed;
  config_read_byte(packed.upper, packed_scale_defs[scale_num].upper);
  config_read_byte(packed.lower, packed_scale_defs[scale_num].lower);

  byte num_notes = 0;
  for (int i = 0; i < 16; i++) {
    if (packed.word & (0x1 << i)) {
      intervals[num_notes] = i;
      num_notes++;
    }
  }
  return num_notes;
}

void scale_init(byte scale_num, int base_note, int num_values, int scale[]) {
  int scale_octave[12];
  int scale_length = unpack_scale(scale_num, scale_octave);

  int j = 0;
  for (int i = 0; i < num_values; i++) {
    if (j >= scale_length) {
      // next octave
      j = 0;
      base_note += 12;
    }
    scale[i] = base_note + scale_octave[j];
    j++;
    //Serial.print(" ");Serial.print(scale[i],DEC);
  }
  //Serial.println();
}

void use_preset(byte key) {
  //  Serial.print("# PRESET ");
  //  Serial.print(num, DEC);
  //  Serial.println();

  int num = -1;
  byte n_presets;
  config_read_byte(n_presets, n_presets);
  for (byte i = 0; i < n_presets; i++) {
    byte preset_key;
    config_read_byte(preset_key, presets[i].key);
    if (preset_key == key) {
      num = i;
      break;
    }
  }
  if (num < 0) {
    //    Serial.println("#  -> no such preset defined. Ignored.");
    return;
  }
  {
    byte scale_num;
    byte base_note;
    config_read_byte(scale_num, presets[num].r_preset.scale);
    config_read_byte(base_note, presets[num].r_preset.base_note);
    scale_init(scale_num,
               base_note,
               MAX_R_STRIP - MIN_R_STRIP + 1,
               r_scale);
    config_read_byte(l_channel, presets[num].l_preset.midi_channel);
  }
  {
    byte scale_num;
    byte base_note;
    config_read_byte(scale_num, presets[num].l_preset.scale);
    config_read_byte(base_note, presets[num].l_preset.base_note);
    scale_init(scale_num,
               base_note,
               MAX_L_STRIP - MIN_L_STRIP + 1,
               l_scale);
    config_read_byte(r_channel, presets[num].r_preset.midi_channel);
  }
}

void config_print() {
  Serial.print("{");
  config_printScales();
  Serial.print(", ");
  config_printPresets();
  Serial.print(", \"maxnotelen\" : ");
  unsigned short i;
  config_read_uint16(i,max_note_length_ms);
  Serial.print(i);
  Serial.print(", \"looptime\" : ");
  byte v;
  config_read_byte(v,loop_time_ms);
  Serial.print(v);
  Serial.print("}");
}

void config_printScales() {
  Serial.print("\"scales\": [");
  byte n_scales;
  config_read_byte(n_scales, n_scales);
  for (byte i = 0; i < n_scales; i++) {
    if (i != 0) {
      Serial.print(",");
    }
    config_printScale(i);
  }
  Serial.print("]");
}

void config_printScale(byte packedIndex) {
  int intervals[12];

  int n = unpack_scale(packedIndex, intervals);
  Serial.print("{\"i\":[");
  for (int i = 0; i < n; i++) {
    if (i != 0) {
      Serial.print(", ");
    }
    Serial.print(intervals[i]);
  }
  Serial.print("]}");
}

void config_printPresets() {
  Serial.print("\"presets\": [");
  byte n_presets;
  config_read_byte(n_presets, n_presets);
  for (byte i = 0; i < n_presets; i++) {
    if (i != 0) {
      Serial.print(", ");
    }
    config_printPreset(i);
  }
  Serial.print("]");
}

void config_printPreset(int preset_index) {
  byte v;
  Serial.print("{\"key\": ");
  config_read_byte(v, presets[preset_index].key);
  Serial.print(v);
  Serial.print(",\"l\": {");
  Serial.print("\"base\": ");
  config_read_byte(v, presets[preset_index].l_preset.base_note);
  Serial.print(v);
  Serial.print(", \"scale\": ");
  config_read_byte(v, presets[preset_index].l_preset.scale);
  Serial.print(v);
  Serial.print(", \"chan\": ");
  config_read_byte(v, presets[preset_index].l_preset.midi_channel);
  Serial.print(v);
  Serial.print("}, \"r\": {");
  Serial.print("\"base\": ");
  config_read_byte(v, presets[preset_index].r_preset.base_note);
  Serial.print(v);
  Serial.print(", \"scale\": ");
  config_read_byte(v, presets[preset_index].r_preset.scale);
  Serial.print(v);
  Serial.print(", \"chan\": ");
  config_read_byte(v, presets[preset_index].r_preset.midi_channel);
  Serial.print(v);
  Serial.print("}}");
}
