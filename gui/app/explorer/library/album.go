package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func displayAlbum(album *db.Album) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(album.Title)
	utils.HandleError(err, "Cannot create button")

	btn.Connect("clicked", func() {
		// TODO:
	})

	return btn
}
