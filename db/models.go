package db

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name string
}

type Album struct {
	gorm.Model
	Name     string
	ArtistID uint
	Artist   Artist
}

type Tag struct {
	gorm.Model
	Name string
}

type Track struct {
	gorm.Model
	YoutubeId    string
	Album        Album
	Name         string
	Length       int
	YoutubeTitle string
	Tags         []Tag
}
