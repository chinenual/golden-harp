The Arduino device reads serial data from the harp controller (keyboard) and converts it to MIDI.

In performance mode, it need not be connected to a Windows machine; the USB port can be used as its power supply.

Notes pressed on either of the touch sensitive strips are converted to MIDI note on/off events at a fixed velocity of 64. The note conversion depends on the active key/scale/octave settings.

Those settings are selected from whatever presets are loaded and which can be changed via the keyboard on the controller.

In the studio, the Windows-based scale-manager.exe can connect through the USB port to configure scale presets.
