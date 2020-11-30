package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Pauloo27/neptune/utils"
	"github.com/buger/jsonparser"
)

type ArtistInfo struct {
	Name, MBID string
}

type AlbumInfo struct {
	Title, MBID, ImageURL string
}

type TrackInfo struct {
	Title, MBID string
	Tags        []string
	Artist      *ArtistInfo
	Album       *AlbumInfo
}

const (
	API_KEY  = "12dec50313f885d407cf8132697b8712"
	ENDPOINT = "https://ws.audioscrobbler.com/2.0"
)

func FetchTrackInfo(artist, track string) (*TrackInfo, error) {
	reqPath := utils.Fmt(
		"%s/?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json",
		ENDPOINT, API_KEY, artist, track,
	)

	res, err := http.Get(reqPath)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)

	// artist info
	artistName, err := jsonparser.GetString(buffer, "track", "artist", "name")
	if err != nil {
		return nil, err
	}

	artistMBID, err := jsonparser.GetString(buffer, "track", "artist", "mbid")
	if err != nil {
		return nil, err
	}

	// album info
	albumTitle, err := jsonparser.GetString(buffer, "track", "album", "title")
	if err != nil {
		return nil, err
	}

	albumMBID, err := jsonparser.GetString(buffer, "track", "album", "mbid")
	if err != nil {
		return nil, err
	}

	albumImageURL, err := jsonparser.GetString(buffer, "track", "album", "image", "[3]", "#text")
	if err != nil {
		return nil, err
	}

	// track info
	trackTitle, err := jsonparser.GetString(buffer, "track", "name")
	if err != nil {
		return nil, err
	}

	trackMBID, err := jsonparser.GetString(buffer, "track", "mbid")
	if err != nil {
		return nil, err
	}

	var trackTags []string
	tagsArr, _, _, err := jsonparser.Get(buffer, "track", "toptags", "tag")

	_, err = jsonparser.ArrayEach(tagsArr, func(data []byte, t jsonparser.ValueType, i int, err error) {
		tagName, err := jsonparser.GetString(data, "name")
		trackTags = append(trackTags, tagName)
	})
	if err != nil {
		return nil, err
	}

	return &TrackInfo{
		Title: trackTitle,
		MBID:  trackMBID,
		Tags:  trackTags,
		Album: &AlbumInfo{
			Title:    albumTitle,
			MBID:     albumMBID,
			ImageURL: albumImageURL,
		},
		Artist: &ArtistInfo{
			Name: artistName,
			MBID: artistMBID,
		},
	}, nil
}
