package file

import (
	"bytes"
	"io"
	"strconv"

	"github.com/minio/minio-go"
)

type fStorage interface {
	Get(bucketName, objectName string) (*minio.Object, error)
	Put(
		bucketName string,
		objectName string,
		reader io.Reader,
		objectSize int64,
	) (n int64, err error)
	Delete(bucketName, objectName string) error
}

type fileRepository struct {
	f fStorage
}

func New(fileStorage fStorage) *fileRepository {
	return &fileRepository{
		f: fileStorage,
	}
}

func (r *fileRepository) Get(bucketName string, id int64) ([]byte, error) {
	file, err := r.f.Get(bucketName, r.prepareId(id))
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return b, nil

}

func (r *fileRepository) Put(bucketName string, id int64, file []byte) error {
	reader := bytes.NewReader(file)
	_, err := r.f.Put(bucketName, r.prepareId(id), reader, reader.Size())
	if err != nil {
		return err
	}

	return nil
}

func (r *fileRepository) Delete(bucketName string, id int64) error {
	err := r.f.Delete(bucketName, r.prepareId(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *fileRepository) prepareId(id int64) string {
	return strconv.FormatInt(id, 10)
}
