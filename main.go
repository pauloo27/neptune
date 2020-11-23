package main

import (
	"fmt"
	"os"

	"github.com/Pauloo27/my-tune/db"
	"github.com/Pauloo27/my-tune/gui/app"
	"github.com/Pauloo27/my-tune/player"
	"github.com/Pauloo27/my-tune/utils"
)

const version = "0.0.1"

func main() {
	fmt.Printf("Starting my-tune v%s\n", version)

	// load data folder
	home, err := os.UserHomeDir()
	utils.HandleError(err, "Cannot get user home")

	dataFolder := home + "/.cache/my-tune"

	_, err = os.Stat(dataFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0744)
		utils.HandleError(err, "Cannot create data folder")
	}

	// conect to db
	db.Connect(dataFolder)

	// add hook (not useful yet)
	player.RegisterHook(player.HOOK_PLAYER_INITIALIZED, func() {
		fmt.Println("The player was initialized!")
	})

	// start backend player
	player.Initialize()

	// start gui
	app.Start()
}
