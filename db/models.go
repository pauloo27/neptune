package db

import "gorm.io/gorm"

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

type Track struct {
	gorm.Model
	MBID         string `gorm:"unique"`
	YoutubeId    string
	Album        Album
	Title        string
	Length       int
	YoutubeTitle string
	Tags         []Tag
}
