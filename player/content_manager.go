package player

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func PlayResult(result *youtube.YoutubeEntry) {
	RemoveCurrentFromPlaylist()
	State.Playing = result
	callHooks(HOOK_RESULT_FETCH_STARTED, nil)
	filePath := path.Join(DataFolder, "songs", result.ID+".m4a")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		callHooks(HOOK_RESULT_DOWNLOAD_STARTED, nil)
		go func() {
			fmt.Println("Downloading file...")
			videoInfo, err := youtube.DownloadResult(result, filePath)
			utils.HandleError(err, "Cannot download file")
			trackInfo, err := providers.FetchTrackInfo(videoInfo.Artist, videoInfo.Track)
			utils.HandleError(err, "Cannot fetch track info")
			track, err := db.StoreTrack(videoInfo, trackInfo)
			utils.HandleError(err, "Cannot store track")
			State.Track = track
			LoadFile(filePath)
		}()
	} else {
		track, err := db.TrackFrom(result)
		utils.HandleError(err, "Cannot get track")
		State.Track = track
		LoadFile(filePath)
	}
}
