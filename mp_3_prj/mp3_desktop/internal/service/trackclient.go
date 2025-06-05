package service

import (
	"encoding/json"
	"net/http"
	"time"
)

// Track — структура для получения из API
type Track struct {
	ID        int       `json:"id"`
	Metadata  Metadata  `json:"metadata"`
	ObjectKey string    `json:"object_key"`
	CreatedAt time.Time `json:"created_at"`
}

type Metadata struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration int    `json:"duration"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
}

func GetTracks(apiURL string) ([]Track, error) {
	resp, err := http.Get(apiURL + "/tracks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tracks []Track
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
