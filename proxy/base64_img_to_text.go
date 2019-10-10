package proxy

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/otiai10/gosseract"

	"github.com/kardianos/osext"
	"github.com/teris-io/shortid"
)

// Base64ImgToText takes a base64 image source and attempts to convert it to text
func Base64ImgToText(filePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(filePath)
	return client.Text()
}

var (
	errInvalidImage = errors.New("invalid image")
)

func getImgTypeFromBase64String(str string) string {
	idx := strings.Index(str, ";charset=utf-8;base64,")

	return strings.TrimLeft(str[0:idx], "data:image/")
}

func SaveBase64ImageToDisk(base64Str string) (fileName string, err error) {
	idx := strings.Index(base64Str, ";base64,")

	if idx < 0 {
		return fileName, errInvalidImage
	}

	imgType := getImgTypeFromBase64String(base64Str[0 : idx+8])
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str[idx+8:]))
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(reader)

	if err != nil {
		return fileName, err
	}

	fileName = genImgFilePath(imgType)
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(fmt.Sprintf("unable to open file at %v", fileName))
	}

	switch imgType {
	case "png":
		img, err := png.Decode(bytes.NewReader(buff.Bytes()))

		if err != nil {
			return fileName, err
		}

		png.Encode(f, img)
	case "jpeg":
		im, err := jpeg.Decode(bytes.NewReader(buff.Bytes()))

		if err != nil {
			return fileName, err
		}

		jpeg.Encode(f, im, nil)
	default:
		return fileName, fmt.Errorf("invalid image type :%v", imgType)
	}

	return fileName, err
}

func genImgFilePath(ext string) string {
	id, _ := shortid.Generate()
	folderPath, err := osext.ExecutableFolder()

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%v/images/%v.%v", folderPath, id, ext)
}

func DeleteImage(filePath string) {
	os.Remove(filePath)
}
