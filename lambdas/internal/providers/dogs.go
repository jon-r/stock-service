package providers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type DogApiRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func FetchDogItem(url string) (*DogApiRes, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("FETCH ERROR %v", err)
		return nil, err
	} else {
		log.Printf("request body: %s", string(body))
	}

	var data DogApiRes
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
