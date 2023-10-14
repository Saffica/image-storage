package usecase

import (
	"encoding/base64"
	"fmt"

	"github.com/chai2010/webp"
)

type Client interface {
	GetImgByURL(url string) ([]byte, error)
}

type imgService struct {
	client Client
}

func New(client Client) *imgService {
	return &imgService{
		client: client,
	}
}

func (s *imgService) GetImgByURL(url string) ([]byte, error) {
	newUrl, err := base64.StdEncoding.DecodeString(url)
	if err != nil {
		return nil, err
	}

	data, err := s.client.GetImgByURL(string(newUrl))
	if err != nil {
		return nil, err
	}

	width, height, _, err := webp.GetInfo(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("width = %d, height = %d\n", width, height)
	return data, nil
}
