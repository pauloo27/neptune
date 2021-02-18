package explorer

import (
	"fmt"

	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createQueuePage() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	player.RegisterHook(
		player.HOOK_QUEUE_UPDATE_FINISHED,
		func(params ...interface{}) {
			fmt.Println("Updating queue...")
			glib.IdleAdd(func() {
				container.GetChildren().Foreach(func(item interface{}) {
					item.(*gtk.Widget).Destroy()
				})
				for i := 0; i < len(player.State.Queue); i++ {
					track := player.GetTrackAt(i)

					trackLabel, err := gtk.LabelNew(utils.Fmt("%s - %s", track.Album.Artist.Name, track.Title))
					utils.HandleError(err, "Cannot create label")

					container.PackStart(trackLabel, false, false, 1)
				}
				container.ShowAll()
			})
		},
	)

	return container
}
