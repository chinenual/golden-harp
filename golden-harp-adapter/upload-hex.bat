SET COMPORT=COM1
SET HEXFILE=golden-harp-adapter-1_4.ino.hex
"C:\Program Files\Arduino\hardware\tools\avr/bin/avrdude" -C"C:\Program Files\Arduino\hardware\tools\avr/etc/avrdude.conf" -v -patmega328p -carduino -P%COMPORT% -b57600 -D -Uflash:w:%HEXFILE%:i 

