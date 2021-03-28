package youtube

import (
	"os/exec"
	"runtime"

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

func parseVideoInfo(buffer []byte) *VideoInfo {
	uploader, _ := jsonparser.GetString(buffer, "uploader")
	uploaderID, _ := jsonparser.GetString(buffer, "uploader_id")
	title, _ := jsonparser.GetString(buffer, "title")
	id, _ := jsonparser.GetString(buffer, "id")
	duration, _ := jsonparser.GetInt(buffer, "duration")
	artist, _ := jsonparser.GetString(buffer, "artist")
	track, _ := jsonparser.GetString(buffer, "track")

	return &VideoInfo{artist, track, uploader, uploaderID, title, id, duration}
}

var youtubeDLPath string

func GetYouTubeDLPath() string {
	if youtubeDLPath != "" {
		return youtubeDLPath
	}
	if runtime.GOOS == "windows" {
		youtubeDLPath = "bin/youtube-dl.exe"
	} else {
		youtubeDLPath = "youtube-dl"
	}
	return youtubeDLPath
}

func FetchInfoAndDownload(result *YoutubeEntry, filePath string) (*VideoInfo, error) {
	cmd := exec.Command(GetYouTubeDLPath(), result.URL(),
		"-f 140", "--add-metadata", "-o", filePath, "--print-json",
	)

	buffer, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseVideoInfo(buffer), nil
}

func FetchInfo(result *YoutubeEntry) (*VideoInfo, error) {
	cmd := exec.Command(GetYouTubeDLPath(), result.URL(), "--dump-json")

	buffer, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseVideoInfo(buffer), nil
}
