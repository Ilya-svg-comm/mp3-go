// В файле main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"mp3/models/psql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=postgres dbname=player sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	trackModel := &psql.TrackModel{DB: db}

	// Исправленный вызов GetAll()
	tracks, err := trackModel.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	a := app.New()
	w := a.NewWindow("MP3-Player")
	w.Resize(fyne.NewSize(600, 400))

	list := widget.NewList(
		func() int { return len(tracks) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			track := tracks[i]
			obj.(*widget.Label).SetText(fmt.Sprintf("%d. %s - %s", track.ID, track.Artist, track.Name))
		},
	)

	w.SetContent(container.NewVBox(list))
	w.ShowAndRun()
}
