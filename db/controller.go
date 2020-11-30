package db

import (
	"path"

	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func TrackFrom(result *youtube.YoutubeEntry) (*Track, error) {
	var track Track
	err := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		First(&track, "youtube_id = ?", result.ID).Error
	return &track, err
}

func StoreTrack(videoInfo *youtube.VideoInfo, trackInfo *providers.TrackInfo) (*Track, error) {
	var err error
	// artist
	artist := Artist{
		MBID: trackInfo.Artist.MBID,
		Name: trackInfo.Artist.Name,
	}
	err = Database.Where(Artist{MBID: trackInfo.Artist.MBID}).
		FirstOrCreate(&artist).Error
	if err != nil {
		return nil, err
	}
	// album
	album := Album{
		MBID:   trackInfo.Album.MBID,
		Title:  trackInfo.Album.Title,
		Artist: artist,
	}
	err = Database.Where(Album{MBID: trackInfo.Album.MBID}).
		FirstOrCreate(&album).Error
	// download album art
	err = utils.DownloadFile(trackInfo.Album.ImageURL,
		path.Join(DataFolder, "albums", trackInfo.Album.MBID+".png"),
	)
	utils.HandleError(err, "Cannot download album image")
	if err != nil {
		return nil, err
	}
	// track
	track := Track{
		MBID:         trackInfo.MBID,
		YoutubeID:    videoInfo.ID,
		Album:        album,
		Title:        trackInfo.Title,
		Length:       int(videoInfo.Duration),
		YoutubeTitle: videoInfo.Title,
	}
	err = Database.Where(Track{MBID: trackInfo.MBID}).
		FirstOrCreate(&track).Error
	if err != nil {
		return nil, err
	}
	// tags
	for _, tagName := range trackInfo.Tags {
		tag := Tag{
			Name: tagName,
		}
		err = Database.Where(Tag{Name: tagName}).
			FirstOrCreate(&tag).Error
		if err != nil {
			return nil, err
		}
		trackTag := TrackTag{
			Track: track,
			Tag:   tag,
		}
		err = Database.Create(&trackTag).Error
		if err != nil {
			return nil, err
		}
	}
	return &track, nil
}
