package explorer

import (
	"github.com/Pauloo27/neptune/gui/app/explorer/library"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var explorerContainer *gtk.Notebook

func CreateExplorer() *gtk.Notebook {
	var err error
	explorerContainer, err = gtk.NotebookNew()
	utils.HandleError(err, "Cannot create notebook")

	libraryLabel, err := gtk.LabelNew("Library")
	explorerContainer.AppendPage(library.CreateLibraryPage(), libraryLabel)

	searchLabel, err := gtk.LabelNew("Search")
	explorerContainer.AppendPage(createSearchPage(), searchLabel)

	return explorerContainer
}
