package providers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Pauloo27/neptune/providers/youtube"
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

func FetchTrackInfo(info *youtube.VideoInfo) (*TrackInfo, error) {
	fmt.Printf("Fetching track info for %s by %s\n", info.Track, info.Artist)

	/*
		here are some old stuff that I used... maybe it will be necessary here
		as well
	*/
	// var parenthesisRegex = regexp.MustCompile(`\s?\(.+\)`)
	// fix track with '(stuff)'
	//trackName := parenthesisRegex.ReplaceAllString(videoInfo.Track, "")
	// fix for "artist" list (splitted by ',')
	//artist := strings.Split(videoInfo.Artist, ",")[0]

	// escape params
	escapedArtist := url.QueryEscape(info.Artist)
	escapedTrack := url.QueryEscape(info.Track)

	reqPath := utils.Fmt(
		"%s/?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json",
		ENDPOINT, API_KEY, escapedArtist, escapedTrack,
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
		artistName = info.Artist
	}

	artistMBID, err := jsonparser.GetString(buffer, "track", "artist", "mbid")
	if err != nil {
		artistMBID = "*" + artistName
	}

	// album info
	albumTitle, err := jsonparser.GetString(buffer, "track", "album", "title")
	if err != nil {
		albumTitle = utils.Fmt("Unknown %s's album", artistName)
	}

	albumMBID, err := jsonparser.GetString(buffer, "track", "album", "mbid")
	if err != nil {
		albumMBID = artistMBID + "|" + albumTitle
	}

	albumImageURL, err := jsonparser.GetString(buffer, "track", "album", "image", "[3]", "#text")
	if err != nil {
		albumImageURL = info.GetThumbnail()
	}

	// track info
	trackTitle, err := jsonparser.GetString(buffer, "track", "name")
	if err != nil {
		trackTitle = info.Track
	}

	trackMBID, err := jsonparser.GetString(buffer, "track", "mbid")
	if err != nil {
		trackMBID = albumMBID + "|" + trackTitle
	}

	var trackTags []string
	tagsArr, _, _, err := jsonparser.Get(buffer, "track", "toptags", "tag")

	_, err = jsonparser.ArrayEach(tagsArr, func(data []byte, t jsonparser.ValueType, i int, err error) {
		tagName, err := jsonparser.GetString(data, "name")
		trackTags = append(trackTags, tagName)
	})

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
