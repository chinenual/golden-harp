package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Settings struct {
	SerialPort string
	SerialBaud uint
}

var userSettings = Settings {
	SerialPort: "COM1",
	SerialBaud: 9600,
}

var settingsPathname = getWorkingDirectory() + "/settings.json"

func getWorkingDirectory() (path string) {
	// don't do this if we are running from the source tree
	_, err := os.Stat("build.bat")
	if !os.IsNotExist(err) {
		// running from source directory
		path = "."
		return
	}
	path, _ = os.UserConfigDir()
	path = path + "/GoldenHarpScaleManager"

	// create it if necessary
	_ = os.MkdirAll(path, os.ModePerm)
	return
}

func LoadSettings() (err error) {

	_, err = os.Stat(settingsPathname)
	if os.IsNotExist(err) {
		log.Printf("Settings file (%s) does not exist.  Using defaults %#v\n", settingsPathname, userSettings)
		err = nil
		return
	}

	var b []byte
	if b, err = ioutil.ReadFile(settingsPathname); err != nil {
		log.Printf("Error loading settings.  Using defaults %#v: %v\n", userSettings, err)
		return
	}
	if err = json.Unmarshal(b, &userSettings); err != nil {
		log.Printf("Error parsing settings.  Using defaults %#v: %v\n", userSettings, err)
		return
	}
	log.Printf("Loaded settings %#v from file %s\n", userSettings, settingsPathname)
	return
}

func SaveSettings() (err error) {
	var b []byte
	if b, err = json.MarshalIndent(userSettings, "", " "); err != nil {
		log.Printf("Error saving settings: %v\n", err)
	}
	log.Printf("Save settings %#v to file %s\n", userSettings, settingsPathname)
	if err = ioutil.WriteFile(settingsPathname, b, 0644); err != nil {
		log.Printf("Error saving settings: %v\n", err)
	}
	return
}
