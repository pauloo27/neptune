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

func displayTrack(track *db.Track, showArtist bool) *gtk.Box {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	hbox.SetMarginStart(5)
	hbox.SetMarginEnd(5)

	playButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	playButton.Connect("clicked", func() {
		go player.PlayTrack(track)
	})

	addToQueueButton, err := gtk.ButtonNewFromIconName("list-add", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	addToQueueButton.Connect("clicked", func() {
		go player.AddTrackToQueue(track)
	})

	var fullTitle string
	if showArtist {
		fullTitle = utils.Fmt("%s from %s", utils.EnforceSize(track.Title, 50), utils.EnforceSize(track.Album.Artist.Name, 20))
	} else {
		fullTitle = track.Title
	}

	songLabel, err := gtk.LabelNew(fullTitle)
	utils.HandleError(err, "Cannot create label")

	hbox.PackStart(playButton, false, false, 1)
	hbox.PackStart(addToQueueButton, false, false, 1)
	hbox.PackStart(songLabel, false, false, 1)
	return hbox
}

func appendTracks(container *gtk.Grid, tracks []*db.Track, offset int) int {
	for i, track := range tracks {
		container.Attach(displayTrack(track, true), 0, i+offset, 10, 1)
	}

	container.ShowAll()
	return len(tracks)
}

func showTracks() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	offset := 0
	page := 0

	tracksContainer, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	tracksContainer.SetRowSpacing(1)
	container.SetRowSpacing(1)

	container.SetColumnHomogeneous(true)

	loadPage := func(page int) {
		tracks, err := db.ListTracks(page)
		utils.HandleError(err, "Cannot list songs")

		glib.IdleAdd(func() {
			offset += appendTracks(tracksContainer, tracks, offset)
		})
	}

	loadMoreButton, err := gtk.ButtonNewWithLabel("Load more")
	utils.HandleError(err, "Cannot create button")

	loadMoreButton.Connect("clicked", func() {
		page++
		go loadPage(page)
	})

	container.Attach(tracksContainer, 0, 0, 1, 1)
	container.Attach(loadMoreButton, 0, 1, 1, 1)

	go loadPage(page)

	return container
}
