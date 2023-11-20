package file

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Saffica/image-storage/pkg/fileStorage"
)

func TestFile(t *testing.T) {
	const minioEndpoint = "localhost:9000"
	const accessKeyID = "tvkPFzqWyavI2q4hwfE0"
	const secretAccessKey = "hcNgl4U8nSwU8C0EpenBepd0zIl7o4RUCxEdLknN"
	const useSSL = false
	const objectName = 1
	const bucketName = "main"
	const badBucketName = "not-found"
	fStorage, err := fileStorage.New(minioEndpoint, accessKeyID, secretAccessKey, useSSL)

	require.NoError(t, err)

	fileRepository := New(fStorage)

	defer func() {
		err = fileRepository.Delete(bucketName, objectName)
		require.NoError(t, err)
	}()

	t.Run("Put OK", func(t *testing.T) {
		err = fileRepository.Put(bucketName, objectName, []byte("File from application"))
		require.NoError(t, err)
	})

	t.Run("Put with error", func(t *testing.T) {
		err = fileRepository.Put(badBucketName, objectName, []byte("File from application"))
		require.Error(t, err)
	})

	t.Run("Get OK", func(t *testing.T) {
		f, err := fileRepository.Get(bucketName, objectName)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(f), 0)
	})

	t.Run("Get with error", func(t *testing.T) {
		_, err := fileRepository.Get(bucketName, 0)
		require.Error(t, err)
	})

	t.Run("Delete OK", func(t *testing.T) {
		err := fileRepository.Delete(bucketName, objectName)
		require.NoError(t, err)
	})

	t.Run("Delete with error", func(t *testing.T) {
		err := fileRepository.Delete(badBucketName, 0)
		require.Error(t, err)
	})
}
