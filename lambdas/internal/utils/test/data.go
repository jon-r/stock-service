package test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
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

func ReadJsonToString(testFile string) string {
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

	return string(byteValue)
}

func unpackArray(s any) []any {
	v := reflect.ValueOf(s)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}
