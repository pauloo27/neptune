package youtube

import (
	"os/exec"
)

var youtubeDLPath = ""

func GetYouTubeDLPath() string {
	return "/usr/bin/youtube-dl"
}

func DownloadResult(result *YoutubeEntry, filePath string) error {
	if youtubeDLPath == "" {
		youtubeDLPath = GetYouTubeDLPath()
	}

	cmd := exec.Command(youtubeDLPath, result.URL(),
		"-f 140", "--add-metadata", "-o", filePath,
	)

	err := cmd.Run()

	return err
}
