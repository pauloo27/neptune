package content_manager

import (
	"github.com/Pauloo27/neptune/providers/youtube"
)

func Stream(result *youtube.YoutubeEntry) {
	store(result, false)
}
