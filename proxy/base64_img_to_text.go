package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Base64ImgToText takes a base64 image source and calls the OCR service to
// extract the image text
func Base64ImgToText(base64Str string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"base64_string": getBase64StringData(base64Str),
		"extension":     getBase64StringExtension(base64Str),
	})

	if err != nil {
		return "", err
	}

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	request, err := http.NewRequest("POST", "http://localhost:8080/base64_image_text", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-type", "application/json")

	if err != nil {
		return "", err
	}

	resp, err := client.Do(request)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getBase64StringExtension(str string) string {
	idx := strings.Index(str, ";charset=utf-8;base64,")

	return strings.TrimLeft(str[0:idx], "data:image/")
}

func getBase64StringData(str string) string {
	return strings.Split(str, ";charset=utf-8;base64,")[1]
}
