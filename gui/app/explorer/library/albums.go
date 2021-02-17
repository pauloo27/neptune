package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var albumsPage = &LibraryPage{"Albums", showAlbums}

func displayAlbum(album *db.Album, showArtistName bool) *gtk.Button {
	var displayTitle string
	if showArtistName {
		displayTitle = utils.Fmt("%s: %s", album.Artist.Name, album.Title)
	} else {
		displayTitle = album.Title
	}
	btn, err := gtk.ButtonNewWithLabel(displayTitle)
	utils.HandleError(err, "Cannot create button")

	btn.Connect("clicked", func() {
		displayPage(createAlbumPage(album))
	})

	return btn
}
func showAlbums() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	container.SetRowSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	go func() {
		albums, err := db.ListAlbums()
		utils.HandleError(err, "Cannot list albums")

		glib.IdleAdd(func() {
			for i, album := range albums {
				container.Attach(displayAlbum(album, true), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

	return container
}
