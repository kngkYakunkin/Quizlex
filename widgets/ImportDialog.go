package widgets

import (
	"fmt"
	"image/color"
	"log"
	"quizlex/db"
	"quizlex/layouts"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ImportDialog struct {
	dialog *dialog.CustomDialog
	window fyne.Window
}

func NewImportDialog(window fyne.Window, content string) *ImportDialog {
	var importDialogTitle string = "Importing"
	var importExitButtonTitle string = "Cancel"
	// Dialog width and height
	var dialog_w float32 = 400
	var dialog_h float32 = 700

	//var descriptionText string = "You can export words from quizlet or other resources."

	// Dialog description
	descriptionTextElement := canvas.NewText("You can export words from quizlet or other resources.", theme.ForegroundColor())
	descriptionTextElement.Alignment = fyne.TextAlignCenter

	descriptionDownTextElement := canvas.NewText("Just paste them here!", theme.ForegroundColor())
	descriptionDownTextElement.Alignment = fyne.TextAlignCenter
	descriptionDownTextElement.TextSize = 18
	descriptionDownTextElement.Color = color.RGBA{0, 128, 0, 255}

	// Importing entry
	importingArea := widget.NewEntry()
	importingArea.MultiLine = true

	moduleNameLabel := canvas.NewText("Module Name", theme.ForegroundColor())
	moduleNameLabel.Alignment = fyne.TextAlignCenter
	moduleNameLabel.Color = color.RGBA{0, 128, 0, 255}

	moduleNameEntry := widget.NewEntry()

	saveButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		if moduleNameEntry.Text != "" || importingArea.Text != "" {
			saveWords(importingArea.Text, moduleNameEntry.Text)
		}
	})
	saveButton.Importance = widget.SuccessImportance

	dismissButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), nil)
	dismissButton.Importance = widget.DangerImportance

	// Content
	dialogContent := container.NewBorder(
		container.NewVBox(descriptionTextElement, moduleNameLabel, moduleNameEntry, descriptionDownTextElement),
		container.NewVBox(
			layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.4, 0.4, 0.1}, layout.NewSpacer(), dismissButton, saveButton, layout.NewSpacer()),
		),
		nil,
		nil,
		container.NewMax(importingArea),
	)

	dialog := dialog.NewCustom(importDialogTitle, importExitButtonTitle, dialogContent, window)

	dialog.SetButtons([]fyne.CanvasObject{})
	dismissButton.OnTapped = func() {
		dialog.Hide()
	}

	importDialog := &ImportDialog{
		window: window,
		dialog: dialog,
	}

	dialog.Resize(fyne.NewSize(dialog_w, dialog_h))

	return importDialog
}

func (d *ImportDialog) Show() {
	d.dialog.Show()
}

func saveWords(content string, moduleName string) {
	var some []string = strings.Split(content, "\n")
	data := [][2]string{}

	for _, oneString := range some {
		stringParts := strings.Split(oneString, "	")
		data = append(data, [2]string{stringParts[0], stringParts[1]})
	}

	db, _ := db.ConnectDB()
	var moduleId int
	err := db.QueryRow(`INSERT INTO modules (name, created_at) VALUES ($1, $2) RETURNING id`, moduleName, time.Now()).Scan(&moduleId)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(moduleId)
	for i := 0; i < len(data); i++ {
		db.Exec(`INSERT INTO words (term, definition, module_id) VALUES ($1, $2, $3)`, data[i][0], data[i][1], moduleId)
	}
}
