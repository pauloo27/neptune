package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var tagsPage = &LibraryPage{"Tags", showTags}

func displayTag(tag *db.Tag) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(tag.Name)
	utils.HandleError(err, "Cannot create button")

	btn.Connect("clicked", func() {
		displayPage(createTagPage(tag))
	})

	return btn
}

func showTags() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	container.SetRowSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	go func() {
		tags, err := db.ListTags(1)
		utils.HandleError(err, "Cannot list tags")

		glib.IdleAdd(func() {
			for i, tag := range tags {
				container.Attach(displayTag(tag), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

	return container
}
