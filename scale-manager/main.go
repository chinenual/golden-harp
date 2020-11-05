package main

import (
	"log"
	"os"
)

func main() {
	if err := LoadSettings(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	if err := LoadConfig(getWorkingDirectory() + "/HarpConfig.xlsx"); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	if err := SerialInit(userSettings.SerialPort, userSettings.SerialBaud); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	if err := SerialWriteCommand([]byte("{cmd: \"version\"")); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	var bytes []byte
	var err error
	if bytes,err = SerialReadResponse(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf(" version: %s\n", string(bytes) )
	return
}
