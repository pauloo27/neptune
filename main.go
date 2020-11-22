package main

import (
	"fmt"
	"os"

	"github.com/Pauloo27/my-tune/db"
	"github.com/Pauloo27/my-tune/gui/app"
	"github.com/Pauloo27/my-tune/utils"
)

const version = "0.0.1"

func main() {
	fmt.Printf("Starting my-tune v%s\n", version)

	home, err := os.UserHomeDir()
	utils.HandleError(err, "Cannot get user home")

	dataFolder := home + "/.cache/my-tune"

	_, err = os.Stat(dataFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0744)
		utils.HandleError(err, "Cannot create data folder")
	}

	db.Connect(dataFolder)

	app.Start()
}
