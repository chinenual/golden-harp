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

/*
// same as the winc.NewComboBox except uses CBS_DROPDOWN instead of CBS_DROPDOWNLIST style
func NewComboBox(parent winc.Controller) *winc.ComboBox {
	if true {
		return winc.NewComboBox(parent)
	}

	cb := new(winc.ComboBox)

	cb.InitControl("COMBOBOX", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_DROPDOWN)
	winc.RegMsgHandler(cb)

	cb.SetFont(winc.DefaultFont)
	cb.SetSize(200, 400)
	return cb
}

func makePresetControl(index int, parent *winc.Panel) {
	CONTROL_HEIGHT := 20
	NOTE_NAMES := []string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}
	name := fmt.Sprintf("%s%d", NOTE_NAMES[index%len(NOTE_NAMES)], index/len(NOTE_NAMES))
	y := 10 + index*CONTROL_HEIGHT

	label := winc.NewLabel(parent)
	label.SetPos(10, y)
	label.SetText(name)

	octaveL := NewComboBox(parent)
	octaveL.SetPos(40, y)
	octaveL.SetSize(40, 10)
	for i := 0; i < 8; i++ {
		octaveL.InsertItem(i, fmt.Sprintf("%d", i))
	}

	keyL := NewComboBox(parent)
	keyL.SetPos(90, y)
	keyL.SetSize(40, 10)
	for i := 0; i < len(NOTE_NAMES); i++ {
		keyL.InsertItem(i, NOTE_NAMES[i])
	}

	scaleL := NewComboBox(parent)
	scaleL.SetPos(150, y)
	scaleL.SetSize(200, 10)
	// FAKE data
	scaleL.InsertItem(0, "Aeolian")
	scaleL.InsertItem(1, "Major")
	scaleL.InsertItem(2, "Minor")

	octaveR := NewComboBox(parent)
	octaveR.SetPos(440, y)
	octaveR.SetSize(40, 10)
	for i := 0; i < 8; i++ {
		octaveR.InsertItem(i, fmt.Sprintf("%d", i))
	}

	keyR := NewComboBox(parent)
	keyR.SetPos(490, y)
	keyR.SetSize(40, 10)
	for i := 0; i < len(NOTE_NAMES); i++ {
		keyR.InsertItem(i, NOTE_NAMES[i])
	}

	scaleR := NewComboBox(parent)
	scaleR.SetPos(550, y)
	scaleR.SetSize(200, 10)
	// FAKE data
	scaleR.InsertItem(0, "Aeolian")
	scaleR.InsertItem(1, "Major")
	scaleR.InsertItem(2, "Minor")

}
*/

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
		presetLabel := fmt.Sprintf("%s%d", NOTE_NAMES[i%12], i/12)

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

	/*
		// --- Tabs
		tabs := winc.NewTabView(mainWindow)
		presetsPanel := tabs.AddPanel("Presets")
		scalesPanel := tabs.AddPanel("Scales")

		channelLabel := winc.NewLabel(presetsPanel)
		channelLabel.SetPos(800, 10)
		channelLabel.SetText("MIDI Channel")

		channelControl := NewComboBox(presetsPanel)
		channelControl.SetPos(880, 10)
		channelControl.SetSize(40, 10)
		for i := 0; i < 16; i++ {
			channelControl.InsertItem(i, fmt.Sprintf("%d", i+1))
		}
		channelControl.SetSelectedItem(0)
		for i := 0; i < 36; i++ {
			makePresetControl(i, presetsPanel)
		}
		//imlist := winc.NewImageList(16, 16)
		//imlist.AddResIcon(12)

		ls := winc.NewListView(scalesPanel)
		//ls.SetImageList(imlist)
		//ls.EnableEditLabels(false)
		//ls.SetCheckBoxes(true)
		ls.EnableFullRowSelect(true)
		//ls.EnableHotTrack(true)
		ls.EnableEditLabels(true)

		//	ls.EnableSortHeader(true)
		//	ls.EnableSortAscending(true)

		ls.AddColumn("Name", 120)
		ls.AddColumn("Intervals", 120)
		ls.SetPos(10, 180)
		p1 := &Item{[]string{"Minor Pentatonic", "1 b3 4 5 b7"}}
		ls.AddItem(p1)
		p2 := &Item{[]string{"Aeolian", "1 2 b3 4 5 b6 b7"}}
		ls.AddItem(p2)
		p3 := &Item{[]string{"Dorian", "1 2 b3 4 5 6 b7"}}
		ls.AddItem(p3)
		for i := 0; i < 20; i++ {
			p4 := &Item{[]string{fmt.Sprintf("Funky Scale %03d", i+1), "???"}}
			ls.AddItem(p4)
		}

		// --- Dock
		dock2 := winc.NewSimpleDock(scalesPanel)
		dock2.Dock(ls, winc.Fill)
		tabs.SetCurrent(0)
	*/
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
