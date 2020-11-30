package player

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/utils"
	"github.com/Pauloo27/neptune/youtube"
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
			err = db.StoreTrack(videoInfo, trackInfo)
			utils.HandleError(err, "Cannot store track")
			LoadFile(filePath)
			fmt.Println(trackInfo)
		}()
	} else {
		LoadFile(filePath)
	}
}
