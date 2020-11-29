package app

import (
	"github.com/Pauloo27/neptune/gui/app/library"
	"github.com/Pauloo27/neptune/gui/app/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var appWin *gtk.Window

func Start() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err, "Cannot create window")

	win.SetTitle("My Tune")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	appWin = win

	baseContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	// main content container
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	mainContainer.SetHomogeneous(true)
	mainContainer.PackStart(player.CreatePlayer(), false, true, 1)
	mainContainer.PackEnd(library.CreateLibrary(), false, true, 1)

	win.Add(baseContainer)

	baseContainer.PackStart(library.CreateSearchHeader(), false, false, 0)
	baseContainer.PackEnd(mainContainer, true, true, 0)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	gtk.Main()
}
