Serial port names on Steve's mac:

board2: /dev/cu.usbserial-AG0JV596 (first box sent to Iasos)
board3: /dev/cu.usbserial-AG0JV6J0 (second box sent to Iasos)
board4: /dev/cu.usbserial-AG0JV27K (spare board for testing/support)
uno:    /dev/cu.usbserial-AL05OC8S (prototype)

Required libraries:
    ArduinoJson.h - https://arduinojson.org/
    SendOnlySoftwareSerial.h - https://github.com/nickgammon/SendOnlySoftwareSerial

Arduino IDE setup for firmware upload: 
    Tools->Board Arduino Nano
    Tools->Processor ATmega328P (Old bootloader)
