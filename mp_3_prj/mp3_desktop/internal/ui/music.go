package ui

import (
	"fmt"
	"mp3_desktop/mp3_desktop/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MusicTab() fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			tracks, _ := service.GetTracks("http://localhost:8080")
			return len(tracks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("...")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			tracks, _ := service.GetTracks("http://localhost:8080")
			o.(*widget.Label).SetText(fmt.Sprintf("%s — %s", tracks[i].Metadata.Artist, tracks[i].Metadata.Title))
		},
	)

	return container.NewVBox(
		widget.NewLabel("Список треков"),
		list,
	)
}
