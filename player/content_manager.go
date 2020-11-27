package player

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/my-tune/youtube"
)

func PlayResult(result *youtube.YoutubeEntry) {
	State.Playing = result
	callHooks(HOOK_RESULT_FETCH_STARTED, nil)
	filePath := path.Join(DataFolder, "songs", result.ID+".m4a")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// TODO: Download
		fmt.Println("Downloading file...")
		os.Exit(-1)
	}
	LoadFile(filePath)
}
