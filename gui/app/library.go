package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreateLibrary() *gtk.Notebook {
	libraryContainer, err := gtk.NotebookNew()
	utils.HandleError(err, "Cannot create notebook")

	searchLabel, err := gtk.LabelNew("Search")
	libraryContainer.AppendPage(CreateSearchPage(), searchLabel)

	return libraryContainer
}
