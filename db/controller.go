package db

import (
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/youtube"
)

func StoreTrack(videoInfo *youtube.VideoInfo, trackInfo *providers.TrackInfo) error {
	// artist
	artist := Artist{
		MBID: trackInfo.Artist.MBID,
		Name: trackInfo.Artist.Name,
	}
	Database.FirstOrCreate(&artist)
	// album
	album := Album{
		MBID:   trackInfo.Album.MBID,
		Title:  trackInfo.Album.Title,
		Artist: artist,
	}
	Database.FirstOrCreate(&album)
	// track
	track := Track{
		MBID:         trackInfo.MBID,
		YoutubeID:    videoInfo.ID,
		Album:        album,
		Title:        trackInfo.Title,
		Length:       int(videoInfo.Duration),
		YoutubeTitle: videoInfo.Title,
	}
	Database.FirstOrCreate(&track)
	// tags
	for _, tagName := range trackInfo.Tags {
		tag := Tag{
			Name: tagName,
		}
		Database.FirstOrCreate(&tag)
		trackTag := TrackTag{
			Track: track,
			Tag:   tag,
		}
		Database.Create(&trackTag)
	}
	return nil
}
