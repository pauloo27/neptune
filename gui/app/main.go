package app

import (
	"os"

	"github.com/Pauloo27/neptune/gui/app/explorer"
	"github.com/Pauloo27/neptune/gui/app/player"
	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var appWin *gtk.Window

func Start(onExit func()) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err, "Cannot create window")

	win.SetTitle("Neptune")
	win.Connect("destroy", func() {
		onExit()
		gtk.MainQuit()
		os.Exit(0)
	})

	appWin = win

	baseContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	// main content container
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	mainContainer.SetHomogeneous(true)
	mainContainer.PackStart(player.CreatePlayer(), false, true, 1)
	mainContainer.PackEnd(explorer.CreateExplorer(), false, true, 1)

	win.Add(baseContainer)

	baseContainer.PackStart(explorer.CreateSearchHeader(), false, false, 0)
	baseContainer.PackEnd(mainContainer, true, true, 0)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	hook.RegisterHook(hook.HOOK_REQUEST_EXIT, func(params ...interface{}) {
		onExit()
		gtk.MainQuit()
		os.Exit(0)
	})

	hook.RegisterHook(hook.HOOK_REQUEST_SHOW_HIDE, func(params ...interface{}) {
		win.SetVisible(!win.GetVisible())
	})

	hook.CallHooks(hook.HOOK_GUI_STARTED)

	gtk.Main()
}
