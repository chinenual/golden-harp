package main

import (
	"flag"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"time"
)

const AppName = "Golden Harp Manager"

var applog *log.Logger

var settingsflag = flag.String("settings", "", "Settings file override")
var versionflag = flag.Bool("getversion", false, "get the Arduino's firmware build date")
var getflag = flag.Bool("getconfig", false, "get the config from the attached Arduino")
var setflag = flag.String("setconfig", "", "set the config to the attached Arduino")
var settimingflag = flag.Bool("settiming", false, "set the timing config to the attached Arduino")

func ConnectToArduino() (err error) {
	if !SerialConnected() {
		if err = SerialInit(userSettings.SerialPort, userSettings.SerialBaud); err != nil {
			applog.Printf("ERROR: %v\n", err)
			return
		}
		// give the Arduino time to initialize (connecting seems to cause an unwanted RESET):
		applog.Printf("Waiting for arduino to initialize...\n")
		time.Sleep(time.Second * 5)
	}
	return
}

// platform specific config to ensure logs and preferences go to reasonable locations
func getWorkingDirectory() (path string) {
	// don't do this if we are running from the source tree
	_, err := os.Stat("main.go")
	if !os.IsNotExist(err) {
		// running from source directory
		path = "."
		return
	}

	path, _ = os.UserConfigDir()
	path = path + "/GoldenHarpManager"

	// create it if necessary
	_ = os.MkdirAll(path, os.ModePerm)
	return
}

func initLog() {
	logPath := getWorkingDirectory() + "/harpmgr.log"

	multi := io.MultiWriter(
		&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    5, // megabytes
			MaxBackups: 2,
			Compress:   false,
		},
		os.Stderr)
	applog = log.New(multi, "", log.LstdFlags)
}

func main() {
	flag.Parse()
	initLog()
	applog.Printf("Starting " + AppName + " version " + Version)

	var err error
	if err = LoadSettings(*settingsflag); err != nil {
		applog.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	ranACommand := false

	if *versionflag {
		ranACommand = true
		if err = ConnectToArduino(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		var version string
		var timestamp string
		if version, timestamp, err = CmdVersion(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		applog.Printf("Arduino version: %s, build date: %s\n", version, timestamp)
		fmt.Printf("Arduino version: %s, build date: %s\n", version, timestamp)
	}

	if *setflag != "" {
		ranACommand = true
		if err = LoadConfig(*setflag); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		if err = ConnectToArduino(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		for i, _ := range packedScales {
			CmdSetScale(len(packedScales), i, packedScales[i])
		}
		for i, _ := range packedPresets {
			CmdSetPreset(len(packedPresets), i, packedPresets[i])
		}
	}
	if *settimingflag {
		ranACommand = true
		if err = ConnectToArduino(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		if err = CmdSetTimingParams(userSettings.MaxNoteLen, userSettings.LoopTime); err != nil {
			applog.Printf("ERROR: could not set timing params: %v\n", err)
		}

	}
	if *getflag {
		ranACommand = true
		if err = ConnectToArduino(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		var presets []Preset
		var scales []Scale
		var maxNoteLen int
		var loopTime int
		if presets, scales, maxNoteLen, loopTime, err = CmdGetConfig(); err != nil {
			applog.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		applog.Printf("presets: %#v\n", presets)
		applog.Printf("scales: %#v\n", scales)
		applog.Printf("maxnotelen: %d looplen: %d\n", maxNoteLen, loopTime)
		fmt.Printf("presets: %#v\n", presets)
		fmt.Printf("scales: %#v\n", scales)
		fmt.Printf("maxnotelen: %d looplen: %d\n", maxNoteLen, loopTime)
		return
	}
	if !ranACommand {
		WindowsUI()
	}
}
