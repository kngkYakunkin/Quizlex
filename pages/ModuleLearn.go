package pages

import (
	"image/color"
	"math/rand"
	"quizlex/layouts"
	"quizlex/widgets"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ModuleLearn(window fyne.Window, moduleName string, data [][2]string, pastData [][2]string, mistakes []string, mainPage func(window fyne.Window, searchingResults []string)) {
	var page_title string = "Learn Module"
	window.SetTitle(page_title)

	// Header
	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		ModuleRead(window, moduleName, mainPage)
	})

	moduleTitle := canvas.NewText(moduleName, theme.ForegroundColor())
	moduleTitle.TextSize = 25
	moduleTitle.TextStyle.Bold = true
	moduleTitle.Alignment = fyne.TextAlignCenter
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})

	// Current Term
	var correctAnswerIdx int = rand.Intn(4)
	var currentTermIdx int = rand.Intn(len(pastData))
	currentTermLabel := canvas.NewText(pastData[currentTermIdx][0], theme.ForegroundColor())
	currentTermLabel.TextSize = 30
	currentTermLabel.TextStyle.Bold = true

	// Question label
	answerButtonsTitle := canvas.NewText("Choose the correct answer: ", theme.ForegroundColor())
	answerButtonsTitle.TextSize = 20
	answerButtonsTitle.Alignment = fyne.TextAlignCenter

	// Result label
	resultLabel := canvas.NewText("", color.RGBA{120, 22, 22, 180})
	resultLabel.TextSize = 18
	resultLabel.Alignment = fyne.TextAlignCenter

	// Progress
	var currentProgress int = len(data) - len(pastData) + 1
	progressLabel := canvas.NewText(strconv.Itoa(currentProgress)+" / "+strconv.Itoa(len(data)), theme.ForegroundColor())
	progressLabel.TextSize = 17
	progressLabel.TextStyle.Bold = true
	progressLabel.Alignment = fyne.TextAlignCenter

	// Terms indeces that have already been used
	randIndices := []int{}

	// To prevent button clicks after the answer
	var isEnd bool = false

	var answerButtons []fyne.CanvasObject

	for i := 0; i < 4; i++ {
		if i != correctAnswerIdx {
			// Get random word from all words
			var randIdx int
			var isFound bool = false

			for !isFound {
				randIdx = rand.Intn(len(data))
				// If it's the first element
				if len(randIndices) == 0 {
					// if randIdx is already in use by correct answer
					if data[randIdx] == pastData[currentTermIdx] {
						continue
					}
					randIndices = append(randIndices, randIdx)
					//randIndices[i] = randIdx
					isFound = true
				} else {
					// trying to find unique element
					for n := 0; n < len(randIndices); n++ {
						if randIdx == randIndices[n] || data[randIdx][0] == pastData[currentTermIdx][0] {
							break
						}
						// If it's the last element then randIdx is unique
						if n == len(randIndices)-1 {
							isFound = true
							randIndices = append(randIndices, randIdx)
							//randIndices[i] = randIdx
						}
					}
				}
			}

		}

		if i == correctAnswerIdx {
			// Right answer button
			answerBtn := widgets.NewAnswerButton(theme.ButtonColor(), color.RGBA{255, 255, 255, 40}, pastData[currentTermIdx][1])
			answerBtn.Clicked = func() {
				if !isEnd {
					// Show right answer
					answerBtn.RightAnswer = true
					answerBtn.Refresh()
					// Title
					resultLabel.Text = "Correct!"
					resultLabel.Color = color.RGBA{0, 128, 0, 255}

					// Next
					window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
						if key.Name == fyne.KeyReturn {

							pastData := append(pastData[:currentTermIdx], pastData[currentTermIdx+1:]...)
							// If there's no more words
							if len(pastData) == 0 {
								ModuleLearnResults(window, moduleName, data, len(mistakes), mainPage)
							} else {
								ModuleLearn(window, moduleName, data, pastData, mistakes, mainPage)
							}
						}
					})

					isEnd = true
				}
			}
			answerButtons = append(answerButtons, answerBtn)

		} else {
			// Wrong answer button
			var lastRandIndex int = len(randIndices) - 1
			answerBtn := widgets.NewAnswerButton(theme.ButtonColor(), color.RGBA{255, 255, 255, 40}, data[randIndices[lastRandIndex]][1])
			answerBtn.Clicked = func() {

				if !isEnd {
					// Show wrong answer
					answerBtn.WrongAnswer = true
					answerBtn.Refresh()
					// Show right answer
					correctAnswerButton := answerButtons[correctAnswerIdx].(*widgets.AnswerButton)
					correctAnswerButton.RightAnswer = true
					correctAnswerButton.Refresh()
					// Title
					resultLabel.Text = "Incorrect!"
					resultLabel.Color = color.RGBA{255, 0, 0, 255}

					// Next
					window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
						if key.Name == fyne.KeyReturn {
							var isSameMistake = false
							// Add a unique mistake
							for m := 0; m < len(mistakes); m++ {
								if mistakes[m] == pastData[currentTermIdx][0] {
									isSameMistake = true
									break
								}
							}
							if !isSameMistake {
								mistakes = append(mistakes, pastData[currentTermIdx][0])
							}

							ModuleLearn(window, moduleName, data, pastData, mistakes, mainPage)
						}
					})

					isEnd = true
				}
			}

			answerButtons = append(answerButtons, answerBtn)
		}

	}

	content := container.NewBorder(
		container.NewVBox(layouts.NewAdaptiveGridWithRatios([]float32{0.05, 0.9, 0.05}, backButton, moduleTitle, settingsButton), progressLabel),
		container.NewVBox(
			answerButtonsTitle,
			widget.NewLabel(""),
			layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.4, 0.4, 0.1}, layout.NewSpacer(), answerButtons[0], answerButtons[1], layout.NewSpacer()),
			layouts.NewAdaptiveGridWithRatios([]float32{0.1, 0.4, 0.4, 0.1}, layout.NewSpacer(), answerButtons[2], answerButtons[3], layout.NewSpacer()),
			widget.NewLabel(""),
			resultLabel,
			widget.NewLabel(""),
		),
		nil, nil, container.NewCenter(container.NewVBox(currentTermLabel)),
	)

	window.SetContent(content)
}
