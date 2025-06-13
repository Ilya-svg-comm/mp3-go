// ui/upload_tab.go
package ui

import (
	"fmt"
	"log"
	"mp3_desktop/mp3_desktop/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewUploadTab(win fyne.Window) fyne.CanvasObject {
	var selectedFile string

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Название трека")

	artistEntry := widget.NewEntry()
	artistEntry.SetPlaceHolder("Исполнитель")

	fileLabel := widget.NewLabel("Файл не выбран")

	selectFileBtn := widget.NewButton("Выбрать файл", func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			selectedFile = uc.URI().Path()
			fileLabel.SetText("Выбран: " + selectedFile)
		}, win)
		fd.Show()
	})

	uploadBtn := widget.NewButton("Загрузить", func() {
		if selectedFile == "" || titleEntry.Text == "" || artistEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("все поля должны быть заполнены"), win)
			return
		}

		err := service.UploadTrack(selectedFile, titleEntry.Text, artistEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("ошибка загрузки: %w", err), win)
			log.Println("Upload error:", err)
		} else {
			dialog.ShowInformation("Успех", "Трек успешно загружен", win)
			titleEntry.SetText("")
			artistEntry.SetText("")
			fileLabel.SetText("Файл не выбран")
			selectedFile = ""
		}
	})

	form := container.NewVBox(
		titleEntry,
		artistEntry,
		fileLabel,
		selectFileBtn,
		layout.NewSpacer(),
		uploadBtn,
	)

	return container.NewVBox(
		widget.NewLabelWithStyle("Загрузка трека", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
	)
}

// func storageFilterAudio() fyne.FileFilter {
// 	return storage.NewExtensionFileFilter([]string{".mp3", ".wav", ".flac"})
// }
