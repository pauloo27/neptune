package player

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

var parenthesisRegex = regexp.MustCompile(`\s?\(.+\)`)

func PlayTrack(track *db.Track) {
	RemoveCurrentFromPlaylist()
	filePath := path.Join(DataFolder, "songs", track.YoutubeID+".m4a")

	err := db.PlayTrack(track)
	utils.HandleError(err, "Cannot play track")

	State.Track = track
	LoadFile(filePath)
}

func PlayResult(result *youtube.YoutubeEntry) {
	RemoveCurrentFromPlaylist()
	State.Fetching = result
	callHooks(HOOK_RESULT_FETCH_STARTED, nil)
	filePath := path.Join(DataFolder, "songs", result.ID+".m4a")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		callHooks(HOOK_RESULT_DOWNLOAD_STARTED, nil)
		go func() {
			fmt.Println("Downloading file...")
			videoInfo, err := youtube.DownloadResult(result, filePath)
			utils.HandleError(err, "Cannot download file")
			// fix track with '(stuff)'
			trackName := parenthesisRegex.ReplaceAllString(videoInfo.Track, "")
			// fix for "artist" list (splitted by ',')
			artist := strings.Split(videoInfo.Artist, ",")[0]
			trackInfo, err := providers.FetchTrackInfo(artist, trackName)
			utils.HandleError(err, "Cannot fetch track info")
			track, err := db.StoreTrack(videoInfo, trackInfo)
			utils.HandleError(err, "Cannot store track")
			State.Track = track
			LoadFile(filePath)
		}()
	} else {
		track, err := db.PlayEntry(result)
		utils.HandleError(err, "Cannot get track")
		State.Track = track
		LoadFile(filePath)
	}
	State.Fetching = nil
}
