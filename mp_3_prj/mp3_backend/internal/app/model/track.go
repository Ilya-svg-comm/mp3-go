package model

import (
	"encoding/xml"
	"time"
)

type Track struct {
	ID        int       `json:"id"`
	Metadata  Metadata  `json:"metadata"`
	ObjectKey string    `json:"object_key"`
	CreatedAt time.Time `json:"created_at"`
	AudioURL  string    `json:"audio_url" xml:"audio_url,omitempty"`
}

type Metadata struct {
	XMLName  xml.Name `xml:"metadata"`
	Title    string   `xml:"title"`
	Artist   string   `xml:"artist"`
	Album    string   `xml:"album"`
	Duration int      `xml:"duration"`
	Year     int      `xml:"year"`
	Genre    string   `xml:"genre"`
}
