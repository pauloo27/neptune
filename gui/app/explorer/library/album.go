package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createAlbumPage(album *db.Album) *LibraryPage {
	show := func() *gtk.Grid {
		container, err := gtk.GridNew()
		utils.HandleError(err, "Cannot create grid")

		go func() {
			tracks, err := db.ListTracksIn(album)
			utils.HandleError(err, "Cannot list tracks in album "+album.MBID)

			glib.IdleAdd(func() {
				container.Attach(createPlayAll("Play album", tracks), 0, 0, 1, 1)

				albumArt, err := gtk.ImageNew()
				utils.HandleError(err, "Cannot create image")

				imagePath := album.GetAlbumArtPath()
				imagePix, err := gdk.PixbufNewFromFileAtScale(imagePath, 150, 150, true)
				utils.HandleError(err, "Cannot load image from file")
				albumArt.SetFromPixbuf(imagePix)
				albumArt.SetHAlign(gtk.ALIGN_CENTER)
				albumArt.SetMarginBottom(5)
				container.SetHAlign(gtk.ALIGN_CENTER)

				container.Attach(albumArt, 0, 1, 10, 1)

				for i, track := range tracks {
					container.Attach(displayTrack(track, false), 0, i+2, 10, 1)
				}
				container.ShowAll()
			})
		}()

		return container
	}

	return &LibraryPage{"Album: " + album.Title, show}
}
