package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"io"
	"time"
)

var bufreader *bufio.Reader
var unbuffered io.ReadWriteCloser
var inputChan chan []byte

func SerialConnected() bool {
	return bufreader != nil
}

func SerialInit(port string, baudRate uint) (err error) {
	options := serial.OpenOptions{
		PortName:              port,
		BaudRate:              baudRate,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 500,
		MinimumReadSize:       1,
		DataBits:              8,
		StopBits:              1,
	}
	if unbuffered, err = serial.Open(options); err != nil {
		err = errors.Wrapf(err, "Could not open serial port %s", port)
		return
	}
	rdr := bufio.NewReader(unbuffered)
	bufreader = bufio.NewReader(rdr)

	inputChan = make(chan []byte)
	go readInput()

	return
}
func SerialClose() (err error) {
	if bufreader != nil {
		if err = unbuffered.Close(); err != nil {
			return
		}
		bufreader = nil
	}
	return
}

func writeLine(bytes []byte) (err error) {
	applog.Printf("SEND: \"%s\"...\n", string(bytes))
	if _, err = unbuffered.Write(bytes); err != nil {
		return
	}
	//	if _,err = unbuffered.Write([]byte{'\n'}); err != nil {
	//		return
	//	}
	return
}

// consumes the input - if payload is a DEBUG string, handles it directly
// else assumed to be a command response and puts it on the channel.
// This allows debug strings to come before and after commands without
// the command processor needing to wait for them
func readInput() {
	for {
		fmt.Println("readinput 1")
		bytes, err := readLine()
		fmt.Println("readinput 2")
		// already logged of err != nil
		if err == nil {
			fmt.Println("readinput 3")
			var data map[string]interface{}
			if err = json.Unmarshal([]byte(bytes), &data); err != nil {
				continue
			}
			fmt.Println("readinput 4 " + string(bytes))
			if data["DEBUG"] != nil {
				str := string(bytes)
				fmt.Println("readinput 5")
				applog.Printf("DEBUG: \"%s\"", str)
				DisplayDebug(str)
			} else {
				fmt.Println("readinput 6")
				inputChan <- bytes
				// let the other coroutine wake up
				fmt.Println("readinput 6.1")
				time.Sleep(2000 * time.Millisecond)
				fmt.Println("readinput 6.2")
			}
			fmt.Println("readinput 7")
		}
		fmt.Println("readinput 8")
	}
}

func readLine() (bytes []byte, err error) {
	if bytes, err = bufreader.ReadBytes('\n'); err != nil {
		applog.Println("err")
		return
	}
	return
}

func SerialWriteCommand(json []byte) (err error) {
	if err = writeLine(json); err != nil {
		return
	}
	return
}

func SerialReadResponse() (json []byte, err error) {
	fmt.Println("SerialReadResponse 1")
	json = <-inputChan
	fmt.Println("SerialReadResponse 2")
	applog.Printf("READ: \"%s\"", string(json))
	return
}
