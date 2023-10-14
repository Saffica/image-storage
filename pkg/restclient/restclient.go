package restclient

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	tr     = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client = &http.Client{Transport: tr, Timeout: 10 * time.Second}
)

type restClient struct {
}

func New() *restClient {
	return &restClient{}
}

func (a *restClient) GetImgByURL(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return body, nil
}
