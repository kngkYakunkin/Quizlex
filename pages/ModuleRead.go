package pages

import (
	"fmt"
	"image/color"
	"log"
	"quizlex/db"
	"quizlex/layouts"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ModuleRead(window fyne.Window, moduleName string, mainPage func(window fyne.Window, searcingResults []string)) {

	moduleNameLabel := canvas.NewText(moduleName, theme.ForegroundColor())
	moduleNameLabel.TextSize = 24
	moduleNameLabel.TextStyle.Bold = true
	moduleNameLabel.Alignment = fyne.TextAlignCenter

	data := [][2]string{}

	db, _ := db.ConnectDB()

	rows, err := db.Query("SELECT * FROM modules WHERE name = $1", moduleName)
	if err != nil {
		log.Fatal(err)
	}

	var moduleId int = 0
	for rows.Next() {
		var id int
		var name string
		var createdAt time.Time

		err := rows.Scan(&id, &name, &createdAt)
		if err != nil {
			log.Fatal(err)
		}

		moduleId = id
	}
	defer rows.Close()

	words, _ := db.Query(`SELECT * FROM words WHERE module_id = $1`, moduleId)
	fmt.Println(moduleId)
	for words.Next() {
		var id int
		var term string
		var definition string
		var module_id int

		err := words.Scan(&id, &term, &definition, &module_id)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, [2]string{term, definition})
	}

	defer words.Close()

	backButton := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		mainPage(window, []string{})
	})

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})

	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	errorLabel.Alignment = fyne.TextAlignCenter

	learnButton := widget.NewButtonWithIcon("Learn", theme.FileTextIcon(), func() {
		if len(data) >= 4 {
			pastData := make([][2]string, len(data))
			copy(pastData, data)
			ModuleLearn(window, moduleName, data, pastData, []string{}, mainPage)
		} else {
			errorLabel.Text = "Количество слов в модуле должно быть не менее 4!"
		}
	})

	learnButton.Importance = widget.SuccessImportance

	// Canvas objects for terms and definitions
	var terms []fyne.CanvasObject

	// Generating $terms
	for i := 0; i < len(data); i++ {

		newDefElement := canvas.NewText((data)[i][1], theme.ForegroundColor())
		newDefElement.TextSize = 18
		newDefElement.Alignment = fyne.TextAlignCenter

		newTermElement := canvas.NewText(data[i][0], theme.ForegroundColor())
		newTermElement.TextSize = 18
		newTermElement.Alignment = fyne.TextAlignCenter

		// Separator line
		separator := canvas.NewLine(color.White)
		separator.StrokeWidth = 0.9

		terms = append(terms, fyne.NewContainer())
		terms[i] = container.NewVBox(layouts.NewAdaptiveGridWithRatios([]float32{0.5, 0.5}, newTermElement, newDefElement), layouts.NewAdaptiveGridWithRatios([]float32{1}, separator))
	}

	window.SetContent(container.NewVBox(
		// Module name
		layouts.NewAdaptiveGridWithRatios([]float32{0.05, 0.9, 0.05}, backButton, moduleNameLabel, settingsButton),
		widget.NewLabel(""),
		// Titles
		layouts.NewAdaptiveGridWithRatios([]float32{0.5, 0.5}, setTitle("Terms"), setTitle("Definitions")),
		widget.NewLabel(""),
		// Words labels
		container.NewVBox(terms...),
		layout.NewSpacer(),
		// Back button
		errorLabel,
		layouts.NewAdaptiveGridWithRatios([]float32{0.4, 0.2, 0.4}, layout.NewSpacer(), learnButton, layout.NewSpacer()),
	))
}

func setTitle(text string) *canvas.Text {
	title := canvas.NewText(text, theme.ForegroundColor())
	title.TextSize = 17
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	return title
}
