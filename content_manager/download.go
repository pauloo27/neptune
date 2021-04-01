package content_manager

import (
	"github.com/Pauloo27/neptune/providers/youtube"
)

func Download(result *youtube.YoutubeEntry) {
	store(result, true)
}
