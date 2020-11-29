package library

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreateLibrary() *gtk.Notebook {
	libraryContainer, err := gtk.NotebookNew()
	utils.HandleError(err, "Cannot create notebook")

	searchLabel, err := gtk.LabelNew("Search")
	libraryContainer.AppendPage(createSearchPage(), searchLabel)

	return libraryContainer
}
