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
	reqBody, err := json.Marshal(map[string]string{"base64_string": base64Str})

	if err != nil {
		return "", err
	}

	client := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}

	request, err := http.NewRequest("POST", "http://localhost:8080/base64_image_ocr", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-type", "application/json")
	resp, err := client.Do(request)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var obj map[string]string
	err = json.Unmarshal(body, &obj)

	if err != nil {
		return "", err
	}

	return obj["task_key"], nil
}

func getBase64StringExtension(str string) string {
	idx := strings.Index(str, ";charset=utf-8;base64,")

	return strings.TrimLeft(str[0:idx], "data:image/")
}

func getBase64StringData(str string) string {
	return strings.Split(str, ";charset=utf-8;base64,")[1]
}
