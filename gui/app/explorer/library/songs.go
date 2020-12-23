package library

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var (
	songsPage = &LibraryPage{"Songs", showSongs}
)

func showSongs() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	return container
}
