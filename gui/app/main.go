package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func Start() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleError(err, "Cannot create window")

	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	l, err := gtk.LabelNew("Hello, gotk3!")
	utils.HandleError(err, "Cannot create label")

	win.Add(l)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	gtk.Main()
}
