package player

import (
	"math"
	"math/rand"
	"time"

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
		false,
		[]int{},
		initialVolume,
		0.0,
	}

	// internal hooks
	RegisterHook(HOOK_FILE_LOAD_STARTED, func(params ...interface{}) {
		Play()
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

func GetTrackAt(index int) *db.Track {
	if index >= len(State.Queue) {
		return nil
	}
	if State.Shuffled {
		return State.Queue[State.ShuffIndexes[index]]
	}
	return State.Queue[index]
}

func GetCurrentTrack() *db.Track {
	return GetTrackAt(State.QueueIndex)
}

func AddToQueue(track *db.Track) {
	State.Queue = append(State.Queue, track)
	// TODO: if shuffled, update shuffindexes
}

func AddToTopOfQueue(track *db.Track) {
	// TODO: if shuffled, update shuffindexes
	newQueue := []*db.Track{track}
	newQueue = append(newQueue, State.Queue...)
	State.Queue = newQueue
}

func RemoveFromQueue(index int) {
	if index >= len(State.Queue) {
		return
	}
	newQueue := []*db.Track{}
	for i := 0; i < len(State.Queue); i++ {
		if i == index {
			continue
		}
		newQueue = append(newQueue, GetTrackAt(i))
	}
	State.Queue = newQueue
	MpvInstance.CommandString(utils.Fmt("playlist-remove %d", index))
	callHooks(HOOK_QUEUE_UPDATE_FINISHED)
	// TODO: if shuffled, update shuffindexes
}

func ClearQueue() {
	State.Queue = []*db.Track{}
	State.Shuffled = false
	State.ShuffIndexes = []int{}
	clearPlaylist()
	removeCurrentFromPlaylist()
}

func clearPlaylist() error {
	return MpvInstance.CommandString("playlist-clear")
}

func removeCurrentFromPlaylist() error {
	return MpvInstance.CommandString("playlist-remove current")
}

func LoadFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.CommandString("loadfile " + filePath)
	callHooks(HOOK_FILE_LOAD_STARTED, err, filePath)
	return err
}

func AppendFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.CommandString("loadfile " + filePath + " append")
	callHooks(HOOK_FILE_APPENDED, err, filePath)
	return err
}

func Stop() error {
	return MpvInstance.CommandString("stop")
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

func SetCurrentTrackID(id int) error {
	return MpvInstance.SetProperty("playlist-pos", mpv.FORMAT_INT64, id)
}

func Shuffle() {
	State.Shuffled = true
	State.ShuffIndexes = []int{}
	for i := 0; i < len(State.Queue); i++ {
		State.ShuffIndexes = append(State.ShuffIndexes, i)
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(State.Queue), func(i, j int) {
		State.ShuffIndexes[i], State.ShuffIndexes[j] = State.ShuffIndexes[j], State.ShuffIndexes[i]
	})
	SetCurrentTrackID(0)
}

func PreviousTrack() error {
	return MpvInstance.CommandString("playlist-prev")
}

func Exit() error {
	callHooks(HOOK_PLAYER_EXIT)
	return Stop()
}

func NextTrack() error {
	return MpvInstance.CommandString("playlist-next")
}
