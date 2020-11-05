package main

import (
	"log"
	"os"
	"time"
)

func main() {
	var err error
	if err = LoadSettings(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	if err = LoadConfig(getWorkingDirectory() + "/HarpConfig.xlsx"); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	if err = SerialInit(userSettings.SerialPort, userSettings.SerialBaud); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	// give the Arduino time to initialize (connecting seems to cause an unwanted RESET):
	time.Sleep(time.Second*5)

	var version string
	if version,err = CmdVersion(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf(" version: %s\n", version )


	if err = SerialWriteCommand([]byte("{cmd: \"getconfig\"}")); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	var bytes []byte
	if bytes,err = SerialReadResponse(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf(" config: %s\n", string(bytes) )

	return
}
