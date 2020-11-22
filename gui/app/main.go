package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

const (
	boxSpacing = 1
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

	baseContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, boxSpacing)

	// search bar
	searchBarContainer, err := gtk.HeaderBarNew()
	utils.HandleError(err, "Cannot create header bar")

	searchBarContainer.SetShowCloseButton(false)

	searchInput, err := gtk.EntryNew()
	utils.HandleError(err, "Cannot create entry")

	searchInput.SetPlaceholderText("Search YouTube")
	searchInput.SetHExpand(true)

	searchButton, err := gtk.ButtonNewFromIconName("search", gtk.ICON_SIZE_BUTTON)

	searchBarContainer.PackStart(searchInput)
	searchBarContainer.PackStart(searchButton)

	// main content container
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, boxSpacing)
	utils.HandleError(err, "Cannot create box")

	mainContainer.SetHomogeneous(true)
	mainContainer.PackStart(CreatePlayer(), false, true, boxSpacing)
	mainContainer.PackEnd(CreateLibrary(), false, true, boxSpacing)

	win.Add(baseContainer)

	baseContainer.PackStart(searchBarContainer, false, false, boxSpacing)
	baseContainer.PackEnd(mainContainer, true, true, boxSpacing)

	win.SetDefaultSize(800, 600)

	win.ShowAll()

	gtk.Main()
}
