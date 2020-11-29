package main

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/gui/app"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
)

const version = "0.0.1"

func main() {
	fmt.Printf("Starting neptune v%s\n", version)

	// load data folder
	home, err := os.UserHomeDir()
	utils.HandleError(err, "Cannot get user home")

	dataFolder := path.Join(home, ".cache", "neptune")

	_, err = os.Stat(dataFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0744)
		utils.HandleError(err, "Cannot create data folder")
	}

	// conect to db
	db.Connect(dataFolder)

	// add hook (not useful yet)
	player.RegisterHook(player.HOOK_PLAYER_INITIALIZED, func(params ...interface{}) {
		fmt.Println("The player was initialized")
	})

	// start backend player
	player.Initialize(dataFolder)

	// start gui
	app.Start()
}
