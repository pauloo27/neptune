package utils

import (
	"io"
	"net/http"
	"os"
	"path"
)

var userHome string

func GetUserHome() string {
	if userHome == "" {
		var err error
		userHome, err = os.UserHomeDir()
		HandleError(err, "Cannot get user home")
	}
	return userHome
}

var tmpFolder string

func GetTmpFolder() string {
	if tmpFolder == "" {
		tmpFolder = path.Join(GetUserHome(), ".cache", "neptune", "tmp")
		err := os.MkdirAll(tmpFolder, 0744)
		HandleError(err, "Cannot create tmp folder "+tmpFolder)
	}
	return tmpFolder
}

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
