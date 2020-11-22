package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

const (
	boxSpacing = 5
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

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, boxSpacing)
	utils.HandleError(err, "Cannot create box")

	win.Add(mainContainer)

	mainContainer.SetHomogeneous(true)
	mainContainer.PackStart(CreatePlayer(), false, true, boxSpacing)
	mainContainer.PackEnd(CreateLibrary(), false, true, boxSpacing)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	gtk.Main()
}
