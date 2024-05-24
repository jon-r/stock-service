package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadTestJson(testFile string, target any) {
	jsonFile, err := os.Open(testFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
