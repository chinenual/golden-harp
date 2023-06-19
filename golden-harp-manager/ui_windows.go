//go:build windows || !darwin || !linux
// +build windows !darwin !linux

package main

import (
	"fmt"
	"github.com/tadvi/winc"
	"os"
	"strconv"
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

var arduinoStatusLabel *winc.Label

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

var progressBar *winc.ProgressBar

func openProgressBar(context winc.Controller, max int) {
	progressBar = winc.NewProgressBar(context)
	progressBar.SetPos(300, 10)
	progressBar.SetRange(0, max)
	progressBar.SetValue(1)
}
func updateProgressBar(val int) {
	progressBar.SetValue(val)
}
func closeProgressBar() {
	if progressBar != nil {
		progressBar.Close()
	}
	progressBar = nil
}

func drawPresets(list *winc.ListView, presets []Preset, scales []Scale) {
	list.DeleteAllItems()
	for i := 36; i >= 0; i-- {
		var preset Preset
		found := false
		for _, p := range presets {
			if p.KeyPosition == i {
				preset = p
				found = true
				break
			}
		}
		numLabel := fmt.Sprintf("%d", i+1)
		presetLabel := fmt.Sprintf("%s%d", NOTE_NAMES[i%12], (i/12)+1)

		if found {
			leftRootLabel := fmt.Sprintf("%s%d", NOTE_NAMES[preset.Left.Base%12], (preset.Left.Base/12)-2)
			rightRootLabel := fmt.Sprintf("%s%d", NOTE_NAMES[preset.Right.Base%12], (preset.Right.Base/12)-2)

			p1 := &Item{[]string{numLabel, presetLabel,
				leftRootLabel, scaleString(scales[preset.Left.Scale]), fmt.Sprintf("%d", preset.Left.Channel+1),
				rightRootLabel, scaleString(scales[preset.Right.Scale]), fmt.Sprintf("%d", preset.Right.Channel+1)}}
			list.AddItem(p1)

		} else {
			p1 := &Item{[]string{numLabel, presetLabel}}
			list.AddItem(p1)
		}
	}
}

var statusbar *winc.Panel

func setStatus(status string, version string) {
	msg := status
	if version != "" {
		msg = msg + "; Adapter version: " + version
	}
	arduinoStatusLabel.SetText(msg)
}

func showSettingsDialog(context winc.Controller) {
	dlg := winc.NewDialog(context)
	dlg.SetText("Settings")
	dlg.SetSize(350, 200)

	lbl1 := winc.NewLabel(dlg)
	lbl1.SetPos(10, 20)
	lbl1.SetText("Serial Port")

	txt1 := winc.NewEdit(dlg)
	txt1.SetText(userSettings.SerialPort)
	txt1.SetPos(130, 20)

	lbl2 := winc.NewLabel(dlg)
	lbl2.SetPos(10, 50)
	lbl2.SetText("Min Note Length (ms)")

	txt2 := winc.NewEdit(dlg)
	txt2.SetText(strconv.Itoa(userSettings.MinNoteLen))
	txt2.SetPos(130, 50)

	lbl3 := winc.NewLabel(dlg)
	lbl3.SetPos(10, 80)
	lbl3.SetText("Max Note Length (ms)")

	txt3 := winc.NewEdit(dlg)
	txt3.SetText(strconv.Itoa(userSettings.MaxNoteLen))
	txt3.SetPos(130, 80)

	lbl4 := winc.NewLabel(dlg)
	lbl4.SetPos(10, 110)
	lbl4.SetText("Note Resolution (ms)")

	txt4 := winc.NewEdit(dlg)
	txt4.SetText(strconv.Itoa(userSettings.LoopTime))
	txt4.SetPos(130, 110)

	cancelBtn := winc.NewPushButton(dlg)
	cancelBtn.SetPos(100, 140)
	cancelBtn.SetText("Cancel")

	saveBtn := winc.NewPushButton(dlg)
	saveBtn.SetPos(0, 140)
	saveBtn.SetText("Save")
	//	dlg.SetButtons(saveBtn, cancelBtn)

	cancelBtn.OnClick().Bind(func(e *winc.Event) {
		dlg.Close()
	})

	saveBtn.OnClick().Bind(func(e *winc.Event) {
		SerialClose()
		setStatus("Not connected", "")
		userSettings.SerialPort = txt1.Text()
		var err error
		if userSettings.MinNoteLen, err = strconv.Atoi(txt2.Text()); err != nil {
			applog.Printf("ERROR: could not parse Min Note Length as integer: %v\n", err)
			winc.Errorf(dlg, "Error: could not parse Min Note Length as integer: %v", err)
			return
		}
		if userSettings.MaxNoteLen, err = strconv.Atoi(txt3.Text()); err != nil {
			applog.Printf("ERROR: could not parse Max Note Length as integer: %v\n", err)
			winc.Errorf(dlg, "Error: could not parse Max Note Length as integer: %v", err)
			return
		}
		if userSettings.LoopTime, err = strconv.Atoi(txt4.Text()); err != nil {
			applog.Printf("ERROR: could not parse Note Resolution as integer: %v\n", err)
			winc.Errorf(dlg, "Error: could not parse Note Resolution as integer: %v", err)
			return
		}

		defer dlg.Close()
		if err = SaveSettings(); err != nil {
			applog.Printf("ERROR: could not save settings: %v\n", err)
			winc.Errorf(dlg, "Error: could not save settings: %v", err)
			return
		}
	})

	dlg.Show()
}

func showAboutDialog(context winc.Controller) {
	winc.MsgBoxOk(context,
		"About Golden Harp Manager",
		"Version "+Version+"\nCopyright 2023 Steve Tynor (steve.tynor@chinenual.com)")
}

func WindowsUI() {
	winc.SetAppIcon(int(resourceIds["icon_app.ico"]))
	mainWindow := winc.NewForm(nil)
	dock := winc.NewSimpleDock(mainWindow)

	mainWindow.SetSize(1000, 800)
	mainWindow.SetText(AppName)

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	settingsMn := fileMn.AddItem("Settings...", winc.NoShortcut)
	debugMn := fileMn.AddItem("Enable Debug", winc.NoShortcut)
	exitMn := fileMn.AddItem("Exit", winc.NoShortcut)

	helpMn := menu.AddSubMenu("Help")
	aboutMn := helpMn.AddItem("About Golden Harp Manager", winc.NoShortcut)

	aboutMn.OnClick().Bind(func(e *winc.Event) {
		showAboutDialog(mainWindow)
	})
	settingsMn.OnClick().Bind(func(e *winc.Event) {
		showSettingsDialog(mainWindow)
	})
	debugMn.OnClick().Bind(func(e *winc.Event) {
		println("debug-enable click")
		if err := ConnectToArduino(); err != nil {
			applog.Printf("ERROR: could not connect to Arduino: %v\n", err)
			winc.Errorf(mainWindow, "Error: could not connect to Arduino: %v", err)
			return
		}
		if err := CmdSetDebug(true, true, true); err != nil {
			applog.Printf("ERROR: could not enable debug logging: %v\n", err)
			winc.Errorf(mainWindow, "Error: could not enable debug logging: %v", err)
		} else {
			setStatus("Debug Enabled", "")
		}
	})
	exitMn.OnClick().Bind(func(e *winc.Event) {
		mainWindow.Close()
		os.Exit(0)
	})
	menu.Show()

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

	statusbar = winc.NewPanel(mainWindow)
	statusbar.SetSize(0, 25)
	//	runBtn.OnClick().Bind(func(e *winc.Event) {
	//		println("runBtn click")
	//	})

	arduinoStatusLabel = winc.NewLabel(statusbar)
	arduinoStatusLabel.SetPos(10, 5)
	arduinoStatusLabel.SetSize(300, 25)

	setStatus("Not connected", "")

	ls := winc.NewListView(mainWindow)
	ls.AddColumn("Preset", 50)
	ls.AddColumn("Key", 40)
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
	dock.Dock(statusbar, winc.Bottom)

	downloadBtn.OnClick().Bind(func(e *winc.Event) {
		println("downloadBtn click")
		openProgressBar(mainWindow, 4)
		defer closeProgressBar()
		if err := ConnectToArduino(); err != nil {
			applog.Printf("ERROR: could not connect to Arduino: %v\n", err)
			winc.Errorf(mainWindow, "Error: could not connect to Arduino: %v", err)
			return
		}

		updateProgressBar(2)
		if version, _, err := CmdVersion(); err != nil {
			applog.Printf("ERROR: could not get Arduino version info: %v\n", err)
			winc.Errorf(mainWindow, "Error: could not get Arduino version info: %v", err)
		} else {
			setStatus("Connected", version)
		}

		updateProgressBar(3)

		if presets, scales, minNoteLen, maxNoteLen, loopTime, err := CmdGetConfig(); err != nil {
			applog.Printf("ERROR: could not get config from Arduino %v\n", err)
			winc.Errorf(mainWindow, "Error: could not get config from Arduino: %v", err)
			return
		} else {
			applog.Printf("presets: %#v\n", presets)
			applog.Printf("scales: %#v\n", scales)
			applog.Printf("minnotelen: %d maxnotelen: %d looplen: %d\n", minNoteLen, maxNoteLen, loopTime)
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
				applog.Printf("ERROR: could not load config: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not load config file: %v", err)
				return
			}
			applog.Printf("Loaded %s\n", filePath)

			openProgressBar(mainWindow, 7)
			defer closeProgressBar()
			if err := ConnectToArduino(); err != nil {
				applog.Printf("ERROR: could not connect to Arduino: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not connect to Arduino: %v", err)
				return
			}

			updateProgressBar(2)
			if version, _, err := CmdVersion(); err != nil {
				applog.Printf("ERROR: could not get Arduino version info: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not get Arduino version info: %v", err)
			} else {
				setStatus("Connected", version)
			}

			updateProgressBar(3)
			if err := CmdSetTimingParams(userSettings.MinNoteLen, userSettings.MaxNoteLen, userSettings.LoopTime); err != nil {
				applog.Printf("ERROR: could not set timing params: %v\n", err)
				winc.Errorf(mainWindow, "Error: could not set timing params: %v", err)
			}

			updateProgressBar(4)
			for i, _ := range packedScales {
				if err := CmdSetScale(len(packedScales), i, packedScales[i]); err != nil {
					applog.Printf("ERROR: could not get send scale config to Arduino %v\n", err)
					winc.Errorf(mainWindow, "Error: could not get send scale config to Arduino: %v", err)
					return
				}
			}
			updateProgressBar(5)
			for i, _ := range packedPresets {
				if err := CmdSetPreset(len(packedPresets), i, packedPresets[i]); err != nil {
					applog.Printf("ERROR: could not get send preset config to Arduino %v\n", err)
					winc.Errorf(mainWindow, "Error: could not send preset config to Arduino: %v", err)
					return
				}
			}
			updateProgressBar(6)
			if presets, scales, minNoteLen, maxNoteLen, loopTime, err := CmdGetConfig(); err != nil {
				applog.Printf("ERROR: could not get config from Arduino %v\n", err)
				winc.Errorf(mainWindow, "Error: could not get config from Arduino: %v", err)
				return
			} else {
				applog.Printf("presets: %#v\n", presets)
				applog.Printf("scales: %#v\n", scales)
				applog.Printf("minnotelen: %d maxnotelen: %d looplen: %d\n", minNoteLen, maxNoteLen, loopTime)
				drawPresets(ls, presets, scales)
			}

		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop()
}
