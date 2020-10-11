package main

import (
	"fmt"
	"github.com/tadvi/winc/w32"

	"github.com/tadvi/winc"
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
func (item *Item) SetChecked(checked bool) {  }
func (item Item) ImageIndex() int          { return 0 }


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
	name := fmt.Sprintf("%s%d", NOTE_NAMES[index%len(NOTE_NAMES)], index / len(NOTE_NAMES))
	y := 10 + index * CONTROL_HEIGHT;

	label := winc.NewLabel(parent)
	label.SetPos(10, y)
	label.SetText(name)

	octaveL := NewComboBox(parent)
	octaveL.SetPos(40, y)
	octaveL.SetSize(40,10)
	for i := 0; i < 8; i++ {
		octaveL.InsertItem(i, fmt.Sprintf("%d",i))
	}

	keyL := NewComboBox(parent)
	keyL.SetPos(90, y)
	keyL.SetSize(40,10)
	for i := 0; i < len(NOTE_NAMES); i++ {
		keyL.InsertItem(i, NOTE_NAMES[i])
	}

	scaleL := NewComboBox(parent)
	scaleL.SetPos(150, y)
	scaleL.SetSize(200,10)
	// FAKE data
	scaleL.InsertItem(0, "Aeolian")
	scaleL.InsertItem(1, "Major")
	scaleL.InsertItem(2, "Minor")

	octaveR := NewComboBox(parent)
	octaveR.SetPos(440, y)
	octaveR.SetSize(40,10)
	for i := 0; i < 8; i++ {
		octaveR.InsertItem(i, fmt.Sprintf("%d",i))
	}

	keyR := NewComboBox(parent)
	keyR.SetPos(490, y)
	keyR.SetSize(40,10)
	for i := 0; i < len(NOTE_NAMES); i++ {
		keyR.InsertItem(i, NOTE_NAMES[i])
	}

	scaleR := NewComboBox(parent)
	scaleR.SetPos(550, y)
	scaleR.SetSize(200,10)
	// FAKE data
	scaleR.InsertItem(0, "Aeolian")
	scaleR.InsertItem(1, "Major")
	scaleR.InsertItem(2, "Minor")

}

func main() {
	mainWindow := winc.NewForm(nil)
	dock := winc.NewSimpleDock(mainWindow)

	mainWindow.SetSize(1000, 800)
	mainWindow.SetText("Controls Demo")

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

	imlistTb := winc.NewImageList(16, 16)
	imlistTb.AddResIcon(10)
	imlistTb.AddResIcon(12)
	imlistTb.AddResIcon(15)

	// --- Toolbar
	toolbar := winc.NewToolbar(mainWindow)
	toolbar.SetImageList(imlistTb)
	uploadBtn := toolbar.AddButton("Upload", 1)
//	toolbar.AddSeparator()
//	runBtn := toolbar.AddButton("Run Now Fast", 2)
	toolbar.Show()

//	runBtn.OnClick().Bind(func(e *winc.Event) {
//		println("runBtn click")
//	})

	uploadBtn.OnClick().Bind(func(e *winc.Event) {
		println("uploadBtn click")
	})

	// --- Tabs
	tabs := winc.NewTabView(mainWindow)
	presetsPanel := tabs.AddPanel("Presets")
	scalesPanel := tabs.AddPanel("Scales")

	channelLabel := winc.NewLabel(presetsPanel)
	channelLabel.SetPos(800,10)
	channelLabel.SetText("MIDI Channel")

	channelControl := NewComboBox(presetsPanel)
	channelControl.SetPos(880,10)
	channelControl.SetSize(40,10)
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
		p4 := &Item{[]string{fmt.Sprintf("Funky Scale %03d", i+1), "???" }}
		ls.AddItem(p4)
	}

	// --- Dock
	dock2 := winc.NewSimpleDock(scalesPanel)
	dock2.Dock(ls, winc.Fill)
	tabs.SetCurrent(0)

	dock.Dock(toolbar, winc.Top)        // toolbars always dock to the top
	dock.Dock(tabs, winc.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), winc.Fill) // tab panels dock just below tabs and fill area

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop()
}
