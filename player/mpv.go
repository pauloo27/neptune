package player

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/Pauloo27/my-tune/youtube"
	"github.com/YouROK/go-mpv/mpv"
)

var mpvInstance *mpv.Mpv

func Initialize() {
	var err error

	// create a mpv instance
	mpvInstance = mpv.Create()

	// set options
	// disable video
	err = mpvInstance.SetOptionString("video", "no")
	utils.HandleError(err, "Cannot set mpv video option")

	// disable cache
	err = mpvInstance.SetOptionString("cache", "no")
	utils.HandleError(err, "Cannot set mpv cache option")

	// set quality to worst
	err = mpvInstance.SetOptionString("ytdl-format", "worst")
	utils.HandleError(err, "Cannot set mpv ytdl-format option")

	// start the player
	mpvInstance.Initialize()

	callHooks(HOOK_PLAYER_INITIALIZED)
}

func Play(result *youtube.YouTubeResult) error {
	err := mpvInstance.Command([]string{"loadfile", result.URL()})
	callHooks(HOOK_FILE_LOADED)
	return err
}
