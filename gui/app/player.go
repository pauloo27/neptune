package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func createPlayer() *gtk.Box {
	playerContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	//playerContainer.PackEnd(bottomContainer, false, false, 1)

	return playerContainer
}
