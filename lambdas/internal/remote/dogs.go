package remote

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

func FetchDogItem() (*DogApiRes, error) {
	res, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
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
