package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

var scaleMap map[string]Scale
var packedScales []Scale
var packedPresets []Preset

type Scale struct {
	Name      string `json:"-"`
	Intervals []int  `json:"i"`
}

type StripPreset struct {
	Scale   int `json:"scale"` // packed index
	Base    int `json:"base"`
	Channel int `json:"chan"`
}
type Preset struct {
	KeyPosition int         `json:"key"`
	Left        StripPreset `json:"l"`
	Right       StripPreset `json:"r"`
}

func LoadConfig(filename string) (err error) {
	scaleMap = make(map[string]Scale)
	var f *excelize.File
	if f, err = excelize.OpenFile(filename); err != nil {
		return
	}

	if err = readScales(f); err != nil {
		return
	}
	if err = readPresets(f); err != nil {
		return
	}

	//fmt.Printf("Scales: %v\n", scaleMap)
	//fmt.Printf("PackedScales: %#v\n", packedScales)
	//fmt.Printf("PackedPresets: %#v\n", packedPresets)
	return
}

func readScales(f *excelize.File) (err error) {
	rows, err := f.Rows("Scales")
	if err != nil {
		fmt.Println(err)
		return
	}
	var count = 0
	var blankRowCount = 0
	for rows.Next() {
		// skip the header rows
		var colVals []string
		if colVals, err = rows.Columns(); err != nil {
			return err
		}
		if count >= 1 {
			if len(colVals) > 12 {
				if strings.TrimSpace(colVals[0]) == "" {
					blankRowCount++
					if blankRowCount > 3 {
						break
					}
				} else {
					blankRowCount = 0
					var scale Scale
					scale.Name = colVals[0]
					for i := 0; i < 12; i++ {
						val := colVals[i+1]
						if strings.TrimSpace(val) != "" {
							scale.Intervals = append(scale.Intervals, i)
						}
					}
					scaleMap[scale.Name] = scale
				}
			}
		}
		count = count + 1
	}
	return
}

func readPresets(f *excelize.File) (err error) {
	rows, err := f.Rows("Presets")
	if err != nil {
		fmt.Println(err)
		return
	}
	var rowNum = 0
	var keyPosition = 0
	for rows.Next() {
		// skip the header rows
		var colVals []string
		if colVals, err = rows.Columns(); err != nil {
			return err
		}
		fmt.Println(colVals)
		if len(colVals) > 0 {
			if rowNum >= 2 { // first two rows are headers
				if strings.TrimSpace(colVals[0]) == "" {
					break
				}
				keyPosition, _ = strconv.Atoi(colVals[0])
				keyPosition -= 1 // one based in the spreadsheet; convert to zero-based

				if len(colVals) >= 8 {
					if strings.TrimSpace(colVals[2]) != "" {
						var preset Preset
						preset.KeyPosition = keyPosition
						if preset.Left.Scale, err = useScale(colVals[2]); err != nil {
							return err
						}
						if preset.Left.Base, err = parseRoot(colVals[3]); err != nil {
							return err
						}
						if preset.Left.Channel, err = parseChannel(colVals[4]); err != nil {
							return err
						}
						if preset.Right.Scale, err = useScale(colVals[6]); err != nil {
							return err
						}
						if preset.Right.Base, err = parseRoot(colVals[7]); err != nil {
							return err
						}
						if len(colVals) > 8 { // channel may be blank and the reader won't include that in the columns
							if preset.Right.Channel, err = parseChannel(colVals[8]); err != nil {
								return err
							}
						} else {
							preset.Right.Channel = 0
						}
						packedPresets = append(packedPresets, preset)
					}
				}
			}
		}
		rowNum = rowNum + 1
	}
	return
}

func useScale(name string) (packedIndex int, err error) {
	scale, ok := scaleMap[name]
	if !ok {
		err = errors.Errorf("Reference to unknown scale \"%s\"", name)
	}
	for idx, packedScale := range packedScales {
		if scale.Name == packedScale.Name {
			packedIndex = idx
			return
		}
	}
	packedScales = append(packedScales, scale)
	packedIndex = len(packedScales) - 1
	return
}

func parseRoot(offsetName string) (root int, err error) {
	// Use Yamaha conventions: C3 == Middle-C == MIDI 60

	re := regexp.MustCompile("^([A-G])([b#]?)([-]?\\d+)$")
	trimmed := strings.TrimSpace(offsetName)
	var valstring []string
	if valstring = re.FindStringSubmatch(trimmed); valstring == nil {
		err = errors.Errorf("Invalid transpose root %s - bad syntax. should be a note name with an optional # or b suffix and required octave (e.g. C#2 or Bf-1)", trimmed)
		return
	}
	// valstring[0] == the overall match
	// valstring[1] == the note name
	// valstring[2] == optional sharp or flat
	// valstring[3] == octave
	var m = map[string]int{
		"C": 0,
		"D": 2,
		"E": 4,
		"F": 5,
		"G": 7,
		"A": 9,
		"B": 11,
	}
	base := m[valstring[1]]
	sharpflat := 0
	if valstring[2] == "#" {
		sharpflat = 1
	} else if valstring[2] == "b" {
		sharpflat = -1
	}
	octave, _ := strconv.Atoi(valstring[3])
	root = base + ((octave + 2) * 12) + sharpflat
	if root < 0 || root > 127 {
		err = errors.Errorf("Invalid transpose root %s - must be in range C-2 .. G8", trimmed)
	}
	return
}

// string is 1-based, result is 0-based
func parseChannel(channelString string) (channel int, err error) {
	if strings.TrimSpace(channelString) == "" {
		channel = 0
		return
	}
	if channel, err = strconv.Atoi(channelString); err != nil {
		return
	}
	if channel < 1 || channel > 16 {
		err = errors.Errorf("Channel %s out of range - must be 1 .. 16", channelString)
	}
	channel = channel - 1
	return
}
