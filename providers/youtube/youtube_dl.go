package youtube

import (
	"os/exec"

	"github.com/buger/jsonparser"
)

type VideoInfo struct {
	Artist, Track, Uploader, Title, ID string
	Duration                           int64
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

	uploader, _ := jsonparser.GetString(buffer, "uploader")
	title, _ := jsonparser.GetString(buffer, "title")
	id, _ := jsonparser.GetString(buffer, "id")
	duration, _ := jsonparser.GetInt(buffer, "duration")
	artist, _ := jsonparser.GetString(buffer, "artist")
	track, _ := jsonparser.GetString(buffer, "track")

	return &VideoInfo{artist, track, uploader, title, id, duration}, nil
}
