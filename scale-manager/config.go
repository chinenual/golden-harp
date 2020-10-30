package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

var scales []string
var intervals [][]int

func main() {
	var f *excelize.File
	var err error
	if f, err = excelize.OpenFile("HarpConfig.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	readScales(f)
	readPresets(f)

	fmt.Printf("INTERVALS: %v\n", intervals)
}

func readScales(f *excelize.File) {
	rows, err := f.Rows("Scales")
	if err != nil {
		fmt.Println(err)
		return
	}
	var count = 0
	for rows.Next() {
		// skip the header rows
		colVals, err := rows.Columns()
		if count >= 1 {
			if err != nil {
				fmt.Println(err)
			}
			if strings.TrimSpace(colVals[0]) == "" {
				break
			}
			fmt.Printf("name: %s, intervals: %s\n", colVals[0], colVals[1])
			scales = append(scales, colVals[0])
			var newscale []int
			if newscale, err = parseIntervals(colVals[1]); err != nil {
				fmt.Printf("ERROR: Row %d: %v\n",count+1,err)
				//return
			}
			intervals = append(intervals, newscale)
		}
		count = count + 1
	}
}

func readPresets(f *excelize.File) {
	rows, err := f.Rows("Presets")
	if err != nil {
		fmt.Println(err)
		return
	}
	var count = 0
	for rows.Next() {
		// skip the header rows
		colVals, err := rows.Columns()
		if count >= 2 {
			if err != nil {
				fmt.Println(err)
			}
			if strings.TrimSpace(colVals[2]) == "" {
				break
			}
			fmt.Printf("left: %s, base: %s, chan: %s   right: %s, base: %s, chan: %s\n",
				colVals[2], colVals[3], colVals[4],
				colVals[6], colVals[7], colVals[8])
		}
		count = count + 1
	}
}

func parseIntervals(s string) (result []int, err error) {
	re := regexp.MustCompile("^[b#]?(\\d+)[b#]?$")
	raw := strings.Fields(s)
	for _, r := range raw {
		var offset = 0
		var valstring []string
		if valstring = re.FindStringSubmatch(r); valstring == nil {
			err = errors.Errorf("Invalid interval %s - bad syntax. should be a number with an optional # or b suffix/prefix", r)
			return
		}
		if strings.HasPrefix(r, "#") || strings.HasSuffix(r, "#") {
			offset = 1
		} else if strings.HasPrefix(r, "b") || strings.HasSuffix(r, "b") {
			offset = -1
		}
		val, _ := strconv.Atoi(valstring[1])
		// if raw was "b3", val is now "3"
		var majorIntervals = []int{-1, 0, 2, 4, 5, 7, 9, 11, 12};
		if val < 0 || val >= len(majorIntervals) {
			err = errors.Errorf("Invalid interval %s - out of range (1 .. 8)", r)
			return
		}
		val = majorIntervals[val];
		// val is now the semitone offset in the octave for the 3rd degree of the major scale
		val = val + offset
		// val is now flattened or sharped if it had a b or # suffix/prefix.
		if val < 0 || val > 12 {
			err = errors.Errorf("Invalid interval %s - out of range", r)
			return
		}
		if val < 12 {
			// 8 (i.e. the octave - is acceptable in the config file, but its unused in the interval processing
			result = append(result, val)
		}
	}
	return
}
