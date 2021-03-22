package trayicon

import (
	"fmt"

	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/trayicon/icon"
	"github.com/getlantern/systray"
)

func LoadTrayIcon() {
	systray.Run(onReady, onExit)
}

func onExit() {
	// TODO
}

func onReady() {
	fmt.Println("Ready")
	systray.SetTitle("Neptune")
	systray.SetTooltip("Neptune")
	systray.SetIcon(icon.ICON_DATA)

	mQuit := systray.AddMenuItem("Quit", "Quit")
	mShowHide := systray.AddMenuItem("Show/Hide player", "Show/Hide player")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				hook.CallHooks(hook.HOOK_REQUEST_EXIT)
			case <-mShowHide.ClickedCh:
				hook.CallHooks(hook.HOOK_REQUEST_SHOW_HIDE)
			}
		}
	}()
}
