package youtubedl

import (
	"os/exec"

	"github.com/Pauloo27/my-tune/youtube"
)

var youtubeDLPath = ""

func GetYouTubeDLPath() string {
	return "/usr/bin/youtube-dl"
}

func DownloadResult(result *youtube.YoutubeEntry, filePath string) error {
	if youtubeDLPath == "" {
		youtubeDLPath = GetYouTubeDLPath()
	}

	cmd := exec.Command(youtubeDLPath, result.URL(),
		"-f 140", "--add-metadata", "-o", filePath,
	)

	err := cmd.Run()

	return err
}
