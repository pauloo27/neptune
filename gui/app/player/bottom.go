package player

import (
	"github.com/Pauloo27/my-tune/player"
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/gtk"
)

func createProgressBar() *gtk.Scale {
	progressBar, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	utils.HandleError(err, "Cannot create scale")

	progressBar.SetDrawValue(false)
	progressBar.SetHExpand(true)

	return progressBar
}

func createVolumeController() *gtk.Box {
	volumeContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	volumeIcon, err := gtk.ImageNewFromIconName("audio-volume-medium", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	volumeController, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0.0, 100.0, 1.0)
	utils.HandleError(err, "Cannot create box")

	volumeController.SetValue(50.0)
	volumeController.SetDrawValue(false)
	volumeController.Connect("value-changed", func() {
		player.SetVolume(volumeController.GetValue())
	})

	volumeContainer.PackStart(volumeIcon, false, false, 0)
	volumeContainer.PackEnd(volumeController, true, true, 0)

	return volumeContainer
}

func createDurationLabel() *gtk.Label {
	durationLabel, err := gtk.LabelNew("3:21")
	utils.HandleError(err, "Cannot create label")

	durationLabel.SetHAlign(gtk.ALIGN_END)

	return durationLabel
}

func createPositionLabel() *gtk.Label {
	positonLabel, err := gtk.LabelNew("1:23")
	utils.HandleError(err, "Cannot create label")

	positonLabel.SetHAlign(gtk.ALIGN_START)

	return positonLabel
}

func createSongLabel() *gtk.Label {
	songLabel, err := gtk.LabelNew("Testing - I have no ideia")
	utils.HandleError(err, "Cannot create label")

	return songLabel
}

func createButtonsContainer() *gtk.Box {
	buttonsContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot creat ebox")

	pauseButton, err := gtk.ButtonNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	buttonsContainer.SetHAlign(gtk.ALIGN_CENTER)

	buttonsContainer.PackStart(pauseButton, false, true, 0)

	return buttonsContainer
}

func createPlayerBottom() *gtk.Grid {
	bottomContainer, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	// row 0
	bottomContainer.Attach(createProgressBar(), 0, 0, 10, 1)
	// row 1
	bottomContainer.Attach(createPositionLabel(), 0, 1, 1, 1)
	bottomContainer.Attach(createSongLabel(), 1, 1, 8, 1)
	bottomContainer.Attach(createDurationLabel(), 9, 1, 1, 1)
	// row 2
	bottomContainer.Attach(createVolumeController(), 0, 2, 3, 1)
	bottomContainer.Attach(createButtonsContainer(), 3, 2, 4, 1)

	return bottomContainer
}
