package db

import (
	"path"

	"github.com/Pauloo27/neptune/utils"
	"gorm.io/gorm"
)

type NeptuneVersion struct {
	gorm.Model
	Version string
}

type Artist struct {
	gorm.Model
	MBID string `gorm:"unique"`
	Name string
}

type Album struct {
	gorm.Model
	MBID     string `gorm:"unique"`
	Title    string
	ArtistID uint
	Artist   Artist
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type TrackTag struct {
	gorm.Model
	TagID   uint
	Tag     Tag
	TrackID uint
	Track   Track
}

type Track struct {
	gorm.Model
	MBID         string
	Downloaded   bool
	YoutubeID    string `gorm:"unique"`
	AlbumID      uint
	Album        Album
	Title        string
	Length       int
	PlayCount    int
	YoutubeTitle string
	Tags         []TrackTag
}

func (t *Track) GetLocalPath() string {
	return path.Join(DataFolder, "albums", utils.Fmt("%d", t.Album.ID), t.YoutubeID+".m4a")
}

func (t *Track) GetYouTubeURL() string {
	return utils.Fmt("https://www.youtube.com/watch?v=%s", t.YoutubeID)
}

func (t *Track) GetPath() string {
	if t.Downloaded {
		return t.GetLocalPath()
	}
	return t.GetYouTubeURL()
}

func (a *Album) GetAlbumArtPath() string {
	return path.Join(a.GetAlbumPath(), ".folder.png")
}

func (a *Album) GetAlbumPath() string {
	return path.Join(DataFolder, "albums", utils.Fmt("%d", a.ID))
}
