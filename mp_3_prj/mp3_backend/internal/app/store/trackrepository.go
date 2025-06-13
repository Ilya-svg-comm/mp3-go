package store

import (
	"backend/internal/app/model"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type TrackRepository struct {
	store *Store
	minio *MinioClient
}

func (r *TrackRepository) Create(t *model.Track, audioData []byte) (*model.Track, error) {
	// Генерируем безопасный objectKey
	objectKey := fmt.Sprintf("tracks/%s/%s.mp3",
		sanitize(t.Metadata.Artist),
		sanitize(t.Metadata.Title))

	// Загружаем в MinIO
	_, err := r.minio.UploadBytes(audioData, objectKey)
	if err != nil {
		return nil, fmt.Errorf("minio upload failed: %v", err)
	}

	// Сохраняем метаданные в БД
	metadataXML, err := xml.Marshal(t.Metadata)
	if err != nil {
		return nil, fmt.Errorf("xml marshal failed: %v", err)
	}

	err = r.store.db.QueryRow(
		`INSERT INTO "default".tracks (metadata, object_key) 
		 VALUES ($1, $2)
		 RETURNING id`,
		metadataXML,
		objectKey,
	).Scan(&t.ID)

	if err != nil {
		// Откатываем загрузку в MinIO при ошибке
		_ = r.minio.DeleteObject(objectKey)
		return nil, fmt.Errorf("db insert failed: %v", err)
	}

	return t, nil
}

func (r *TrackRepository) GetAll() ([]*model.Track, error) {
	rows, err := r.store.db.Query(
		`SELECT id, metadata, object_key 
		 FROM "default".tracks 
		 ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %v", err)
	}
	defer rows.Close()

	var tracks []*model.Track
	for rows.Next() {
		var t model.Track
		var metadataXML []byte

		if err := rows.Scan(&t.ID, &metadataXML, &t.ObjectKey); err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		if err := xml.Unmarshal(metadataXML, &t.Metadata); err != nil {
			return nil, fmt.Errorf("xml unmarshal failed: %v", err)
		}

		t.AudioURL = r.minio.GetFileURL(t.ObjectKey)
		tracks = append(tracks, &t)
	}

	return tracks, nil
}

func sanitize(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), " ", "_")
}

func (r *TrackRepository) FindByID(id int) (*model.Track, error) {
	var t model.Track
	var metadataXML []byte

	err := r.store.db.QueryRow(
		`SELECT id, metadata, object_key FROM "default".tracks WHERE id = $1`, id,
	).Scan(&t.ID, &metadataXML, &t.ObjectKey)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(metadataXML, &t.Metadata); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TrackRepository) GetAudioStream(objectKey string) (io.ReadSeekCloser, int64, error) {
	return r.minio.DownloadStream(objectKey)
}
