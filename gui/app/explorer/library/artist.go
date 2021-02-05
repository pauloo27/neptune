package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createArtistPage(artist *db.Artist) *LibraryPage {
	show := func() *gtk.Grid {
		container, err := gtk.GridNew()
		utils.HandleError(err, "Cannot create box")

		container.SetRowSpacing(5)
		container.SetColumnHomogeneous(true)
		container.SetMarginStart(5)
		container.SetMarginEnd(5)

		albumsLabel, err := gtk.LabelNew("Albums:")
		utils.HandleError(err, "Cannot create label")

		container.Attach(albumsLabel, 0, 0, 1, 1)

		go func() {
			albums, err := db.ListAlbumsBy(artist, 1)
			utils.HandleError(err, "Cannot list albums by artist "+artist.MBID)

			glib.IdleAdd(func() {
				for i, album := range albums {
					container.Attach(displayAlbum(album), 0, i+1, 1, 1)
				}
				container.ShowAll()
			})
		}()

		return container
	}
	return &LibraryPage{"Artist: " + artist.Name, show}
}
