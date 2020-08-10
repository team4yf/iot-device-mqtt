package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/oliveagle/jsonpath"
)

func GetBodyByte(readCloser io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

func GetBodyString(readCloser io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	defer readCloser.Close()

	return string(body), nil
}

func GetBodyMap(readCloser io.ReadCloser) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := json.NewDecoder(readCloser).Decode(&body)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

func GetBodyStruct(readCloser io.ReadCloser, obj interface{}) error {
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}
	fmt.Printf("from pkg:%v", obj)
	return nil
}

func GetJsonPathData(data, jp string) (interface{}, error) {
	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)
	res, err := jsonpath.JsonPathLookup(jsonData, jp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Sha256Encode(origin string) string {
	sum := sha256.Sum256(([]byte)(origin))
	return fmt.Sprintf("%x", sum)
}
