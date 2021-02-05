package player

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

var parenthesisRegex = regexp.MustCompile(`\s?\(.+\)`)

func PlayTrack(track *db.Track) {
	RemoveCurrentFromPlaylist()

	filePath := track.GetPath()
	err := db.PlayTrack(track)
	utils.HandleError(err, "Cannot play track")

	State.Track = track
	LoadFile(filePath)
}

func PlayResult(result *youtube.YoutubeEntry) {
	RemoveCurrentFromPlaylist()
	State.Fetching = result
	callHooks(HOOK_RESULT_FETCH_STARTED, nil)
	track, err := db.PlayEntry(result)
	utils.HandleError(err, "Cannot find track")
	if track == nil {
		callHooks(HOOK_RESULT_DOWNLOAD_STARTED, nil)
		go func() {
			fmt.Println("Downloading file...")
			tmpFile := path.Join(DataFolder, "wip_downloads", result.ID+".m4a")
			videoInfo, err := youtube.DownloadResult(result, tmpFile)
			utils.HandleError(err, "Cannot download file")
			var trackInfo *providers.TrackInfo
			if videoInfo.Artist == "" || videoInfo.Track == "" {
				artistInfo := providers.ArtistInfo{
					Name: videoInfo.Uploader,
					MBID: "!YT:" + videoInfo.UploaderID,
				}
				albumInfo := providers.AlbumInfo{
					Title: "YouTube video by " + videoInfo.Uploader,
					MBID:  "!YT:" + videoInfo.ID,
					ImageURL: utils.Fmt(
						"https://i1.ytimg.com/vi/%s/hqdefault.jpg", videoInfo.ID,
					),
				}
				trackInfo = &providers.TrackInfo{
					Artist: &artistInfo,
					Album:  &albumInfo,
					Title:  videoInfo.Title,
					MBID:   "!YT:" + videoInfo.ID,
				}
			} else {
				var err error
				// fix track with '(stuff)'
				trackName := parenthesisRegex.ReplaceAllString(videoInfo.Track, "")
				// fix for "artist" list (splitted by ',')
				artist := strings.Split(videoInfo.Artist, ",")[0]
				trackInfo, err = providers.FetchTrackInfo(artist, trackName)
				utils.HandleError(err, "Cannot fetch track info")
			}
			albumPath := path.Join(DataFolder, "albums", trackInfo.Album.MBID)
			if _, err := os.Stat(albumPath); os.IsNotExist(err) {
				err = os.MkdirAll(albumPath, 0744)
				utils.HandleError(err, "Cannot create album folder at"+albumPath)
			}
			filePath := path.Join(albumPath, result.ID+".m4a")
			err = os.Rename(tmpFile, filePath)
			utils.HandleError(err, "Cannot move tmp download file to cache file")
			track, err = db.StoreTrack(videoInfo, trackInfo)
			utils.HandleError(err, "Cannot store track")

			State.Track = track
			LoadFile(filePath)
		}()
	} else {
		State.Track = track
		LoadFile(track.GetPath())
	}
	State.Fetching = nil
}
