package utils

import (
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
)

var userHome string

func GetUserHome() string {
	if userHome != "" {
		return userHome
	}
	var err error
	userHome, err = os.UserHomeDir()
	HandleError(err, "Cannot get user home")
	return userHome
}

func GetTmpFolder() string {
	if runtime.GOOS == "windows" {
		return path.Join(GetUserHome(), "AppData", "Local", "Temp")
	}
	return "/tmp/"
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
