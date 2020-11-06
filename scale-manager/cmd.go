package main

import (
	"encoding/json"
	"log"
	"os"
)

func CmdVersion() (version string, err error) {
	if err = SerialWriteCommand([]byte("{cmd: \"version\"}")); err != nil {
		return
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		return
	}
	log.Printf(" version json: %s\n", string(bytes))

	var data map[string]interface{}

	if err = json.Unmarshal([]byte(bytes), &data); err != nil {
		return
	}

	version = data["version"].(string)
	return
}

func CmdGetConfig() (presets []Preset, scales []Scale, err error) {

	if err = SerialWriteCommand([]byte("{cmd: \"getconfig\"}")); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	log.Printf(" config json: %s\n", string(bytes))

	var config struct {
		Scales  []Scale
		Presets []Preset
	}
	if err = json.Unmarshal([]byte(bytes), &config); err != nil {
		return
	}

	presets = config.Presets
	scales = config.Scales
	return
}
