package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var settingsflag = flag.String("settings", "", "Settings file override")
var versionflag = flag.Bool("getversion", false, "get the Arduino's firmware build date")
var getflag = flag.Bool("getconfig", false, "get the config from the attached Arduino")
var setflag = flag.String("setconfig", "", "set the config to the attached Arduino")

func ConnectToArduino() (err error) {
	if !SerialConnected() {
		if err = SerialInit(userSettings.SerialPort, userSettings.SerialBaud); err != nil {
			log.Printf("ERROR: %v\n", err)
			return
		}
		// give the Arduino time to initialize (connecting seems to cause an unwanted RESET):
		log.Printf("Waiting for arduino to initialize...\n")
		time.Sleep(time.Second * 5)
	}
	return
}

func main() {
	flag.Parse()

	log.Printf("Starting version " + Version)
	
	var err error
	if err = LoadSettings(*settingsflag); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	ranACommand := false

	if *versionflag {
		ranACommand = true
		if err = ConnectToArduino(); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		var version string
		if version, err = CmdVersion(); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		log.Printf("Arduino build date: %s\n", version)
	}

	if *setflag != "" {
		ranACommand = true
		if err = LoadConfig(*setflag); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		if err = ConnectToArduino(); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		for i, _ := range packedScales {
			CmdSetScale(len(packedScales), i, packedScales[i])
		}
		for i, _ := range packedPresets {
			CmdSetPreset(len(packedPresets), i, packedPresets[i])
		}
	}
	if *getflag {
		ranACommand = true
		if err = ConnectToArduino(); err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

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
	if !ranACommand {
		WindowsUI()
	}
}
