// +build windows, !darwin, !linux

package main

import (
	"fmt"
	"github.com/tadvi/winc"
	"log"
)

func btnOnClick(arg *winc.Event) {
	fmt.Println("Button clicked")
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}

type Item struct {
	T []string
}

func (item Item) Text() []string    { return item.T }
func (item *Item) SetText(s string) { item.T[0] = s }

func (item Item) Checked() bool            { return false }
func (item *Item) SetChecked(checked bool) {}
func (item Item) ImageIndex() int          { return 0 }

var resourceIds = map[string]uint16{
	"app.manifest":      1,
	"icon_app.ico":      7,
	"icon_download.ico": 13,
	"icon_upload.ico":   19,
}

var NOTE_NAMES = []string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}
var INTERVAL_NAMES = []string{"1", "b2", "2", "b3", "3", "4", "b5", "5", "b6", "6", "b7", "7"}

func scaleString(scale Scale) (result string) {
	for _, interval := range scale.Intervals {
		result += "  " + INTERVAL_NAMES[interval]
	}
	// Iasos's convention always appends the octave at the end of the scale:
	result += "  8"
	return
}

func drawPresets(list *winc.ListView, presets []Preset, scales []Scale) {
	list.DeleteAllItems()
	for i := 0; i < 37; i++ {
		var preset Preset
		found := false
		for _, p := range presets {
			if p.KeyPosition == i {
				preset = p
				found = true
				break
			}
		}
		presetLabel := fmt.Sprintf("%s%d", NOTE_NAMES[i%12], (i/12)+1)

		if found {
			leftRootLabel := fmt.Sprintf("%s%d", NOTE_NAMES[preset.Left.Base%12], (preset.Left.Base/12)-2)
			rightRootLabel := fmt.Sprintf("%s%d", NOTE_NAMES[preset.Right.Base%12], (preset.Right.Base/12)-2)

			p1 := &Item{[]string{presetLabel,
				leftRootLabel, scaleString(scales[preset.Left.Scale]), fmt.Sprintf("%d", preset.Left.Channel+1),
				rightRootLabel, scaleString(scales[preset.Right.Scale]), fmt.Sprintf("%d", preset.Right.Channel+1)}}
			list.AddItem(p1)

		} else {
			p1 := &Item{[]string{presetLabel}}
			list.AddItem(p1)
		}
	}
}

func WindowsUI() {
	winc.SetAppIcon(int(resourceIds["icon_app.ico"]))
	mainWindow := winc.NewForm(nil)
	dock := winc.NewSimpleDock(mainWindow)

	mainWindow.SetSize(1000, 800)
	mainWindow.SetText("Golden Harp Manager")

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	fileMn.AddItem("New", winc.NoShortcut)
	editMn := menu.AddSubMenu("Edit")
	cutMn := editMn.AddItem("Cut", winc.Shortcut{winc.ModControl, winc.KeyX})
	copyMn := editMn.AddItem("Copy", winc.NoShortcut)
	pasteMn := editMn.AddItem("Paste", winc.NoShortcut)
	menu.Show()
	copyMn.SetCheckable(true)
	copyMn.SetChecked(true)
	pasteMn.SetEnabled(false)

	cutMn.OnClick().Bind(func(e *winc.Event) {
		println("cut click")
	})

	imlistTb := winc.NewImageList(24, 24)
	imlistTb.AddResIcon(resourceIds["icon_download.ico"])
	imlistTb.AddResIcon(resourceIds["icon_upload.ico"])

	// --- Toolbar
	toolbar := winc.NewToolbar(mainWindow)
	toolbar.SetImageList(imlistTb)
	downloadBtn := toolbar.AddButton("Connect", 0)
	uploadBtn := toolbar.AddButton("Upload", 1)
	//	toolbar.AddSeparator()
	//	runBtn := toolbar.AddButton("Run Now Fast", 2)
	toolbar.Show()

	//	runBtn.OnClick().Bind(func(e *winc.Event) {
	//		println("runBtn click")
	//	})

	ls := winc.NewListView(mainWindow)
	ls.AddColumn("Preset", 60)
	ls.AddColumn("L Base", 60)
	ls.AddColumn("L Intervals", 200)
	ls.AddColumn("L Channel", 60)
	ls.AddColumn("R Base", 60)
	ls.AddColumn("R Intervals", 200)
	ls.AddColumn("R Channel", 60)

	drawPresets(ls, nil, nil)

	dock.Dock(toolbar, winc.Top) // toolbars always dock to the top
	//dock.Dock(tabs, winc.Top)           // tabs should prefer docking at the top
	//dock.Dock(tabs.Panels(), winc.Fill) // tab panels dock just below tabs and fill area
	dock.Dock(ls, winc.Fill)

	downloadBtn.OnClick().Bind(func(e *winc.Event) {
		println("downloadBtn click")
		if err := ConnectToArduino(); err != nil {
			log.Printf("ERROR: could not connect to Arduino: %v\n", err)
			winc.Errorf(mainWindow, "Error: could not connect to Arduino: %v", err)
			return
		}

		if presets, scales, err := CmdGetConfig(); err != nil {
			log.Printf("ERROR: could not get config from Arduino %v\n", err)
			winc.Errorf(mainWindow, "Error: could not get config from Arduino: %v", err)
			return
		} else {
			log.Printf("presets: %#v\n", presets)
			log.Printf("scales: %#v\n", scales)
			drawPresets(ls, presets, scales)
		}
		return
	})
	uploadBtn.OnClick().Bind(func(e *winc.Event) {
		println("uploadBtn click")
		if filePath, ok := winc.ShowOpenFileDlg(mainWindow,
			"Select Harp Scale/Preset configuration file",
			"Config files (*.xlsx)|*.xlsx|All files (*.*)|*.*",
			0, ""); ok {
			if err := LoadConfig(filePath); err != nil {
				log.Printf("ERROR: could not load config: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not load config file: %v", err)
				return
			}
			log.Printf("Loaded %s\n", filePath)
			if err := ConnectToArduino(); err != nil {
				log.Printf("ERROR: could not connect to Arduino: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not connect to Arduino: %v", err)
				return
			}

			for i, _ := range packedScales {
				if err := CmdSetScale(len(packedScales), i, packedScales[i]); err != nil {
					log.Printf("ERROR: could not get send scale config to Arduino %v\n", err)
					winc.Errorf(mainWindow, "Error: could not get send scale config to Arduino: %v", err)
					return
				}
			}
			for i, _ := range packedPresets {
				if err := CmdSetPreset(len(packedPresets), i, packedPresets[i]); err != nil {
					log.Printf("ERROR: could not get send preset config to Arduino %v\n", err)
					winc.Errorf(mainWindow, "Error: could not get send preset config to Arduino: %v", err)
					return
				}
			}
			if presets, scales, err := CmdGetConfig(); err != nil {
				log.Printf("ERROR: could not get config from Arduino %v\n", err)
				winc.Errorf(mainWindow, "Error: could not get config from Arduino: %v", err)
				return
			} else {
				log.Printf("presets: %#v\n", presets)
				log.Printf("scales: %#v\n", scales)
				drawPresets(ls, presets, scales)
			}

		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop()
}
