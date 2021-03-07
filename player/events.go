package player

import (
	"fmt"
	"unsafe"

	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/player/mpv"
	"github.com/Pauloo27/neptune/utils"
)

func handlePropertyChange(data *mpv.EventProperty) {
	switch data.Name {
	case "volume":
		volume := *(*float64)(data.Data.(unsafe.Pointer))
		hook.CallHooks(hook.HOOK_VOLUME_CHANGED, volume)
	case "loop-file":
		if data.Data.(unsafe.Pointer) == nil {
			State.loopFile = true
		} else {
			State.loopFile = *(*bool)(data.Data.(unsafe.Pointer))
		}
		hook.CallHooks(hook.HOOK_LOOP_STATUS_CHANGED)
	case "loop-playlist":
		if data.Data.(unsafe.Pointer) == nil {
			State.loopPlaylist = true
		} else {
			State.loopPlaylist = *(*bool)(data.Data.(unsafe.Pointer))
		}
		hook.CallHooks(hook.HOOK_LOOP_STATUS_CHANGED)
	default:
		fmt.Printf("Property %s changed\n", data.Name)
	}
}

func startEventHandler() {
	go func() {
		for {
			event := MpvInstance.WaitEvent(60)
			switch event.Event_Id {
			case mpv.EVENT_NONE:
				continue
			case mpv.EVENT_PROPERTY_CHANGE:
				data := event.Data.(*mpv.EventProperty)
				handlePropertyChange(data)
			case mpv.EVENT_FILE_LOADED:
				duration, err := MpvInstance.GetProperty("duration", mpv.FORMAT_DOUBLE)
				utils.HandleError(err, "Cannot get duration")
				State.Duration = duration.(float64)
				hook.CallHooks(hook.HOOK_FILE_LOADED, err, duration)
			case mpv.EVENT_PAUSE:
				State.Paused = true
				hook.CallHooks(hook.HOOK_PLAYBACK_PAUSED)
			case mpv.EVENT_END_FILE:
				hook.CallHooks(hook.HOOK_FILE_ENDED)
			case mpv.EVENT_UNPAUSE:
				State.Paused = false
				hook.CallHooks(hook.HOOK_PLAYBACK_RESUMED)
			default:
				fmt.Println(event)
			}
		}
	}()
}
