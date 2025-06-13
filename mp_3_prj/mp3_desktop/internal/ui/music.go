package ui

import (
	"fmt"
	"mp3_desktop/mp3_desktop/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MusicTab() fyne.CanvasObject {
	tracks, _ := service.GetTracks("http://localhost:8080")
	var selectedIndex = -1

	// Контролы
	btnPlay := widget.NewButton("Play", nil)
	btnStop := widget.NewButton("Stop", func() {
		service.StopAudio()
	})
	btnPlay.Disable()

	list := widget.NewList(
		func() int { return len(tracks) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s — %s", tracks[i].Metadata.Artist, tracks[i].Metadata.Title))
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		btnPlay.Enable()
	}

	btnPlay.OnTapped = func() {
		if selectedIndex < 0 {
			return
		}
		url := fmt.Sprintf("http://localhost:8080/tracks/%d/audio", tracks[selectedIndex].ID)
		err := service.PlayAudio(url)
		if err != nil {
			fmt.Println("Error playing:", err)
		}
	}

	controlBar := container.NewHBox(btnPlay, btnStop)

	return container.NewBorder(
		widget.NewLabelWithStyle("Список треков", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		controlBar,
		nil, nil,
		list,
	)
}
