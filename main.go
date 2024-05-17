package main

import (
	"fmt"
	"log"
	"quizlex/db"
	l "quizlex/layouts"
	p "quizlex/pages"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lib/pq"
)

func main() {
	// Screen size
	var win_x float32 = 1000
	var win_y float32 = 800
	var win_title string = "QuizLex"

	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow(win_title)
	window.Resize(fyne.NewSize(win_x, win_y))

	mainPage(window, []string{})
	window.ShowAndRun()

}
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func mainPage(window fyne.Window, searchingResults []string) {
	var modulesData []string

	db, err := db.ConnectDB()
	rows, err := db.Query("SELECT * FROM modules")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var createdAt time.Time

		err := rows.Scan(&id, &name, &createdAt)
		if err != nil {
			log.Fatal(err)
		}

		modulesData = append(modulesData, name)
	}

	// Searching
	searchingInput := widget.NewEntry()
	searchingInput.SetPlaceHolder("Searching...")
	search_button := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		searchingResults = searchingResults[:0]
		for _, item := range modulesData {
			if contains(item, searchingInput.Text) {
				searchingResults = append(searchingResults, item)
				fmt.Println(searchingResults)
			}
		}
		mainPage(window, searchingResults)
	})

	var modules []fyne.CanvasObject
	var visibleModules []string

	if len(searchingResults) > 0 {
		visibleModules = searchingResults
	} else {
		visibleModules = modulesData
	}

	for i := 0; i < len(visibleModules); i++ {

		removeModuleButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {

			dialog.ShowConfirm("Confirmation", "Are you sure you want to delete this module?", func(b bool) {
				if b {
					fmt.Println("remove")
					rows, err := db.Query("SELECT * FROM modules WHERE name = $1", visibleModules[i])
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

					db.Exec(`DELETE FROM words WHERE module_id = $1`, moduleId)

					db.Exec(`DELETE FROM modules WHERE id = $1`, moduleId)

					mainPage(window, []string{})

				}
			}, window)

		})

		moduleTemplate := widget.NewButton(visibleModules[i], func() {
			p.ModuleRead(window, visibleModules[i], mainPage)
		})

		// Spacer
		if i == 0 {
			modules = append(modules, l.NewAdaptiveGridWithRatios([]float32{1}, layout.NewSpacer()))
		}
		modules = append(modules, l.NewAdaptiveGridWithRatios([]float32{0.15, 0.6, 0.05, 0.05}, layout.NewSpacer(), moduleTemplate, removeModuleButton, layout.NewSpacer()))
	}

	add_button := widget.NewButtonWithIcon("Create", theme.ContentAddIcon(), func() {
		p.ModuleCreatingPage(window, mainPage)
	})

	add_button.Importance = widget.SuccessImportance

	window_content := container.NewVBox(
		l.NewAdaptiveGridWithRatios([]float32{0.15, 0.6, 0.05, 0.05}, layout.NewSpacer(), searchingInput, search_button),
		container.NewVBox(modules...),
		layout.NewSpacer(),
		l.NewAdaptiveGridWithRatios([]float32{0.4, 0.2, 0.4}, layout.NewSpacer(), add_button),
	)

	window.SetContent(window_content)
}
