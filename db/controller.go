package db

import (
	"errors"
	"path"

	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
	"gorm.io/gorm"
)

func PlayTrack(track *Track) error {
	return Database.Model(track).Update("play_count", gorm.Expr("play_count + 1")).Error
}

func PlayEntry(result *youtube.YoutubeEntry) (*Track, error) {
	track, err := TrackFrom(result)

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return track, err
	}

	err = PlayTrack(track)

	return track, err
}

func TrackFrom(result *youtube.YoutubeEntry) (*Track, error) {
	var track Track
	err := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		First(&track, "youtube_id = ?", result.ID).Error

	return &track, err
}

func ListTracks(page int) ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		Order("play_count desc").Find(&tracks)

	return tracks, result.Error
}

func ListArtists(page int) ([]*Artist, error) {
	var artists []*Artist

	result := Database.Find(&artists)

	return artists, result.Error
}

func ListAlbumsBy(artist *Artist, page int) ([]*Album, error) {
	var albums []*Album

	result := Database.Find(&albums, "artist_id = ?", artist.ID)

	return albums, result.Error
}

func ListTracksBy(artist *Artist, page int) ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").Joins("Album").
		Order("play_count desc").Find(&tracks, "Album__artist_id = ?", artist.ID)

	return tracks, result.Error
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
		path.Join(DataFolder, "albums", trackInfo.Album.MBID, ".folder.png"),
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
		PlayCount:    1,
		Length:       int(videoInfo.Duration),
		YoutubeTitle: videoInfo.Title,
	}
	err = Database.Where(Track{YoutubeID: videoInfo.ID}).
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
