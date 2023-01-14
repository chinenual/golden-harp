//go:build darwin || !windows || !linux
// +build darwin !windows !linux

package main

import "fmt"

func WindowsInit() {
	// do nothing function to allow command-line only variant to build on Mac
}

func WindowsUI() {
	// do nothing function to allow command-line only variant to build on Mac
	fmt.Printf("Had this been on Windows, would have launched the UI\n")
}

func DisplayDebug(str string) {
	// nop for non-Windows
}
