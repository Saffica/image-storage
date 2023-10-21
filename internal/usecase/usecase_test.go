package usecase

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/Saffica/image-storage/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestUsecase(t *testing.T) {
	t.Run("getImg OK", func(t *testing.T) {
		inputUrl := "https://test"
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}

		_, err := imgService.GetImg(imgRequest)
		require.NoError(t, err)
	})

	t.Run("getImg validate error", func(t *testing.T) {
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})
		decodedUrl := "test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(decodedUrl)),
		}

		_, err := imgService.GetImg(imgRequest)
		require.ErrorIs(t, err, models.ErrBadHash)
	})

	t.Run("getImg OK with download from external service", func(t *testing.T) {
		inputUrl := "https://test"
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetMetaDataNotFoundError{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		_, err := imgService.GetImg(imgRequest)
		require.NoError(t, err)
	})

	t.Run("getImg donwload from external service with error", func(t *testing.T) {
		inputUrl := "https://test"
		imgService := New(
			&mockClientWithGetImgByURLError{},
			&mockMdRepositoryGetMetaDataNotFoundError{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
	})

	t.Run("getImg get meta data error", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetMetaDataOtherError{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
	})

	t.Run("getImg OK retry to download", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetMetaData24hours{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.NoError(t, err)
	})

	t.Run("getImg retry to download with error", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(
			&mockClientWithGetImgByURLError{},
			&mockMdRepositoryGetMetaData24hours{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
	})

	t.Run("getImg await time to download", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryAwaitTime{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
		require.ErrorIs(t, err, models.ErrImageNotFound)
	})

	t.Run("getImg get file error", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetDownloadedMetaData{},
			&mockImageModifyerService{},
			&mockFileRepositoryGetError{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
	})

	t.Run("getImg OK scale", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash:   base64.StdEncoding.EncodeToString([]byte(inputUrl)),
			Width:  1,
			Height: 1,
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetDownloadedMetaData{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.NoError(t, err)
	})

	t.Run("getImg scale error", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash:   base64.StdEncoding.EncodeToString([]byte(inputUrl)),
			Width:  1,
			Height: 1,
		}
		imgService := New(
			&mockClient{},
			&mockMdRepositoryGetDownloadedMetaData{},
			&mockImageModifyerServiceScaleError{},
			&mockFileRepository{},
		)
		_, err := imgService.GetImg(imgRequest)
		require.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		inputUrl := "https://test"
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte(inputUrl)),
		}
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})

		outputUrl, err := imgService.validate(imgRequest)
		require.NoError(t, err)
		require.Equal(t, inputUrl, outputUrl)
	})

	t.Run("decode error", func(t *testing.T) {
		imgRequest := &models.ImageRequest{
			Hash: "тест",
		}
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})

		_, err := imgService.validate(imgRequest)
		require.ErrorIs(t, err, models.ErrBadHash)
	})

	t.Run("parse url error", func(t *testing.T) {
		imgRequest := &models.ImageRequest{
			Hash: base64.StdEncoding.EncodeToString([]byte("тест")),
		}
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})

		_, err := imgService.validate(imgRequest)
		require.ErrorIs(t, err, models.ErrBadHash)
	})
}

func TestDownloadAndPrepareFile(t *testing.T) {
	t.Run("OK for exist image", func(t *testing.T) {
		mData := &models.MetaData{
			ID: 1,
		}
		mockMdR := &mockMdRepository{}
		imgService := New(&mockClient{}, mockMdR, &mockImageModifyerService{}, &mockFileRepository{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.NoError(t, err)
		require.NotEmpty(t, mockMdR.UpdateMetaData.UpdatedAt)
		require.True(t, mockMdR.UpdateMetaData.Downloaded)
	})

	t.Run("OK for not exist image", func(t *testing.T) {
		mData := &models.MetaData{}
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerService{}, &mockFileRepository{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.NoError(t, err)
	})

	t.Run("metaData insert error", func(t *testing.T) {
		mData := &models.MetaData{}
		imgService := New(&mockClient{}, &mockMdRepositoryWithInsertError{}, &mockImageModifyerService{}, &mockFileRepository{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.Error(t, err)
	})

	t.Run("getImgByURL error", func(t *testing.T) {
		mData := &models.MetaData{}
		imgService := New(
			&mockClientWithGetImgByURLError{},
			&mockMdRepository{},
			&mockImageModifyerService{},
			&mockFileRepository{},
		)
		_, err := imgService.downloadAndPrepareFile(mData)
		require.Error(t, err)
	})

	t.Run("convert to webp error", func(t *testing.T) {
		mData := &models.MetaData{}
		imgService := New(&mockClient{}, &mockMdRepository{}, &mockImageModifyerServiceWithError{}, &mockFileRepository{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.Error(t, err)
	})

	t.Run("file insert error", func(t *testing.T) {
		mData := &models.MetaData{}
		mdr := &mockMdRepository{}
		imgService := New(&mockClient{}, mdr, &mockImageModifyerService{}, &mockFileRepositoryInsertError{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.Error(t, err)
		require.False(t, mdr.UpdateMetaData.Downloaded)
	})

	t.Run("metaData update error", func(t *testing.T) {
		mData := &models.MetaData{}
		imgService := New(&mockClient{}, &mockMdRepositoryUpdateError{}, &mockImageModifyerService{}, &mockFileRepository{})
		_, err := imgService.downloadAndPrepareFile(mData)
		require.Error(t, err)
	})

}

type mockClient struct{}

func (c *mockClient) GetImgByURL(url string) ([]byte, error) {
	return nil, nil
}

type mockMdRepository struct {
	UpdateMetaData *models.MetaData
}

func (r *mockMdRepository) Get(url string) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}
func (r *mockMdRepository) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}
func (r *mockMdRepository) Update(metaData *models.MetaData) (*models.MetaData, error) {
	r.UpdateMetaData = metaData
	return &models.MetaData{}, nil
}

type mockImageModifyerService struct{}

func (r *mockImageModifyerService) ConvertToWebp(image []byte) ([]byte, error) {
	return nil, nil
}

func (r *mockImageModifyerService) Scale(image []byte, width, height uint16) ([]byte, error) {
	return nil, nil
}

type mockFileRepository struct{}

func (r *mockFileRepository) Get(id int64) ([]byte, error) {
	return nil, nil
}

func (r *mockFileRepository) Insert(id int64, file []byte) error {
	return nil
}

type mockMdRepositoryWithInsertError struct{}

func (r *mockMdRepositoryWithInsertError) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, errors.New("error")
}

func (r *mockMdRepositoryWithInsertError) Get(url string) (*models.MetaData, error) {
	return nil, nil
}
func (r *mockMdRepositoryWithInsertError) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}

type mockClientWithGetImgByURLError struct{}

func (r *mockClientWithGetImgByURLError) GetImgByURL(url string) ([]byte, error) {
	return nil, errors.New("error")
}

type mockImageModifyerServiceWithError struct{}

func (r *mockImageModifyerServiceWithError) ConvertToWebp(image []byte) ([]byte, error) {
	return nil, errors.New("error")
}

func (r *mockImageModifyerServiceWithError) Scale(image []byte, width, height uint16) ([]byte, error) {
	return nil, nil
}

type mockFileRepositoryInsertError struct{}

func (r *mockFileRepositoryInsertError) Get(id int64) ([]byte, error) {
	return nil, nil
}

func (r *mockFileRepositoryInsertError) Insert(id int64, file []byte) error {
	return errors.New("error")
}

type mockMdRepositoryUpdateError struct{}

func (r *mockMdRepositoryUpdateError) Get(url string) (*models.MetaData, error) {
	return nil, nil
}

func (r *mockMdRepositoryUpdateError) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

func (r *mockMdRepositoryUpdateError) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, errors.New("error")
}

type mockMdRepositoryGetMetaDataNotFoundError struct{}

func (r *mockMdRepositoryGetMetaDataNotFoundError) Get(url string) (*models.MetaData, error) {
	return nil, models.ErrMetaDataNotFound
}

func (r *mockMdRepositoryGetMetaDataNotFoundError) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

func (r *mockMdRepositoryGetMetaDataNotFoundError) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}

type mockMdRepositoryGetMetaDataOtherError struct{}

func (r *mockMdRepositoryGetMetaDataOtherError) Get(url string) (*models.MetaData, error) {
	return nil, errors.New("error")
}

func (r *mockMdRepositoryGetMetaDataOtherError) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}

func (r *mockMdRepositoryGetMetaDataOtherError) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}

type mockMdRepositoryGetMetaData24hours struct{}

func (r *mockMdRepositoryGetMetaData24hours) Get(url string) (*models.MetaData, error) {
	return &models.MetaData{
		ID:         1,
		Downloaded: false,
		UpdatedAt:  time.Now().Add(-25 * time.Hour),
	}, nil
}

func (r *mockMdRepositoryGetMetaData24hours) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

func (r *mockMdRepositoryGetMetaData24hours) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

type mockMdRepositoryAwaitTime struct{}

func (r *mockMdRepositoryAwaitTime) Get(url string) (*models.MetaData, error) {
	return &models.MetaData{
		ID:         1,
		Downloaded: false,
		UpdatedAt:  time.Now(),
	}, nil
}

func (r *mockMdRepositoryAwaitTime) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

func (r *mockMdRepositoryAwaitTime) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

type mockFileRepositoryGetError struct{}

func (r *mockFileRepositoryGetError) Get(id int64) ([]byte, error) {
	return nil, errors.New("error")
}

func (r *mockFileRepositoryGetError) Insert(id int64, file []byte) error {
	return nil
}

type mockMdRepositoryGetDownloadedMetaData struct{}

func (r *mockMdRepositoryGetDownloadedMetaData) Get(url string) (*models.MetaData, error) {
	return &models.MetaData{
		ID:         1,
		Downloaded: true,
	}, nil
}
func (r *mockMdRepositoryGetDownloadedMetaData) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}
func (r *mockMdRepositoryGetDownloadedMetaData) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return &models.MetaData{}, nil
}

type mockImageModifyerServiceScaleError struct{}

func (r *mockImageModifyerServiceScaleError) ConvertToWebp(image []byte) ([]byte, error) {
	return nil, nil
}

func (r *mockImageModifyerServiceScaleError) Scale(image []byte, width, height uint16) ([]byte, error) {
	return nil, errors.New("error")
}
