An updated interface to [Iasos's](https://iasos.com/) "Golden Harp".

The golden harp hardware was built by a now-defunct company named
UME.  It consisted of a small keyboard with two touch
sensitive strips which the musician can stroke like strumming a harp.
Output from this controller was fed through a proprietary serial
connection to a Commodore 64 where the signals were mapped to MIDI
notes with configurable key, octave and scale settings.

This project replaces the C64 based processing with an Arduino based
adapter.  Scales can be configured and uploaded to the adapter via a
Windows application (which need not be connected to the adapter during
performance). 

Thanks to Don Miller (@no-carrier) for reverse engineering the serial 
interface to the harp controller and creating the first Arduino implementation
to validate the controller serial interface.

Original implementation was a Commodore 64 program by Richard Wolton whose 
source code is lost in the mists of time.
