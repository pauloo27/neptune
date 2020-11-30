package db

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	MBID string
	Name string
}

type Album struct {
	gorm.Model
	MBID     string
	Title    string
	ArtistID uint
	Artist   Artist
}

type Tag struct {
	gorm.Model
	Name string
}

type Track struct {
	gorm.Model
	MBID         string
	YoutubeId    string
	Album        Album
	Title        string
	Length       int
	YoutubeTitle string
	Tags         []Tag
}
