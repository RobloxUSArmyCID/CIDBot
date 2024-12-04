package cidbot

import (
	"io"
	"encoding/json"
	"fmt"
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

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code - %d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseBody T

	err = json.Unmarshal(responseBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}