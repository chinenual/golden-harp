package main

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"io"
)

var bufreader *bufio.Reader
var unbuffered io.ReadWriteCloser

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
	for {
		if json, err = readLine(); err != nil {
			return
		}

		// Arduino writes "debug" msgs with a leading #
		if len(json) > 0 && json[0] != '#' {
			applog.Printf("READ: \"%s\"", string(json))
			return
		}
		// includes a newline so don't include it on the printf
		applog.Printf("DEBUG: \"%s\"", string(json))
	}
	return
}
