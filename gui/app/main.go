package app

import (
	"github.com/Pauloo27/my-tune/utils"
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

	baseContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)

	// main content container
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	mainContainer.SetHomogeneous(true)
	mainContainer.PackStart(CreatePlayer(), false, true, 1)
	mainContainer.PackEnd(CreateLibrary(), false, true, 1)

	win.Add(baseContainer)

	baseContainer.PackStart(CreateSearchHeader(), false, false, 1)
	baseContainer.PackEnd(mainContainer, true, true, 1)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	gtk.Main()
}
