package fileStorage

import (
	"io"

	"github.com/minio/minio-go"
)

type fileStorage struct {
	client *minio.Client
}

func New(
	storageEndpoint,
	accessKeyID,
	secretAccessKey string,
	useSSL bool,
) (*fileStorage, error) {
	client, err := minio.New(storageEndpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return &fileStorage{client: client}, nil
}

func (r *fileStorage) Get(bucketName, objectName string) (*minio.Object, error) {
	f, err := r.client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *fileStorage) Put(
	bucketName,
	objectName string,
	reader io.Reader,
	objectSize int64,
) (n int64, err error) {
	n, err = r.client.PutObject(
		bucketName,
		objectName,
		reader,
		objectSize,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (r *fileStorage) Delete(bucketName, objectName string) error {
	err := r.client.RemoveObject(bucketName, objectName)
	if err != nil {
		return err
	}

	return nil
}
