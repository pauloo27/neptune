package player

import (
	"fmt"

	"github.com/Pauloo27/my-tune/utils"
	"github.com/YouROK/go-mpv/mpv"
)

func startEventHandler() {
	go func() {
		for {
			event := MpvInstance.WaitEvent(60)
			switch event.Event_Id {
			case mpv.EVENT_NONE:
				continue
			case mpv.EVENT_FILE_LOADED:
				duration, err := MpvInstance.GetProperty("duration", mpv.FORMAT_DOUBLE)
				utils.HandleError(err, "Cannot get duration")
				State.Duration = duration.(float64)
				callHooks(HOOK_FILE_LOADED, err, duration)
			case mpv.EVENT_PAUSE:
				State.Paused = true
			case mpv.EVENT_UNPAUSE:
				State.Paused = false
			default:
				fmt.Println(event)
			}
		}
	}()
}
