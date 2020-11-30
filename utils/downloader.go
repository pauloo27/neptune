package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(fileURL, targetFilePath string) error {
	res, err := http.Get(fileURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	file, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}
