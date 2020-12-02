package library

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreateLibraryPage() *gtk.Box {
	libraryContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleError(err, "Cannot create box")

	return libraryContainer
}
