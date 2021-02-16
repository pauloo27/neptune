package player

import (
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createPlayerTop() *gtk.Box {
	topContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleError(err, "Cannot create box")

	// album art
	albumArt, err := gtk.ImageNew()
	utils.HandleError(err, "Cannot create image")

	albumArt.SetVAlign(gtk.ALIGN_CENTER)

	player.RegisterHook(player.HOOK_FILE_LOADED, func(params ...interface{}) {
		imagePath := player.GetCurrentTrack().Album.GetAlbumArtPath()
		glib.IdleAdd(func() {
			imagePix, err := gdk.PixbufNewFromFileAtScale(imagePath, 300, 300, true)
			utils.HandleError(err, "Cannot load image from file")
			albumArt.SetFromPixbuf(imagePix)
		})
	})

	topContainer.PackStart(albumArt, true, true, 1)

	return topContainer
}
