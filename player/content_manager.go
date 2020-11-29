package player

import (
	"fmt"
	"os"
	"path"

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
			info, err := youtube.DownloadResult(result, filePath)
			utils.HandleError(err, "Cannot download file")
			LoadFile(filePath)
			fmt.Println("artist:", info.Artist, "track", info.Track)
		}()
	} else {
		LoadFile(filePath)
	}
}
