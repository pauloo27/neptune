package player

import (
	"time"

	"github.com/Pauloo27/my-tune/player"
	"github.com/Pauloo27/my-tune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var positionLabel *gtk.Label
var progressBar *gtk.Scale
var currentPosition float64

func createProgressBar() *gtk.Scale {
	var err error

	progressBar, err = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0.0, 1.0, 0.01)
	utils.HandleError(err, "Cannot create scale")

	progressBar.SetDrawValue(false)
	progressBar.SetHExpand(true)
	progressBar.Connect("value-changed", func() {
		if currentPosition == progressBar.GetValue() {
			return
		}
		player.SetPosition(progressBar.GetValue() * player.State.Duration)
	})

	return progressBar
}

func createVolumeController() *gtk.Box {
	volumeContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	volumeIcon, err := gtk.ImageNewFromIconName("audio-volume-medium", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	volumeController, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0.0, 100.0, 1.0)
	utils.HandleError(err, "Cannot create box")

	volumeController.SetDrawValue(false)
	volumeController.SetValue(player.State.Volume)

	player.RegisterHook(player.HOOK_VOLUME_CHANGED, func(err error, params ...interface{}) {
		volume := params[0].(float64)
		if volume != volumeController.GetValue() {
			volumeController.SetValue(volume)
		}
	})

	volumeController.Connect("value-changed", func() {
		player.SetVolume(volumeController.GetValue())
	})

	volumeContainer.PackStart(volumeIcon, false, false, 0)
	volumeContainer.PackEnd(volumeController, true, true, 0)

	return volumeContainer
}

func createDurationLabel() *gtk.Label {
	durationLabel, err := gtk.LabelNew("--:--")
	utils.HandleError(err, "Cannot create label")

	durationLabel.SetHAlign(gtk.ALIGN_END)
	player.RegisterHook(player.HOOK_FILE_LOADED, func(err error, params ...interface{}) {
		duration := params[0].(float64)

		durationLabel.SetText(utils.FormatDuration(time.Duration(duration) * time.Second))
	})

	return durationLabel
}

func updatePosition(pos float64) {
	positionLabel.SetText(utils.FormatDuration(time.Duration(pos) * time.Second))
	currentPosition = pos / player.State.Duration
	progressBar.SetValue(currentPosition)
}

func createPositionLabel() *gtk.Label {
	var err error
	positionLabel, err = gtk.LabelNew("--:--")
	utils.HandleError(err, "Cannot create label")

	positionLabel.SetHAlign(gtk.ALIGN_START)

	go func() {
		for {
			pos, err := player.GetPosition()
			if err == nil {
				glib.IdleAdd(func() {
					updatePosition(pos)
				})
			}
			time.Sleep(500 * time.Microsecond)
		}
	}()

	return positionLabel
}

func createSongLabel() *gtk.Label {
	songLabel, err := gtk.LabelNew("Testing - I have no ideia")
	utils.HandleError(err, "Cannot create label")

	return songLabel
}

func createButtonsContainer() *gtk.Box {
	buttonsContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	pausedIcon, err := gtk.ImageNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	playingIcon, err := gtk.ImageNewFromIconName("media-playback-pause", gtk.ICON_SIZE_BUTTON)

	pauseButton, err := gtk.ButtonNew()
	utils.HandleError(err, "Cannot create button")

	pauseButton.SetImage(playingIcon)

	player.RegisterHook(player.HOOK_PLAYBACK_PAUSED, func(err error, params ...interface{}) {
		pauseButton.SetImage(pausedIcon)
	})
	player.RegisterHook(player.HOOK_PLAYBACK_RESUMED, func(err error, params ...interface{}) {
		pauseButton.SetImage(playingIcon)
	})

	pauseButton.Connect("clicked", func() {
		player.PlayPause()
	})

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
