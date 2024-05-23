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
	fmt.Println("Successfully Opened " + testFile)
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("%+v\n", jsonFile.Name())

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
