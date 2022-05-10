package filesystem

import (
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

func InsertPhoto(fileHeader *multipart.FileHeader, dirname string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	filename := strings.Replace(uuid.New().String(), "-", "", -1) + ".webp"
	filepath := dirname + filename

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return "", err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 25})
	if err != nil {
		return "", err
	}

	err = bimg.Write("./"+filepath, processed)
	if err != nil {
		return "", err
	}

	return "/" + filepath, nil
}

func InsertPhotos(filesHeaders []*multipart.FileHeader, dirname string) ([]string, error) {
	var imgUrls []string

	for _, fileHeader := range filesHeaders {
		imgUrl, err := InsertPhoto(fileHeader, dirname)
		if err != nil {
			return nil, err
		}

		imgUrls = append(imgUrls, imgUrl)
	}

	return imgUrls, nil
}

func RemovePhoto(imgUrl string) error {
	if imgUrl == "" {
		return nil
	}

	err := os.Remove("./" + imgUrl)
	if err != nil {
		return err
	}

	return nil
}

func RemovePhotos(imgUrls []string) error {
	if len(imgUrls) == 0 {
		return nil
	}

	for _, imgUrl := range imgUrls {
		err := RemovePhoto(imgUrl)
		if err != nil {
			return err
		}
	}

	return nil
}
