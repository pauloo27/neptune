package player

import (
	"math"
	"math/rand"
	"time"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/hook"
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
	err = MpvInstance.ObserveProperty(0, "loop-file", mpv.FORMAT_FLAG)
	utils.HandleError(err, "Cannot observer loop-file property")
	err = MpvInstance.ObserveProperty(0, "loop-playlist", mpv.FORMAT_FLAG)
	utils.HandleError(err, "Cannot observer loop-playlist property")

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
		false, false,
	}

	// internal hooks
	hook.RegisterHook(hook.HOOK_FILE_LOAD_STARTED, func(params ...interface{}) {
		Play()
	})
	hook.RegisterHook(hook.HOOK_VOLUME_CHANGED, func(params ...interface{}) {
		State.Volume = params[0].(float64)
	})
	hook.RegisterHook(hook.HOOK_FILE_ENDED, func(params ...interface{}) {
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

	hook.CallHooks(hook.HOOK_PLAYER_INITIALIZED, err)
}

func PlayTrack(track *db.Track) {
	clearQueue()

	addToTopOfQueue(track)
	loadFile(track.GetPath())

	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func PlayTracks(tracks []*db.Track) {
	clearQueue()

	if len(tracks) == 0 {
		return
	}

	addToTopOfQueue(tracks[0])
	loadFile(tracks[0].GetPath())

	for _, track := range tracks[1:] {
		addToQueue(track)
		appendFile(track.GetPath())
	}
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func AddTrackToQueue(track *db.Track) {
	addToQueue(track)
	appendFile(track.GetPath())
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func GetTrackAt(index int) *db.Track {
	if index >= len(State.Queue) {
		return nil
	}
	return State.Queue[index]
}

func GetCurrentTrack() *db.Track {
	return GetTrackAt(State.QueueIndex)
}

func addToQueue(track *db.Track) {
	State.Queue = append(State.Queue, track)
}

func addToTopOfQueue(track *db.Track) {
	newQueue := []*db.Track{track}
	newQueue = append(newQueue, State.Queue...)
	State.Queue = newQueue
}

func moveInQueue(index int, up bool) error {
	offset := 1
	if up {
		offset = -1
	}
	State.Queue[index+offset], State.Queue[index] = State.Queue[index], State.Queue[index+offset]

	var from, to int
	if up {
		from, to = index, index+offset
	} else {
		to, from = index, index+offset
	}
	return MpvInstance.CommandString(utils.Fmt("playlist-move %d %d", from, to))
}

func MoveUpInQueue(index int) error {
	if index == 0 {
		return nil
	}
	err := moveInQueue(index, true)
	if err != nil {
		return err
	}
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
	return nil
}

func MoveDownInQueue(index int) error {
	if index >= len(State.Queue)-1 {
		return nil
	}
	err := moveInQueue(index, false)
	if err != nil {
		return err
	}
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
	return nil
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
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
	return nil
}

func ClearQueue() error {
	err := clearQueue()
	if err != nil {
		return err
	}
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
	return nil
}

func clearQueue() error {
	State.Queue = []*db.Track{}
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
	hook.CallHooks(hook.HOOK_FILE_LOAD_STARTED, err, filePath)
	return err
}

func appendFile(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath, "append"})
	hook.CallHooks(hook.HOOK_FILE_APPENDED, err, filePath)
	return err
}

func appendFileAndPlay(filePath string) error {
	loadMPRIS()
	err := MpvInstance.Command([]string{"loadfile", filePath, "append-play"})
	hook.CallHooks(hook.HOOK_FILE_APPENDED, err, filePath)
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
	hook.CallHooks(hook.HOOK_POSITION_CHANGED, err, pos)
	return err
}

func setCurrentTrackID(id int) error {
	return MpvInstance.SetProperty("playlist-pos", mpv.FORMAT_INT64, id)
}

func NextLoopStatus() error {
	newLoopStatus := LOOP_NONE
	switch GetLoopStatus() {
	case LOOP_NONE:
		newLoopStatus = LOOP_TRACK
	case LOOP_TRACK:
		newLoopStatus = LOOP_QUEUE
	}
	return SetLoopStatus(newLoopStatus)
}

func SetLoopStatus(loopStatus LoopStatus) error {
	track, queue := false, false

	if loopStatus == LOOP_TRACK {
		track = true
	}
	if loopStatus == LOOP_QUEUE {
		queue = true
	}

	err := MpvInstance.SetProperty("loop-file", mpv.FORMAT_FLAG, track)
	if err != nil {
		return err
	}
	return MpvInstance.SetProperty("loop-playlist", mpv.FORMAT_FLAG, queue)
}

func GetLoopStatus() LoopStatus {
	if State.loopFile {
		return LOOP_TRACK
	}
	if State.loopPlaylist {
		return LOOP_QUEUE
	}
	return LOOP_NONE
}

func Shuffle() {
	newQueue := State.Queue

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(State.Queue), func(i, j int) {
		newQueue[i], newQueue[j] = newQueue[j], newQueue[i]
	})

	PlayTracks(newQueue)
}

func PreviousTrack() error {
	return MpvInstance.Command([]string{"playlist-prev"})
}

func Exit() error {
	hook.CallHooks(hook.HOOK_PLAYER_EXIT)
	return MpvInstance.CommandString("exit 0")
}

func NextTrack() error {
	return MpvInstance.Command([]string{"playlist-next"})
}
