package player

import (
	"math"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/player/mpv"
	"github.com/Pauloo27/neptune/utils"
)

const (
	maxVolume = 150.0
)

var MpvInstance *mpv.Mpv
var State *PlayerState
var DataFolder string

func Initialize(dataFolder string) {
	var err error

	DataFolder = dataFolder

	initialVolume := 50.0

	// create a mpv instance
	MpvInstance = mpv.Create()

	// set options
	// disable video
	err = MpvInstance.SetOptionString("video", "no")
	utils.HandleError(err, "Cannot set mpv video option")

	// disable cache
	err = MpvInstance.SetOptionString("cache", "no")
	utils.HandleError(err, "Cannot set mpv cache option")

	// set default volume value
	err = MpvInstance.SetOption("volume", mpv.FORMAT_DOUBLE, initialVolume)
	utils.HandleError(err, "Cannot set mpv volume option")

	// set default volume value
	err = MpvInstance.SetOption("volume-max", mpv.FORMAT_DOUBLE, maxVolume)
	utils.HandleError(err, "Cannot set mpv volume-max option")

	// set quality to worst
	err = MpvInstance.SetOptionString("ytdl-format", "worst")
	utils.HandleError(err, "Cannot set mpv ytdl-format option")

	// add observers
	err = MpvInstance.ObserveProperty(0, "volume", mpv.FORMAT_DOUBLE)
	utils.HandleError(err, "Cannot observer volume property")

	// start event listener
	startEventHandler()

	// create the state
	State = &PlayerState{
		false,
		nil,
		nil,
		0,
		initialVolume,
		0.0,
	}

	// internal hooks
	RegisterHook(HOOK_FILE_LOAD_STARTED, func(params ...interface{}) {
		Play()
		ClearPlaylist()
	})
	RegisterHook(HOOK_VOLUME_CHANGED, func(params ...interface{}) {
		State.Volume = params[0].(float64)
	})
	RegisterHook(HOOK_FILE_ENDED, func(params ...interface{}) {
		index, err := MpvInstance.GetProperty("playlist-pos", mpv.FORMAT_INT64)
		if err != nil {
			utils.HandleError(err, "Cannot get playlist-pos")
		}
		State.QueueIndex = int(index.(int64))
		if State.QueueIndex == -1 {
			State.QueueIndex = 0
		}
	})

	// start the player
	err = MpvInstance.Initialize()
	utils.HandleError(err, "Cannot initialize mpv")

	callHooks(HOOK_PLAYER_INITIALIZED, err)
}

func GetCurrentTrack() *db.Track {
	if len(State.Queue) == 0 {
		return nil
	}
	return State.Queue[State.QueueIndex]
}

func AddToQueue(track *db.Track) {
	State.Queue = append(State.Queue, track)
}

func AddToTopOfQueue(track *db.Track) {
	newQueue := []*db.Track{track}
	newQueue = append(newQueue, State.Queue...)
	State.Queue = newQueue
}

func ClearQueue() {
	State.Queue = []*db.Track{}
}

func ClearPlaylist() error {
	return MpvInstance.Command([]string{"playlist-clear"})
}

func RemoveCurrentFromPlaylist() error {
	return MpvInstance.Command([]string{"playlist-remove", "current"})
}

func LoadFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath})
	callHooks(HOOK_FILE_LOAD_STARTED, err, filePath)
	return err
}

func AppendFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath, "append"})
	callHooks(HOOK_FILE_APPENDED, err, filePath)
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
	return MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, true)
}

func Play() error {
	return MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, false)
}

func SetVolume(volume float64) error {
	volume = math.Min(maxVolume, volume)
	err := MpvInstance.SetProperty("volume", mpv.FORMAT_DOUBLE, volume)
	return err
}

func GetPosition() (float64, error) {
	position, err := MpvInstance.GetProperty("time-pos", mpv.FORMAT_DOUBLE)
	if err != nil {
		return 0.0, err
	}
	return position.(float64), err
}

func SetPosition(pos float64) error {
	err := MpvInstance.SetProperty("time-pos", mpv.FORMAT_DOUBLE, pos)
	callHooks(HOOK_POSITION_CHANGED, err, pos)
	return err
}
