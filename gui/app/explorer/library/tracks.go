package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	tracksPage = &LibraryPage{"Tracks", showTracks}
)

func displayTrack(track *db.Track) *gtk.Box {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	hbox.SetMarginStart(5)
	hbox.SetMarginEnd(5)

	playButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create label")

	playButton.Connect("clicked", func() {
		go player.PlayTrack(track)
	})

	songLabel, err := gtk.LabelNew(utils.Fmt("%s: %s", track.Album.Artist.Name, track.Title))
	utils.HandleError(err, "Cannot create label")

	hbox.PackStart(playButton, false, false, 1)
	hbox.PackStart(songLabel, false, false, 1)
	return hbox
}

func showTracks() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	container.SetRowSpacing(1)

	go func() {
		tracks, err := db.ListTracks(1)
		utils.HandleError(err, "Cannot list songs")

		glib.IdleAdd(func() {
			for i, track := range tracks {
				container.Attach(displayTrack(track), 0, i, 10, 1)
			}

			container.ShowAll()
		})
	}()

	return container
}
