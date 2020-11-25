package player

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/Pauloo27/my-tune/youtube"
	"github.com/YouROK/go-mpv/mpv"
)

var MpvInstance *mpv.Mpv
var State *PlayerState

func Initialize() {
	var err error

	// create a mpv instance
	MpvInstance = mpv.Create()

	// set options
	// disable video
	err = MpvInstance.SetOptionString("video", "no")
	utils.HandleError(err, "Cannot set mpv video option")

	// disable cache
	err = MpvInstance.SetOptionString("cache", "no")
	utils.HandleError(err, "Cannot set mpv cache option")

	// set quality to worst
	err = MpvInstance.SetOptionString("ytdl-format", "worst")
	utils.HandleError(err, "Cannot set mpv ytdl-format option")

	// add observers
	//MpvInstance.ObserveProperty(1, "pause", mpv.FORMAT_FLAG)

	// start event listener
	startEventHandler()

	// create the state
	State = &PlayerState{
		false,
		nil,
	}

	// start the player
	MpvInstance.Initialize()

	callHooks(HOOK_PLAYER_INITIALIZED)
}

func Load(result *youtube.YoutubeEntry) error {
	State.Playing = result
	err := MpvInstance.Command([]string{"loadfile", result.URL()})
	callHooks(HOOK_FILE_LOADED)
	return err
}

func PlayPause() error {
	if State.Paused {
		return Play()
	} else {
		return Pause()
	}
}

func Pause() error {
	err := MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, true)
	callHooks(HOOK_PLAYBACK_PAUSED)
	return err
}

func Play() error {
	err := MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, false)
	callHooks(HOOK_PLAYBACK_RESUMED)
	return err
}
