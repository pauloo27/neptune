package library

import (
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/gtk"
)

var (
	titleHeader      *gtk.Label
	goBackBtn        *gtk.Button
	contentContainer *gtk.Grid
	libraryContainer *gtk.Box

	homePage = &LibraryPage{"Home", showHome}
)

type ShowFunction func() *gtk.Grid

type LibraryPage struct {
	PageTitle string
	ShowPage  ShowFunction
}

func displayPage(page *LibraryPage) {
	titleHeader.SetText(page.PageTitle)

	if contentContainer != nil {
		contentContainer.Destroy()
	}

	contentContainer = page.ShowPage()

	libraryContainer.PackStart(contentContainer, true, true, 0)

	libraryContainer.ShowAll()
}

func showHome() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	container.SetColumnSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	i := 0
	addBtn := func(page *LibraryPage) {
		button, err := gtk.ButtonNewWithLabel(page.PageTitle)
		utils.HandleError(err, "Cannot create button")

		button.Connect("clicked", func() {
			displayPage(page)
		})

		container.Attach(button, i, 0, 1, 1)
		i++
	}

	addBtn(tracksPage)
	addBtn(artistsPage)
	/*
		addBtn("Artists", nil)
		addBtn("Albums", nil)
		addBtn("Tags", nil)
	*/
	return container
}

func CreateLibraryPage() *gtk.Box {
	var err error
	libraryContainer, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	utils.HandleError(err, "Cannot create box")

	titleHeader, err = gtk.LabelNew("")
	utils.HandleError(err, "Cannot create label")

	titleHeader.SetHAlign(gtk.ALIGN_CENTER)

	goBackBtn, err = gtk.ButtonNewWithLabel("Home")
	utils.HandleError(err, "Cannot create button")

	goBackBtn.Connect("clicked", func() {
		displayPage(homePage)
	})

	libraryContainer.PackStart(goBackBtn, false, false, 0)
	libraryContainer.PackStart(titleHeader, false, false, 0)

	displayPage(homePage)

	return libraryContainer
}
