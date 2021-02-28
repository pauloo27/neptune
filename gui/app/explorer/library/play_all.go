package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func createPlayAll(labelText string, tracks []*db.Track) *gtk.Box {
	return createPlayAllContainer(labelText, func() { player.PlayTracks(tracks) })
}

func createFuturePlayAll(labelText string, loader func() []*db.Track) *gtk.Box {
	return createPlayAllContainer(labelText, func() { player.PlayTracks(loader()) })
}

func createPlayAllContainer(labelText string, onClick func()) *gtk.Box {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	hbox.SetMarginStart(5)
	hbox.SetMarginEnd(5)

	playButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create label")

	playButton.Connect("clicked", onClick)

	label, err := gtk.LabelNew(labelText)
	utils.HandleError(err, "Cannot create label")

	hbox.PackStart(playButton, false, false, 1)
	hbox.PackStart(label, false, false, 1)
	return hbox
}
