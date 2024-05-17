package pages

import (
	"log"
	"quizlex/db"
	"quizlex/layouts"
	"quizlex/widgets"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ModuleCreatingPage(window fyne.Window, mainPage func(window fyne.Window, searcingResults []string)) {

	var page_title string = "Create Module"
	window.SetTitle(page_title)

	// Array with words
	data := [][2]string{{"", ""}}

	var focusId int = 0
	render(window, data, focusId, "")

}

func render(window fyne.Window, data [][2]string, focusId int, moduleNameValue string) {
	moduleNameInput := widget.NewEntry()
	moduleNameInput.Text = moduleNameValue
	moduleNameInput.OnChanged = func(s string) {
		moduleNameValue = s
	}

	var content_width float32 = 0.8
	center_pos := []float32{(1 - content_width) / 2, content_width, (1 - content_width) / 2}

	backButton := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {

	})

	importButton := widget.NewButtonWithIcon("Import", theme.FileTextIcon(), func() {
		myCustomDialog := widgets.NewImportDialog(window, "hello")
		myCustomDialog.Show()
	})

	// Canvas objects for terms and definitions
	var terms []fyne.CanvasObject

	// Needs for finding an element with the focus
	focus_element := widgets.NewDefEntry(window, &data, len(terms), render, focusId, &moduleNameValue)

	// Generating $terms
	for i := 0; i < len(data); i++ {

		newDefEntryElement := widgets.NewDefEntry(window, &data, len(terms), render, focusId, &moduleNameValue)
		newDefEntryElement.Text = (data)[i][1]

		newTermEntryElement := widget.NewEntry()
		newTermEntryElement.Text = (data)[i][0]
		newTermEntryElement.OnChanged = func(text string) {
			(data)[i][0] = text
		}

		removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			if len(data) > 1 {
				data = append(data[:i], data[i+1:]...)
				render(window, data, focusId, moduleNameValue)
			}
		})

		terms = append(terms, fyne.NewContainer())
		terms[i] = container.NewVBox(layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.4, 0.4, 0.05, 0.05}, layout.NewSpacer(), newTermEntryElement, newDefEntryElement, removeButton, layout.NewSpacer()))
		if i == focusId {
			focus_element = newDefEntryElement
		}
	}

	saveModuleButton := widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
		// Validation
		var isError = false

		for i := 0; i < len(data); i++ {
			if data[i][0] == "" || data[i][1] == "" {
				isError = true
				break
			} else {
				for t := 0; t < len(data); t++ {
					if t != i && data[t][0] == data[i][0] {
						isError = true
						break
					}
				}
			}
		}
		if moduleNameInput.Text == "" {
			isError = true
		}

		if !isError {

			db, _ := db.ConnectDB()
			var moduleId int
			err := db.QueryRow(`INSERT INTO modules (name, created_at) VALUES ($1, $2) RETURNING id`, moduleNameInput.Text, time.Now()).Scan(&moduleId)

			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < len(data); i++ {
				db.Exec(`INSERT INTO words (term, definition, module_id) VALUES ($1, $2, $3)`, data[i][0], data[i][1], moduleId)
			}

		}
	})

	saveModuleButton.Importance = widget.SuccessImportance

	window.SetContent(container.NewVBox(
		// Back button
		layouts.NewAdaptiveGridWithRatios([]float32{0.05, 0.1}, backButton, importButton),
		// Name input
		layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.4}, layout.NewSpacer(), widget.NewLabel("Module Name")),
		// Labels
		layouts.NewAdaptiveGridWithRatios(center_pos, layout.NewSpacer(), moduleNameInput, layout.NewSpacer()),
		layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.2, 0.5, 0.2}, layout.NewSpacer(), widget.NewLabel("Term"), layout.NewSpacer(), widget.NewLabel("Definition")),
		// Words entries
		container.NewVBox(terms...),
		layout.NewSpacer(),
		layouts.NewAdaptiveGridWithRatios([]float32{0.45, 0.1, 0.45}, layout.NewSpacer(), saveModuleButton, layout.NewSpacer()),
	))

	// Set focus
	window.Canvas().Focus(focus_element)
	if len(data) > 1 {
		window.Canvas().FocusNext()
		window.Canvas().FocusNext()
	} else {
		window.Canvas().FocusPrevious()
		window.Canvas().FocusPrevious()
	}

}
