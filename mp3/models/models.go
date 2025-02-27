package models

type Track struct {
	ID     int
	Name   string `xml:"name"`
	Artist string `xml:"artist"`
}
