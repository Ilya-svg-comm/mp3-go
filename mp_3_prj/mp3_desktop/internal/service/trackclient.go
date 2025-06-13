package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func UploadTrack(filePath, title, artist string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("не удалось создать поле file: %w", err)
	}
	if _, err := io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("ошибка копирования файла: %w", err)
	}

	_ = writer.WriteField("title", title)
	_ = writer.WriteField("artist", artist)

	if err := writer.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/tracks", &requestBody)
	if err != nil {
		return fmt.Errorf("не удалось создать запрос: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("сервер вернул ошибку: %s", string(body))
	}

	return nil
}
