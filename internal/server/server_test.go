package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Saffica/image-storage/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("getImg OK", func(t *testing.T) {
		router := gin.Default()
		imgService := &mockImgService{
			Data: []byte("getImg OK"),
		}
		port := 8081
		baseAddr := fmt.Sprintf("http://localhost:%d", port)
		s := New(imgService, router)
		defer s.Stop()

		go func() {
			err := s.Run(port)
			require.NoError(t, err)
		}()

		resp, body, err := doRequest(t, http.MethodGet, fmt.Sprintf("%s/img/1", baseAddr), []byte{})
		require.NoError(t, err)
		require.Equal(t, "1", imgService.URL)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "attachment; filename=output.webp", resp.Header.Get("Content-Disposition"))
		require.Equal(t, "application/octet-stream", resp.Header.Get("Content-Type"))
		require.Equal(t, imgService.Data, body)
	})

	t.Run("getImg with bad base64 error", func(t *testing.T) {
		router := gin.Default()
		imgService := &mockImgServiceBadBase64{}
		port := 8081
		baseAddr := fmt.Sprintf("http://localhost:%d", port)
		s := New(imgService, router)
		defer s.Stop()

		go func() {
			err := s.Run(port)
			require.NoError(t, err)
		}()

		resp, _, err := doRequest(t, http.MethodGet, fmt.Sprintf("%s/img/1", baseAddr), []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("getImg internal server error", func(t *testing.T) {
		router := gin.Default()
		imgService := &mockImgServiceInternalServerError{}
		port := 8081
		baseAddr := fmt.Sprintf("http://localhost:%d", port)
		s := New(imgService, router)
		defer s.Stop()

		go func() {
			err := s.Run(port)
			require.NoError(t, err)
		}()

		resp, _, err := doRequest(t, http.MethodGet, fmt.Sprintf("%s/img/1", baseAddr), []byte{})
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func doRequest(
	t *testing.T, method, url string, body []byte,
) (resp *http.Response, respBody []byte, err error) {
	t.Helper()
	newReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, newReader)
	if err != nil {
		return nil, nil, err
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, respBody, nil
}

type mockImgService struct {
	URL  string
	Data []byte
}

func (r *mockImgService) GetImgByURL(url string) ([]byte, error) {
	r.URL = url
	return r.Data, nil
}

type mockImgServiceBadBase64 struct{}

func (r *mockImgServiceBadBase64) GetImgByURL(url string) ([]byte, error) {
	return nil, models.ErrBadBase64
}

type mockImgServiceInternalServerError struct{}

func (r *mockImgServiceInternalServerError) GetImgByURL(url string) ([]byte, error) {
	return nil, errors.New("Chto-to poshlo ne tak")
}