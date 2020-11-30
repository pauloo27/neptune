package db

import (
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/youtube"
	"gorm.io/gorm/clause"
)

func StoreTrack(videoInfo *youtube.VideoInfo, trackInfo *providers.TrackInfo) error {
	// artist
	artist := Artist{
		MBID: trackInfo.Artist.MBID,
		Name: trackInfo.Artist.Name,
	}
	Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&artist)
	// album
	album := Album{
		MBID:   trackInfo.Album.MBID,
		Title:  trackInfo.Album.Title,
		Artist: artist,
	}
	Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&album)
	// track
	track := Track{
		MBID:         trackInfo.MBID,
		YoutubeID:    videoInfo.ID,
		Album:        album,
		Title:        trackInfo.Title,
		Length:       int(videoInfo.Duration),
		YoutubeTitle: videoInfo.Title,
	}
	Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&track)
	// tags
	for _, tagName := range trackInfo.Tags {
		tag := Tag{
			Name: tagName,
		}
		Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&tag)
		trackTag := TrackTag{
			Track: track,
			Tag:   tag,
		}
		Database.Create(&trackTag)
	}
	return nil
}
