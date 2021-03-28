package db

import (
	"errors"

	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"gorm.io/gorm"
)

const (
	PAGE_SIZE = 50
)

func CountPlayFor(track *Track) error {
	return Database.Model(track).Update("play_count", gorm.Expr("play_count + 1")).Error
}

func FindTrackByEntry(result *youtube.YoutubeEntry) (*Track, error) {
	var track Track
	err := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		First(&track, "youtube_id = ?", result.ID).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
	}
	return &track, err
}

func ListAllTracks() ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		Order("title collate nocase asc").
		Find(&tracks)

	return tracks, result.Error
}

func ListTracks(page int) ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		Order("title collate nocase asc").
		Limit(PAGE_SIZE).
		Offset(page * PAGE_SIZE).
		Find(&tracks)

	return tracks, result.Error
}

func ListArtists() ([]*Artist, error) {
	var artists []*Artist

	result := Database.Order("name collate nocase asc").Find(&artists)

	return artists, result.Error
}

func ListAlbumsBy(artist *Artist) ([]*Album, error) {
	var albums []*Album

	result := Database.Preload("Artist").Find(&albums, "artist_id = ?", artist.ID)

	return albums, result.Error
}

func ListAlbums() ([]*Album, error) {
	var albums []*Album

	result := Database.Preload("Artist").Order("title collate nocase asc").Find(&albums)

	return albums, result.Error
}

func ListTracksBy(artist *Artist) ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").Joins("Album").
		Order("tracks.title collate nocase asc").Find(&tracks, "Album__artist_id = ?", artist.ID)

	return tracks, result.Error
}

func ListTracksIn(album *Album) ([]*Track, error) {
	var tracks []*Track

	result := Database.
		Preload("Album.Artist").Preload("Tags.Tag").
		Order("title collate nocase asc").Find(&tracks, "album_id", album.ID)

	return tracks, result.Error
}

func ListTracksWith(tag *Tag) ([]*Track, error) {
	var tracks []*Track

	var trackTags []*TrackTag
	result := Database.Preload("Track.Album.Artist").Preload("Tag").
		Find(&trackTags, "tag_id = ?", tag.ID)

	if result.Error == nil {
		for _, trackTags := range trackTags {
			tracks = append(tracks, &trackTags.Track)
		}
	}

	return tracks, result.Error
}

func ListTags() ([]*Tag, error) {
	var tags []*Tag

	result := Database.Order("name collate nocase asc").Find(&tags)

	return tags, result.Error
}

func StoreTrack(videoInfo *youtube.VideoInfo, trackInfo *providers.TrackInfo, downloaded bool) (*Track, error) {
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
	if err != nil {
		return nil, err
	}

	// track
	track := Track{
		MBID:         trackInfo.MBID,
		Downloaded:   downloaded,
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

func LogStartup(version string) (previousVersion string, err error) {
	var startupLog NeptuneVersion
	res := Database.Last(&startupLog)
	if res.Error != nil && !errors.Is(gorm.ErrRecordNotFound, res.Error) {
		return "", res.Error
	}
	previousVersion = startupLog.Version

	if version != previousVersion {
		newStartupLog := NeptuneVersion{Version: version}
		err = Database.Create(&newStartupLog).Error
	}
	return
}
