package main

import (
	"encoding/json"
	"log"
)

func CmdVersion() (version string, err error) {
	if err = SerialWriteCommand([]byte("{cmd: \"version\"}")); err != nil {
		return
	}
	var bytes []byte
	if bytes, err = SerialReadResponse(); err != nil {
		return
	}
	log.Printf(" version: %s\n", string(bytes))

	var data map[string]interface{}

	if err = json.Unmarshal([]byte(bytes), &data); err != nil {
		return
	}

	version = data["version"].(string)
	return
}