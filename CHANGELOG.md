## 1.5.0

Manager:

* Reinitializes presets and scales vectors when reading a new config file.
* adds explicit error checks for JSON "status: ERROR" responses.

## 1.4.1

Adapter:

* Fixes the firmware upload script and re-versions the firmware

## 1.4.0

Manager:

* adds support for the new minnotelen param

Adapter:

* Added a pull down resistor to the Harp In serial connection
* min_note_length_ms parameter
??????* Verbose debug mode

## 1.3.0

* Adds support for enabling debugging on the arduino 
* Serial config is now async via a goroutine so we can consume DEBUG stream
* Changes Excel parser from github.com/360EntSecGroup-Skylar/excelize to  github.com/uxri/excelize/v2 
 
 Adapter:
 
* Verbose debug mode

## 1.2.0

Manager:

* Invert the order of the presets in the Config XLS file to match the way Iasos views them when playing
* Tolerate presets in any order (use the Preset column to specify the preset number)

## 1.1.0

Manager:

* Support for new timing settings

Adapter: 

* make loop time and max_note_length_ms times configurable via the Manager

## 1.0.0

First release
