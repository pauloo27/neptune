package youtube

import (
	"os/exec"

	"github.com/buger/jsonparser"
)

type VideoInfo struct {
	Artist, Track string
}

var youtubeDLPath = ""

func GetYouTubeDLPath() string {
	return "/usr/bin/youtube-dl"
}

func DownloadResult(result *YoutubeEntry, filePath string) (*VideoInfo, error) {
	if youtubeDLPath == "" {
		youtubeDLPath = GetYouTubeDLPath()
	}

	cmd := exec.Command(youtubeDLPath, result.URL(),
		"-f 140", "--add-metadata", "-o", filePath, "--print-json",
	)

	buffer, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	artist, err := jsonparser.GetString(buffer, "artist")
	if err != nil {
		return nil, err
	}
	track, err := jsonparser.GetString(buffer, "track")
	if err != nil {
		return nil, err
	}

	return &VideoInfo{artist, track}, nil
}
