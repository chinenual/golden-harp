package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Settings struct {
	SerialPort string
	SerialBaud uint
	MaxNoteLen int
	LoopTime   int
}

var userSettings = Settings{
	SerialPort: "COM1",
	SerialBaud: 9600,
	MaxNoteLen: 2500,
	LoopTime:   50,
}

var settingsPathname = getWorkingDirectory() + "/settings.json"

func LoadSettings(path string) (err error) {

	if path != "" {
		settingsPathname = path
	}
	_, err = os.Stat(settingsPathname)
	if os.IsNotExist(err) {
		applog.Printf("Settings file (%s) does not exist.  Using defaults %#v\n", settingsPathname, userSettings)
		err = nil
		return
	}

	var b []byte
	if b, err = ioutil.ReadFile(settingsPathname); err != nil {
		applog.Printf("Error loading settings.  Using defaults %#v: %v\n", userSettings, err)
		return
	}
	if err = json.Unmarshal(b, &userSettings); err != nil {
		applog.Printf("Error parsing settings.  Using defaults %#v: %v\n", userSettings, err)
		return
	}
	applog.Printf("Loaded settings %#v from file %s\n", userSettings, settingsPathname)
	return
}

func SaveSettings() (err error) {
	var b []byte
	if b, err = json.MarshalIndent(userSettings, "", " "); err != nil {
		applog.Printf("Error saving settings: %v\n", err)
	}
	applog.Printf("Save settings %#v to file %s\n", userSettings, settingsPathname)
	if err = ioutil.WriteFile(settingsPathname, b, 0644); err != nil {
		applog.Printf("Error saving settings: %v\n", err)
	}
	return
}
