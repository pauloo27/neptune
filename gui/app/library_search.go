package app

import (
	"github.com/Pauloo27/my-tune/utils"
	"github.com/Pauloo27/my-tune/youtube"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var searchResultsContainer *gtk.Box
var searchStatusLabel *gtk.Label

func doSearch(searchTerm string) {
	if searchTerm == "" {
		return
	}

	searchStatusLabel.SetText("Searching for " + searchTerm)

	go func() {
		results, err := youtube.SearchFor(searchTerm, 10)
		if err != nil {
			glib.IdleAdd(func() {
				searchStatusLabel.SetText("Something went wrong")
			})
		} else {
			glib.IdleAdd(func() {
				searchStatusLabel.SetText("Results:")
				for _, result := range results {
					label, err := gtk.LabelNew(utils.Fmt("%s - %s | %s",
						utils.EnforceSize(result.Title, 40),
						utils.EnforceSize(result.Uploader, 30),
						result.Duration,
					))
					utils.HandleError(err, "Cannot create label")

					searchResultsContainer.PackStart(label, false, false, 1)
				}
				searchResultsContainer.ShowAll()
			})
		}
	}()
}

func CreateSearchHeader() *gtk.HeaderBar {
	searchBarContainer, err := gtk.HeaderBarNew()
	utils.HandleError(err, "Cannot create header bar")

	searchBarContainer.SetShowCloseButton(false)

	searchInput, err := gtk.EntryNew()
	utils.HandleError(err, "Cannot create entry")

	searchInput.SetPlaceholderText("Search YouTube")
	searchInput.SetHExpand(true)

	handleSearch := func() {
		text, err := searchInput.GetText()
		utils.HandleError(err, "Cannot get entry text")

		doSearch(text)
	}

	searchButton, err := gtk.ButtonNewFromIconName("search", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	searchInput.Connect("activate", handleSearch)
	searchButton.Connect("clicked", handleSearch)

	searchBarContainer.PackStart(searchInput)
	searchBarContainer.PackStart(searchButton)

	return searchBarContainer
}

func CreateSearchPage() *gtk.Box {
	var err error
	searchResultsContainer, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	searchStatusLabel, err = gtk.LabelNew("Nothing yet")
	utils.HandleError(err, "Cannot create label")

	searchResultsContainer.PackStart(searchStatusLabel, false, false, 1)

	return searchResultsContainer
}
