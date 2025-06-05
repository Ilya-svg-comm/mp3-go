package store

import (
	"backend/internal/app/model"
	"encoding/xml"
	"fmt"
	"path/filepath"
)

type TrackRepository struct {
	store *Store
	minio *MinioClient
}

func (r *TrackRepository) Create(t *model.Track, audioFile []byte) (*model.Track, error) {
	// Генерируем objectKey (пример: "tracks/{artist}/{title}.mp3")
	objectKey := fmt.Sprintf("tracks/%s/%s.mp3",
		t.Metadata.Artist,
		filepath.Base(t.Metadata.Title))

	// Загружаем аудиофайл в MinIO
	_, err := r.minio.UploadBytes(audioFile, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to upload audio to MinIO: %v", err)
	}
	t.ObjectKey = objectKey

	// Сериализуем метаданные в XML
	metadataXML, err := xml.Marshal(t.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}

	// Сохраняем в PostgreSQL
	query := `INSERT INTO "default".tracks (metadata, object_key) 
			  VALUES ($1, $2) 
			  RETURNING id, created_at`

	err = r.store.db.QueryRow(
		query,
		metadataXML,
		t.ObjectKey,
	).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		// При ошибке откатываем загрузку в MinIO (опционально)
		_ = r.minio.DeleteObject(objectKey)
		return nil, fmt.Errorf("database error: %v", err)
	}

	return t, nil
}

func (r *TrackRepository) GetAll() ([]*model.Track, error) {
	query := `SELECT id, metadata, object_key, created_at 
			  FROM "default".tracks 
			  ORDER BY id DESC`

	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var tracks []*model.Track
	for rows.Next() {
		t := &model.Track{}
		var metadataXML []byte

		err := rows.Scan(
			&t.ID,
			&metadataXML,
			&t.ObjectKey,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		// Десериализуем XML
		if err := xml.Unmarshal(metadataXML, &t.Metadata); err != nil {
			return nil, fmt.Errorf("metadata unmarshal error: %v", err)
		}

		// Добавляем URL для доступа к файлу
		t.AudioURL = r.minio.GetFileURL(t.ObjectKey)
		tracks = append(tracks, t)
	}

	return tracks, nil
}
