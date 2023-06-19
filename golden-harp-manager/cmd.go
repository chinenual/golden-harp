package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func CmdVersion() (version string, timestamp string, err error) {
	if err = SerialWriteCommand([]byte("{cmd: \"version\"}")); err != nil {
		return
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		return
	}
	var data map[string]interface{}

	if err = json.Unmarshal(bytes, &data); err != nil {
		return
	}

	version = data["version"].(string)
	timestamp = data["timestamp"].(string)
	return
}
func CmdGetConfig() (presets []Preset, scales []Scale, minNoteLen int, maxNoteLen int, loopTime int, err error) {

	if err = SerialWriteCommand([]byte("{cmd: \"getconfig\"}")); err != nil {
		applog.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		applog.Printf("ERROR: %v\n", err)
		return
	}
	var config struct {
		Scales     []Scale
		Presets    []Preset
		MinNoteLen int `json:"minnotelen"`
		MaxNoteLen int `json:"maxnotelen"`
		LoopTime   int `json:"looptime"`
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		return
	}

	presets = config.Presets
	scales = config.Scales
	maxNoteLen = config.MaxNoteLen
	minNoteLen = config.MinNoteLen
	loopTime = config.LoopTime
	return
}

func CmdSetTimingParams(minNoteLen int, maxNoteLen int, loopTime int) (err error) {
	val := struct {
		Cmd        string `json:"cmd"`
		MinNoteLen int    `json:"minnotelen"`
		MaxNoteLen int    `json:"maxnotelen"`
		LoopTime   int    `json:"looptime"`
	}{
		Cmd:        "settiming",
		MinNoteLen: minNoteLen,
		MaxNoteLen: maxNoteLen,
		LoopTime:   loopTime,
	}
	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	if err = SerialWriteCommand(bytes); err != nil {
		return
	}
	if bytes, err = SerialReadResponse(); err != nil {
		applog.Printf("ERROR: %v\n", err)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return
	}

	if data["status"].(string) != "OK" {
		err = errors.Errorf("%s", data["msg"].(string))
		applog.Printf("ERROR: %v\n", err)
		return
	}
	return

}

func CmdSetScale(total int, index int, scale Scale) (err error) {
	val := struct {
		Cmd       string `json:"cmd"`
		Total     int    `json:"total_n"`
		N         int    `json:"n"`
		Intervals []int  `json:"i"`
	}{
		Cmd:       "setscale",
		Total:     total,
		N:         index,
		Intervals: scale.Intervals,
	}
	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	if err = SerialWriteCommand(bytes); err != nil {
		return
	}
	if bytes, err = SerialReadResponse(); err != nil {
		applog.Printf("ERROR: %v\n", err)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return
	}

	if data["status"].(string) != "OK" {
		err = errors.Errorf("%s", data["msg"].(string))
		applog.Printf("ERROR: %v\n", err)
		return
	}
	return
}

func CmdSetPreset(total int, index int, preset Preset) (err error) {
	val := struct {
		Cmd    string `json:"cmd"`
		Total  int    `json:"total_n"`
		N      int    `json:"n"`
		Preset Preset `json:"preset"`
	}{
		Cmd:    "setpreset",
		Total:  total,
		N:      index,
		Preset: preset,
	}

	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	if err = SerialWriteCommand(bytes); err != nil {
		return
	}
	if bytes, err = SerialReadResponse(); err != nil {
		applog.Printf("ERROR: %v\n", err)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return
	}

	if data["status"].(string) != "OK" {
		err = errors.Errorf("%s", data["msg"].(string))
		applog.Printf("ERROR: %v\n", err)
		return
	}
	return
}
func CmdSetDebug(general, midi, hw bool) (err error) {
	val := struct {
		Cmd  string `json:"cmd"`
		Gen  bool   `json:"gen"`
		Midi bool   `json:"midi"`
		Hw   bool   `json:"hw"`
	}{
		Cmd:  "setdebug",
		Gen:  general,
		Midi: midi,
		Hw:   hw,
	}

	var bytes []byte
	if bytes, err = json.Marshal(val); err != nil {
		return
	}
	if err = SerialWriteCommand(bytes); err != nil {
		return
	}
	if bytes, err = SerialReadResponse(); err != nil {
		applog.Printf("ERROR: %v\n", err)
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return
	}

	if data["status"].(string) != "OK" {
		err = errors.Errorf("%s", data["msg"].(string))
		applog.Printf("ERROR: %v\n", err)
		return
	}
	return
}
