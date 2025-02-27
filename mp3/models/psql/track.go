package psql

import (
	"database/sql"
	"encoding/xml"
	"fmt"

	"mp3/models"
)

// TrackModel для работы с таблицей tracks
type TrackModel struct {
	DB *sql.DB
}

// GetAll получает все записи из таблицы tracks и выводит их в консоль
func (m *TrackModel) GetAll() ([]models.Track, error) {
	rows, err := m.DB.Query("SELECT id, xmldata FROM tracks")
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var id int
		var xmlData string
		if err := rows.Scan(&id, &xmlData); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %w", err)
		}

		var track models.Track
		if err := xml.Unmarshal([]byte(xmlData), &track); err != nil {
			return nil, fmt.Errorf("ошибка при разборе XML: %w", err)
		}

		track.ID = id // Устанавливаем ID из БД
		tracks = append(tracks, track)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	return tracks, nil
}
