package library

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var titleHeader *gtk.Label

func CreateLibraryPage() *gtk.Box {
	libraryContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleError(err, "Cannot create box")

	titleHeaader, err := gtk.LabelNew("Home")
	utils.HandleError(err, "Cannot create label")

	libraryContainer.PackStart(titleHeaader, false, false, 0)

	return libraryContainer
}
