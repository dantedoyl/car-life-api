package filesystem

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)


func InsertPhoto(fileHeader *multipart.FileHeader, photoPath string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	str, err := os.Getwd()
	if err != nil {
		return "", err
	}

	os.Chdir(photoPath)

	photoID := rand.Uint64()
	photoIDStr := strconv.FormatUint(photoID, 10)

	extension := filepath.Ext(fileHeader.Filename)

	newFile, err := os.OpenFile(photoIDStr+extension, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	os.Chdir(str)

	_, err = io.Copy(newFile, file)
	if err != nil {
		_ = os.Remove(photoIDStr + extension)
		return "", err
	}

	photo := "/" + photoPath + photoIDStr + extension
	return photo, nil
}

func InsertPhotos(filesHeaders []*multipart.FileHeader, photoPath string) ([]string, error) {
	imgUrls := make(map[string][]string)

	for i := range filesHeaders {
		url, err := InsertPhoto(filesHeaders[i], photoPath)
		if err != nil {
			return nil, err
		}

		imgUrls["img"] = append(imgUrls["img"], url)
	}

	return imgUrls["img"], nil
}

func RemovePhoto(imgUrl string) error {
	if imgUrl == "" {
		return nil
	}

	origWd, _ := os.Getwd()
	err := os.Remove(origWd + imgUrl)
	if err != nil {
		return err
	}

	return nil
}

func RemovePhotos(imgUrls []string) error {
	if len(imgUrls) == 0 {
		return nil
	}

	for _, photo := range imgUrls {
		err := RemovePhoto(photo)
		if err != nil {
			return err
		}
	}

	return nil
}