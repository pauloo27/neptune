package youtube

import (
	"os/exec"

	"github.com/Pauloo27/neptune/utils"
	"github.com/buger/jsonparser"
)

type VideoInfo struct {
	Artist, Track, Uploader, UploaderID, Title, ID string
	Duration                                       int64
}

func (v *VideoInfo) GetThumbnail() string {
	return utils.Fmt(
		"https://i1.ytimg.com/vi/%s/hqdefault.jpg", v.ID,
	)
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
	uploaderID, _ := jsonparser.GetString(buffer, "uploader_id")
	title, _ := jsonparser.GetString(buffer, "title")
	id, _ := jsonparser.GetString(buffer, "id")
	duration, _ := jsonparser.GetInt(buffer, "duration")
	artist, _ := jsonparser.GetString(buffer, "artist")
	track, _ := jsonparser.GetString(buffer, "track")

	return &VideoInfo{artist, track, uploader, uploaderID, title, id, duration}, nil
}
