package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/gui/app"
	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/trayicon"
	"github.com/Pauloo27/neptune/utils"
	"github.com/Pauloo27/neptune/version"
)

func main() {
	fmt.Printf("Starting neptune %s\n", version.VERSION)

	// load data folder
	home, err := os.UserHomeDir()
	utils.HandleError(err, "Cannot get user home")

	dataFolder := path.Join(home, ".cache", "neptune")

	_, err = os.Stat(dataFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0744)
		utils.HandleError(err, "Cannot create data folder")
	}

	albumsCacheFolder := path.Join(dataFolder, "albums")
	_, err = os.Stat(albumsCacheFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(albumsCacheFolder, 0744)
		utils.HandleError(err, "Cannot create albums cache folder")
	}

	// conect to db
	db.Connect(dataFolder)

	// save the current version (used in migrations)
	prevVersion, err := db.LogStartup(version.VERSION)
	utils.HandleError(err, "Cannot log current version to db")

	version.MigrateFrom(prevVersion)

	// add hook (not useful yet)
	hook.RegisterHook(hook.HOOK_PLAYER_INITIALIZED, func(params ...interface{}) {
		fmt.Println("The player was initialized")
	})

	// start backend player
	player.Initialize(dataFolder)

	hook.RegisterHook(hook.HOOK_GUI_STARTED, func(params ...interface{}) {
		go func() {
			// to avoid random crash
			// w t f
			time.Sleep(1 * time.Second)
			trayicon.LoadTrayIcon()
		}()
	})

	// start gui
	app.Start(func() {
		player.Exit()
	})
}
