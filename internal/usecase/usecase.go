package usecase

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Saffica/image-storage/pkg/models"
)

type FileRepository interface {
	Get(id int64) ([]byte, error)
	Insert(id int64, file []byte) error
}

type ImageModifyerService interface {
	ConvertToWebp(image []byte) ([]byte, error)
	Scale(image []byte, width, height uint16) ([]byte, error)
}

type MetaDataRepository interface {
	Get(url string) (*models.MetaData, error)
	Insert(metaData *models.MetaData) (*models.MetaData, error)
	Update(metaData *models.MetaData) (*models.MetaData, error)
}

type Client interface {
	GetImgByURL(url string) ([]byte, error)
}

type imgService struct {
	client               Client
	metaDataRepository   MetaDataRepository
	imageModifyerService ImageModifyerService
	fileRepository       FileRepository
}

func New(
	client Client, mdRepository MetaDataRepository, imgMService ImageModifyerService, fileRepository FileRepository,
) *imgService {
	return &imgService{
		client:               client,
		metaDataRepository:   mdRepository,
		imageModifyerService: imgMService,
		fileRepository:       fileRepository,
	}
}

func (s *imgService) GetImg(imageRequest *models.ImageRequest) ([]byte, error) {
	validUrl, err := s.validate(imageRequest)
	if err != nil {
		return nil, err
	}

	metaData, err := s.metaDataRepository.Get(validUrl)
	switch {
	case err == nil:
	case errors.Is(err, models.ErrMetaDataNotFound):
		metaData = &models.MetaData{
			DownloadLink: validUrl,
		}
		img, err := s.downloadAndPrepareFile(metaData)
		if err != nil {
			return nil, err
		}

		return img, nil
	default:
		return nil, err
	}

	day := 24 * time.Hour
	canDownloadAgain := metaData.UpdatedAt.Add(day).Before(time.Now())
	if !metaData.Downloaded && canDownloadAgain {
		img, err := s.downloadAndPrepareFile(metaData)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	if !metaData.Downloaded && !canDownloadAgain {
		return nil, models.ErrImageNotFound
	}

	outputImage, err := s.fileRepository.Get(metaData.ID)
	if err != nil {
		return nil, err
	}

	if imageRequest.Width != 0 && imageRequest.Height != 0 {
		outputImage, err = s.imageModifyerService.Scale(outputImage, imageRequest.Width, imageRequest.Height)
		if err != nil {
			return nil, err
		}
	}

	return outputImage, nil
}

func (s *imgService) validate(imgRequest *models.ImageRequest) (
	outputUrl string, err error,
) {
	decodedUrl, err := base64.StdEncoding.DecodeString(imgRequest.Hash)
	if err != nil {
		return "", fmt.Errorf("%w: %s", models.ErrBadHash, err.Error())
	}

	u, err := url.ParseRequestURI(string(decodedUrl))
	if err != nil {
		return "", fmt.Errorf("%w: %s", models.ErrBadHash, err.Error())
	}

	return u.String(), nil
}

func (s *imgService) downloadAndPrepareFile(metaData *models.MetaData) (
	file []byte, err error,
) {
	var webpImage []byte

	if metaData.ID == int64(0) {
		metaData, err = s.metaDataRepository.Insert(metaData)
		if err != nil {
			return nil, err
		}
	}

	defer func() {
		metaData.UpdatedAt = time.Now()
		if err != nil {
			metaData.Downloaded = false
		} else {
			metaData.Downloaded = true
		}

		_, updateErr := s.metaDataRepository.Update(metaData)
		if updateErr != nil {
			err = updateErr
		}

	}()

	img, err := s.client.GetImgByURL(metaData.DownloadLink)
	if err != nil {
		return nil, err
	}
	//@TODO обработать ситуацию, когда получаем не изображение
	webpImage, err = s.imageModifyerService.ConvertToWebp(img)
	if err != nil {
		return nil, err
	}

	err = s.fileRepository.Insert(metaData.ID, webpImage)
	if err != nil {
		return nil, err
	}

	return webpImage, nil
}
