package player

import "github.com/Pauloo27/my-tune/youtube"

type PlayerState struct {
	Paused   bool
	Playing  *youtube.YoutubeEntry
	Volume   float64
	Duration float64
}
