package app

import (
	"mp3_desktop/mp3_desktop/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type Application struct {
	App    fyne.App
	Window fyne.Window
}

func NewApp() *Application {
	a := app.New()
	w := a.NewWindow("MP3 Player")

	return &Application{
		App:    a,
		Window: w,
	}
}

func (a *Application) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem("Музыка", ui.MusicTab()),
		container.NewTabItem("Плейлисты", ui.PlayListsTab()),
		container.NewTabItem("Загрузка трека", ui.NewUploadTab(a.Window)),
		//container.NewTabItem("Профиль", ui.ProfileTab()),
	)

	a.Window.SetContent(tabs)
	a.Window.Resize(fyne.NewSize(800, 600))
	a.Window.ShowAndRun()
}
