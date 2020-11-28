package player

import (
	"fmt"

	"github.com/Pauloo27/my-tune/player/mpv"
	"github.com/Pauloo27/my-tune/utils"
)

func startEventHandler() {
	go func() {
		for {
			event := MpvInstance.WaitEvent(60)
			switch event.Event_Id {
			case mpv.EVENT_NONE:
				continue
			case mpv.EVENT_PROPERTY_CHANGE:
				data := event.Data.(*mpv.EventProperty)
				fmt.Println(data.Name)
			case mpv.EVENT_FILE_LOADED:
				duration, err := MpvInstance.GetProperty("duration", mpv.FORMAT_DOUBLE)
				utils.HandleError(err, "Cannot get duration")
				State.Duration = duration.(float64)
				callHooks(HOOK_FILE_LOADED, err, duration)
			case mpv.EVENT_PAUSE:
				State.Paused = true
				callHooks(HOOK_PLAYBACK_PAUSED)
			case mpv.EVENT_UNPAUSE:
				State.Paused = false
				callHooks(HOOK_PLAYBACK_RESUMED)
			default:
				fmt.Println(event)
			}
		}
	}()
}
