package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ImageQuizApp struct {
	window      fyne.Window
	imageFolder string
	imageFiles  []string
	score       int
	total       int
	buttons     []*widget.Button
	imageCanvas *canvas.Image
	scoreLabel  *widget.Label
	buttonBox   *fyne.Container
	darkMode    bool
}

func NewImageQuizApp(app fyne.App) *ImageQuizApp {
	mainWindow := app.NewWindow("Student Name Learner")
	quizApp := &ImageQuizApp{
		window:   mainWindow,
		score:    0,
		total:    0,
		darkMode: false,
	}
	quizApp.setupUI()
	return quizApp
}

func (app *ImageQuizApp) setupUI() {
	rand.Seed(time.Now().UnixNano())

	// Menu bar
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Select image folder", func() { app.openPreferences() }),
			fyne.NewMenuItem("About", func() { app.openAbout() }),
		),
		fyne.NewMenu("Appearance",
			fyne.NewMenuItem("Toggle theme", func() { app.toggleDarkMode() }),
		),
	)
	app.window.SetMainMenu(menu)

	// Image display
	app.imageCanvas = &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	app.imageCanvas.SetMinSize(fyne.NewSize(300, 300))

	// Buttons
	app.buttonBox = container.NewHBox()
	for i := 0; i < 3; i++ {
		button := widget.NewButton("", nil)
		button.Hide() // Initially hide buttons
		app.buttons = append(app.buttons, button)
		app.buttonBox.Add(button)
	}

	// Score label
	app.scoreLabel = widget.NewLabel("Score: 0/0")

	// Layout
	mainLayout := container.NewVBox(
		app.imageCanvas,
		app.buttonBox,
		app.scoreLabel,
	)
	app.window.SetContent(mainLayout)
	app.loadImages()
	app.nextImage()
}

func (app *ImageQuizApp) loadImages() {
	if app.imageFolder == "" {
		return
	}
	files, err := os.ReadDir(app.imageFolder)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.imageFiles = nil
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
				app.imageFiles = append(app.imageFiles, file.Name())
			}
		}
	}
}

func (app *ImageQuizApp) nextImage() {
	if len(app.imageFiles) < 3 {
		dialog.ShowInformation("Info", "Not enough images in the folder. Go to File > Select image folder", app.window)
		app.toggleButtons(false) // Hide buttons if not enough images
		return
	}

	app.toggleButtons(true) // Show buttons when enough images are available

	correctIndex := rand.Intn(3)
	correctFile := app.imageFiles[rand.Intn(len(app.imageFiles))]
	correctName := strings.TrimSuffix(correctFile, filepath.Ext(correctFile))

	app.imageCanvas.File = filepath.Join(app.imageFolder, correctFile)
	app.imageCanvas.Refresh()

	options := make([]string, 3)
	options[correctIndex] = correctName
	for i := 0; i < 3; i++ {
		if i == correctIndex {
			continue
		}
		for {
			fakeFile := app.imageFiles[rand.Intn(len(app.imageFiles))]
			fakeName := strings.TrimSuffix(fakeFile, filepath.Ext(fakeFile))
			if fakeName != correctName && !contains(options, fakeName) {
				options[i] = fakeName
				break
			}
		}
	}

	for i, button := range app.buttons {
		button.SetText(options[i])
		button.OnTapped = func(i int) func() {
			return func() {
				app.checkAnswer(i == correctIndex, correctName)
			}
		}(i)
	}
}

func (app *ImageQuizApp) toggleButtons(visible bool) {
	for _, button := range app.buttons {
		if visible {
			button.Show()
		} else {
			button.Hide()
		}
	}
}

func (app *ImageQuizApp) checkAnswer(correct bool, correctName string) {
	app.total++
	if correct {
		app.score++
		dialog.ShowInformation("Correct", "Correct!", app.window)
	} else {
		dialog.ShowInformation("Incorrect", fmt.Sprintf("Incorrect! The correct answer was %s", correctName), app.window)
	}
	app.scoreLabel.SetText(fmt.Sprintf("Score: %d/%d", app.score, app.total))
	app.nextImage()
}

func (app *ImageQuizApp) openPreferences() {
	dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
		if err != nil || folder == nil {
			return
		}
		app.imageFolder = folder.Path()
		app.loadImages()
		app.nextImage()
	}, app.window)
}

func (app *ImageQuizApp) openAbout() {
	dialog.ShowInformation("About", "© Ollie Foster 2024.\nAll rights reserved.", app.window)
}

func (app *ImageQuizApp) toggleDarkMode() {
	app.darkMode = !app.darkMode
	if app.darkMode {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	} else {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func main() {
	myApp := app.New()
	quizApp := NewImageQuizApp(myApp)
	quizApp.window.ShowAndRun()
}
