package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createTagPage(tag *db.Tag) *LibraryPage {
	show := func() *gtk.Grid {
		container, err := gtk.GridNew()
		utils.HandleError(err, "Cannot create box")

		container.SetRowSpacing(5)
		container.SetColumnHomogeneous(true)
		container.SetMarginStart(5)
		container.SetMarginEnd(5)

		go func() {
			tracks, err := db.ListTracksWith(tag)
			utils.HandleError(err, "Cannot list tracks with tag "+tag.Name)

			glib.IdleAdd(func() {
				container.Attach(createPlayAll("Play all", tracks), 0, 0, 1, 1)
				for i, track := range tracks {
					container.Attach(displayTrack(track, true), 0, i+1, 1, 1)
				}
				container.ShowAll()
			})
		}()

		return container
	}

	return &LibraryPage{"Tag: " + tag.Name, show}
}
