package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreatePlayer() *gtk.Box {
	playerContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, boxSpacing)
	utils.HandleError(err, "Cannot create box")

	l, err := gtk.LabelNew("Hello")
	utils.HandleError(err, "Cannot create label")

	playerContainer.PackStart(l, false, false, 0)

	return playerContainer
}
