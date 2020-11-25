package app

import (
	"github.com/Pauloo27/my-tune/player"
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func CreatePlayer() *gtk.Box {
	playerContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	buttonsContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	playPauseButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")
	playPauseButton.SetHAlign(gtk.ALIGN_CENTER)

	playPauseButton.Connect("clicked", func() {
		player.PlayPause()
	})

	buttonsContainer.PackStart(playPauseButton, true, true, 0)

	playerContainer.PackEnd(buttonsContainer, false, false, 0)

	return playerContainer
}
