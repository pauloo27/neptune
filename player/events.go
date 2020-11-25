package player

import (
	"fmt"

	"github.com/YouROK/go-mpv/mpv"
)

func startEventHandler() {
	go func() {
		for {
			event := MpvInstance.WaitEvent(1)
			switch event.Event_Id {
			case mpv.EVENT_NONE:
				continue
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
