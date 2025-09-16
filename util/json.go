package util

import (
	"encoding/json"
)

func ToJson(obj any) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func ToObject[T any](str string) (T, error) {
	data := []byte(str)
	return ByteToObject[T](data)
}

func ByteToObject[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}
