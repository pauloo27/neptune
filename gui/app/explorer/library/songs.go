package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	songsPage = &LibraryPage{"Songs", showSongs}
)

func showSongs() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	go func() {
		tracks, err := db.ListTracks(1)
		utils.HandleError(err, "Cannot list songs")

		// TODO: render
		glib.IdleAdd(func() {
			for i, track := range tracks {
				hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
				utils.HandleError(err, "Cannot create box")

				hbox.SetMarginStart(10)
				hbox.SetMarginEnd(10)

				songLabel, err := gtk.LabelNew(utils.Fmt("%s: %s", track.Album.Artist.Name, track.Title))
				utils.HandleError(err, "Cannot create label")

				hbox.PackStart(songLabel, false, false, 0)

				container.Attach(hbox, 0, i, 5, 1)
			}

			container.ShowAll()
		})
	}()

	return container
}
