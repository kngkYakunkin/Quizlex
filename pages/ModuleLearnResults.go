package pages

import (
	"image/color"
	"quizlex/layouts"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ModuleLearnResults(window fyne.Window, moduleName string, data [][2]string, mistakesCounter int, mainPage func(window fyne.Window, searchingResults []string)) {
	var page_title string = "Results"
	window.SetTitle(page_title)

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {})
	moduleTitle := canvas.NewText(moduleName, theme.ForegroundColor())
	moduleTitle.TextSize = 25
	moduleTitle.TextStyle.Bold = true
	moduleTitle.Alignment = fyne.TextAlignCenter
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})

	resultsTitle := canvas.NewText("Module completed!", theme.ForegroundColor())
	resultsTitle.TextSize = 28
	resultsTitle.TextStyle.Bold = true

	rightLabel := canvas.NewText(strconv.Itoa(len(data)-mistakesCounter)+" right", color.RGBA{0, 128, 0, 255})
	rightLabel.Alignment = fyne.TextAlignCenter

	wrongLabel := canvas.NewText(strconv.Itoa(mistakesCounter)+" wrong", color.RGBA{255, 0, 0, 255})
	wrongLabel.Alignment = fyne.TextAlignCenter

	againButton := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		pastData := make([][2]string, len(data))
		copy(pastData, data)
		ModuleLearn(window, moduleName, data, pastData, []string{}, mainPage)
	})

	content := container.NewBorder(
		container.NewVBox(layouts.NewAdaptiveGridWithRatios([]float32{0.05, 0.9, 0.05}, backButton, moduleTitle, settingsButton)),
		nil, nil, nil,
		container.NewCenter(
			container.NewVBox(
				resultsTitle, layouts.NewAdaptiveGridWithRatios([]float32{0.5, 0.5}, rightLabel, wrongLabel),
				widget.NewLabel(""),
				layouts.NewAdaptiveGridWithRatios([]float32{0.40, 0.2, 0.40}, layout.NewSpacer(), againButton, layout.NewSpacer()),
			),
		),
	)

	window.SetContent(content)

}
