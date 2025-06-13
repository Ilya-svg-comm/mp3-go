package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	return &MinioClient{
		client: client,
		bucket: bucket,
	}, nil
}

func (m *MinioClient) UploadBytes(data []byte, objectName string) (string, error) {
	_, err := m.client.PutObject(
		context.Background(),
		m.bucket,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload bytes: %v", err)
	}
	return objectName, nil
}

func (m *MinioClient) GetFileURL(objectKey string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucket, objectKey)
}

func (m *MinioClient) DeleteObject(objectKey string) error {
	return m.client.RemoveObject(
		context.Background(),
		m.bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
}

func (m *MinioClient) DownloadStream(objectKey string) (io.ReadSeekCloser, int64, error) {
	obj, err := m.client.GetObject(context.Background(), m.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	// Узнаем размер файла
	info, err := obj.Stat()
	if err != nil {
		return nil, 0, err
	}

	return obj, info.Size, nil
}
