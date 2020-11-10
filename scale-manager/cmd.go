package main

import (
	"encoding/json"
	"log"
	"os"
)

const verboseJSON = true

func CmdVersion() (version string, err error) {
	if err = SerialWriteCommand([]byte("{cmd: \"version\"}")); err != nil {
		return
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		return
	}
	if verboseJSON {
		log.Printf("DEBUG: version json: %s\n", string(bytes[:len(bytes)-2]))
	}
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
		return
	}
	if verboseJSON {
		log.Printf("DEBUG: config json: %s\n", string(bytes[:len(bytes)-2]))
	}
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

func CmdSetScale(index int, scale Scale) (err error) {
	val := struct {
		Cmd       string `json:"cmd"`
		N         int    `json:"n"`
		Intervals []int  `json:"i"`
	}{
		Cmd:       "setscale",
		N:         index,
		Intervals: scale.Intervals,
	}
	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	SerialWriteCommand(bytes)
	if bytes, err = SerialReadResponse(); err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	if verboseJSON {
		log.Printf("DEBUG: scale response json: %s\n", string(bytes[:len(bytes)-2]))
	}
	return
}

func CmdSetPreset(index int, preset Preset) (err error) {
	val := struct {
		Cmd    string `json:"cmd"`
		N      int    `json:"n"`
		Preset Preset `json:"preset"`
	}{
		Cmd:    "setpreset",
		N:      index,
		Preset: preset,
	}

	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	SerialWriteCommand(bytes)
	if bytes, err = SerialReadResponse(); err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}
	if verboseJSON {
		log.Printf("DEBUG: preset response json: %s\n", string(bytes[:len(bytes)-2]))
	}
	return
}
