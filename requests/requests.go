package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseData[T any] struct {
	Data []T `json:"data"`
}

func GetRequest[T any](url string) (*T, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code - %d\n\ncontent: %s\n\nendpoint: %s\n\nheaders: %s", response.StatusCode, string(responseBytes), url, response.Header)
	}

	var responseBody T

	err = json.Unmarshal(responseBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func PostRequest[T any](url string, data interface{}) (*T, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code - %d\n\ncontent: %s\n\nendpoint: %s\n\nheaders: %s", response.StatusCode, string(responseBytes), url, response.Header)
	}

	var responseBody T
	err = json.Unmarshal(responseBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}
