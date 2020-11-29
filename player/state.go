package player

import "github.com/Pauloo27/neptune/youtube"

type PlayerState struct {
	Paused   bool
	Playing  *youtube.YoutubeEntry
	Volume   float64
	Duration float64
}
