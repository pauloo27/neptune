package player

import (
	"os"

	"github.com/Pauloo27/neptune/utils"
)

var (
	userHomePaths = []string{"/.config/mpv/scripts/mpris.so"}
	systemPaths   = []string{"/etc/mpv/scripts/mpris.so"}
	loaded        = false
)

func loadMPRIS() {
	if loaded {
		return
	}
	loadScript := func(path string) bool {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}

		err := MpvInstance.CommandString("load-script " + path)
		if err != nil {
			utils.HandleError(err, "Cannot load mpris script at "+path)
		}
		return true
	}

	defer func() { loaded = true }()

	home := utils.GetUserHome()
	for _, path := range userHomePaths {
		if loadScript(home + path) {
			return
		}
	}

	for _, path := range systemPaths {
		if loadScript(path) {
			return
		}
	}
}
