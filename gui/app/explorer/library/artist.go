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

		go func() {
			albums, err := db.ListAlbumsBy(artist, 1)
			utils.HandleError(err, "Cannot list albums by artist "+artist.MBID)

			tracks, err := db.ListTracksBy(artist, 1)
			utils.HandleError(err, "Cannot list tracks by artist "+artist.MBID)

			glib.IdleAdd(func() {
				albumsLabel, err := gtk.LabelNew("Albums:")
				utils.HandleError(err, "Cannot create label")
				container.Attach(albumsLabel, 0, 0, 1, 1)

				i := 1
				for _, album := range albums {
					container.Attach(displayAlbum(album, false), 0, i, 1, 1)
					i++
				}

				tracksLabel, err := gtk.LabelNew("Tracks:")
				utils.HandleError(err, "Cannot create label")
				container.Attach(tracksLabel, 0, i, 1, 1)

				i++
				for _, track := range tracks {
					container.Attach(displayTrack(track), 0, i, 1, 1)
					i++
				}
				container.ShowAll()
			})
		}()

		return container
	}
	return &LibraryPage{"Artist: " + artist.Name, show}
}
