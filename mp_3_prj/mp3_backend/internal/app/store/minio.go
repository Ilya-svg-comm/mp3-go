package store

import (
	"bytes"
	"context"
	"fmt"

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
		return nil, fmt.Errorf("MINIOCLIENT: FAILT, ERROR:", err)
	}
	return &MinioClient{
		client: client,
		bucket: bucket,
	}, nil
}

func (m *MinioClient) UploadFIle(filePath, objectName string) (string, error) {
	_, err := m.client.FPutObject(context.Background(), m.bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	return objectName, nil
}

func (m *MinioClient) GetFileURL(objectKey string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucket, objectKey)
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
	return objectName, err
}

func (m *MinioClient) DeleteObject(objectKey string) error {
	return m.client.RemoveObject(
		context.Background(),
		m.bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
}
