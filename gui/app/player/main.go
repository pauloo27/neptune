package player

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreatePlayer() *gtk.Box {
	playerContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleError(err, "Cannot create box")

	playerContainer.PackStart(createPlayerTop(), true, true, 0)
	playerContainer.PackEnd(createPlayerBottom(), false, false, 0)

	return playerContainer
}
