package explorer

import (
	"runtime"

	"github.com/Pauloo27/neptune/content_manager"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var searchResultsContainer *gtk.Box
var searchStatusLabel *gtk.Label
var searching = false

func doSearch(searchTerm string) {
	if searchTerm == "" || searching {
		return
	}
	searching = true
	explorerContainer.SetCurrentPage(1)

	children := searchResultsContainer.GetChildren()

	id := 0
	children.Foreach(func(item interface{}) {
		wid := item.(*gtk.Widget)
		// ignore the label
		if id != 0 {
			wid.Destroy()
		}
		id++
	})
	searchStatusLabel.SetText("Searching for " + searchTerm)

	go func() {
		results, err := youtube.SearchFor(searchTerm, 10)
		if err != nil {
			glib.IdleAdd(func() {
				searchStatusLabel.SetText("Something went wrong")
				searching = false
			})
		} else {
			glib.IdleAdd(func() {
				searchStatusLabel.SetText("Results:")
				for _, result := range results {
					resultContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
					utils.HandleError(err, "Cannot create label")

					playButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
					utils.HandleError(err, "Cannot create label")

					// var result, at the end of the for, will be the last array element
					currentResult := result
					playButton.Connect("clicked", func() {
						if runtime.GOOS == "windows" {
							content_manager.Download(currentResult)
						} else {
							content_manager.Download(currentResult)
							//content_manager.Stream(currentResult)
						}
					})

					label, err := gtk.LabelNew(utils.Fmt("%s - %s | %s",
						utils.EnforceSize(result.Title, 60),
						utils.EnforceSize(result.Uploader, 20),
						result.DisplayDuration,
					))
					utils.HandleError(err, "Cannot create label")

					resultContainer.PackStart(playButton, false, false, 1)
					resultContainer.PackStart(label, false, false, 1)

					searchResultsContainer.PackStart(resultContainer, false, false, 1)
				}
				searchResultsContainer.ShowAll()
				searching = false
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

func createSearchPage() *gtk.Box {
	var err error
	searchResultsContainer, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	searchStatusLabel, err = gtk.LabelNew("Nothing yet")
	utils.HandleError(err, "Cannot create label")

	searchResultsContainer.PackStart(searchStatusLabel, false, false, 1)

	return searchResultsContainer
}
