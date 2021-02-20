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

func addToQueue(track *db.Track) {
	State.Queue = append(State.Queue, track)
	// TODO: if shuffled, update shuffindexes
}

func addToTopOfQueue(track *db.Track) {
	// TODO: if shuffled, update shuffindexes
	newQueue := []*db.Track{track}
	newQueue = append(newQueue, State.Queue...)
	State.Queue = newQueue
}

func RemoveFromQueue(index int) error {
	if index >= len(State.Queue) {
		return nil
	}
	newQueue := []*db.Track{}
	for i := 0; i < len(State.Queue); i++ {
		if i == index {
			continue
		}
		newQueue = append(newQueue, GetTrackAt(i))
	}
	State.Queue = newQueue
	err := MpvInstance.CommandString(utils.Fmt("playlist-remove %d", index))
	if err != nil {
		return err
	}
	callHooks(HOOK_QUEUE_UPDATE_FINISHED)
	// TODO: if shuffled, update shuffindexes
	return nil
}

func ClearQueue() error {
	err := clearQueue()
	if err != nil {
		return err
	}
	callHooks(HOOK_QUEUE_UPDATE_FINISHED)
	return nil
}

func clearQueue() error {
	State.Queue = []*db.Track{}
	State.Shuffled = false
	State.ShuffIndexes = []int{}
	return clearEntirePlaylist()
}

func clearEntirePlaylist() error {
	err := trimPlaylist()
	if err != nil {
		return err
	}
	return removeCurrentFromPlaylist()
}

func trimPlaylist() error {
	return MpvInstance.Command([]string{"playlist-clear"})
}

func removeCurrentFromPlaylist() error {
	return MpvInstance.Command([]string{"playlist-remove", "current"})
}

func loadFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath})
	callHooks(HOOK_FILE_LOAD_STARTED, err, filePath)
	return err
}

func appendFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath, "append"})
	callHooks(HOOK_FILE_APPENDED, err, filePath)
	return err
}

func appendFileAndPlay(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath, "append-play"})
	callHooks(HOOK_FILE_APPENDED, err, filePath)
	return err
}

func Stop() error {
	return MpvInstance.Command([]string{"stop"})
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

func setCurrentTrackID(id int) error {
	return MpvInstance.SetProperty("playlist-pos", mpv.FORMAT_INT64, id)
}

func Shuffle() error {
	State.ShuffIndexes = []int{}
	for i := 0; i < len(State.Queue); i++ {
		State.ShuffIndexes = append(State.ShuffIndexes, i)
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(State.Queue), func(i, j int) {
		State.ShuffIndexes[i], State.ShuffIndexes[j] = State.ShuffIndexes[j], State.ShuffIndexes[i]
	})

	err := clearEntirePlaylist()
	if err != nil {
		return err
	}

	State.Shuffled = true
	for i := 0; i < len(State.Queue); i++ {
		if i == 0 {
			err := appendFileAndPlay(GetTrackAt(i).GetPath())
			if err != nil {
				return err
			}
		} else {
			err := appendFile(GetTrackAt(i).GetPath())
			if err != nil {
				return err
			}
		}
	}

	callHooks(HOOK_QUEUE_UPDATE_FINISHED)
	return nil
}

func PreviousTrack() error {
	return MpvInstance.Command([]string{"playlist-prev"})
}

func Exit() error {
	callHooks(HOOK_PLAYER_EXIT)
	return MpvInstance.CommandString("exit 0")
}

func NextTrack() error {
	return MpvInstance.Command([]string{"playlist-next"})
}
