SET COMPORT=COM3
SET HEXFILE=golden-harp-adapter-1_4.ino.hex
"C:\Arduino\arduino-1.3.0\hardware\tools\avr\bin\avrdude" -C"C:\Arduino\arduino-1.3.0\hardware\tools\avr\etc\avrdude.conf" -v -patmega328p -carduino -P%COMPORT% -b57600 -D -Uflash:w:%HEXFILE%:i 

