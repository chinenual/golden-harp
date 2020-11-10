package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var getflag = flag.Bool("getconfig", false, "get the config from the attached Arduino")
var setflag = flag.Bool("setconfig", false, "set the config from the attached Arduino")

func main() {
	flag.Parse()

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
	log.Printf("Waiting for arduino to initialize...\n")
	time.Sleep(time.Second * 5)

	var version string
	if version, err = CmdVersion(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Arduino build date: %s\n", version)

	if *setflag {
		// do these in reverse order to reduce writes to EEPROM for "numScales"
		for i := len(packedScales) - 1; i >= 0; i-- {
			CmdSetScale(i, packedScales[i])
		}
		for i := len(packedPresets) - 1; i >= 0; i-- {
			CmdSetPreset(i, packedPresets[i])
		}
	}
	if *getflag {
		var presets []Preset
		var scales []Scale
		if presets, scales, err = CmdGetConfig(); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		log.Printf("presets: %#v\n", presets)
		log.Printf("scales: %#v\n", scales)
		return
	}
}
