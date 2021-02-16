package player

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers/youtube"
)

type PlayerState struct {
	Paused   bool
	Fetching *youtube.YoutubeEntry
	Queue    []*db.Track
	Volume   float64
	Duration float64
}
