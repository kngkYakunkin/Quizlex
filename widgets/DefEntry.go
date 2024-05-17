package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type defEntry struct {
	widget.Entry
	window          fyne.Window
	id              int
	data            *[][2]string
	render          func(fyne.Window, [][2]string, int, string)
	focusId         *int
	moduleNameValue *string
}

func NewDefEntry(window fyne.Window, data *[][2]string, id int, render func(fyne.Window, [][2]string, int, string), focusId int, moduleNameValue *string) *defEntry {
	entry := &defEntry{window: window, data: data, id: id, render: render, focusId: &focusId, moduleNameValue: moduleNameValue}
	entry.ExtendBaseWidget(entry)

	entry.OnChanged = func(text string) {
		(*data)[id][1] = text
	}

	return entry
}

func addWord(id int, data *[][2]string) {
	// Position to put an element
	var putTo int = id + 1
	var temp_line [2]string
	var temp_line_b [2]string

	if putTo == len(*data) {
		*data = append(*data, [2]string{"", ""})
	} else {
		for i := 0; i < len(*data); i++ {

			// continue if i < putTo

			if i == putTo {
				temp_line = (*data)[i]
				(*data)[i] = [2]string{"", ""}
				*data = append(*data, [2]string{"", ""})
			} else if i > putTo {
				temp_line_b = (*data)[i]
				(*data)[i] = temp_line
				temp_line = temp_line_b
			}
		}
	}

}

func (e *defEntry) KeyDown(key *fyne.KeyEvent) {
	// NEED TO IMPROVE
	if key.Name == fyne.KeyReturn {
		addWord(e.id, e.data)
		*e.focusId = e.id
		e.render(e.window, *e.data, *e.focusId, *e.moduleNameValue)
	} else if key.Name == fyne.KeyTab && e.id == len(*e.data)-1 {
		addWord(e.id, e.data)
		*e.focusId = e.id
		e.render(e.window, *e.data, *e.focusId, *e.moduleNameValue)
	}

}
