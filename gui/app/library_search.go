package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreateSearchPage() *gtk.Label {
	l, err := gtk.LabelNew("Search page")
	utils.HandleError(err, "Cannot create label")

	return l
}
