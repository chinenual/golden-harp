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

	var presets []Preset
	var scales []Scale
	if presets, scales, err = CmdGetConfig(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf(" presets: %#v\n", presets )
	log.Printf(" scales: %#v\n", scales )
	return
}
