package content_manager

import (
	"os"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func Stream(result *youtube.YoutubeEntry) {
	player.ClearQueue()

	track, _ := db.FindTrackByEntry(result)
	if track != nil {
		player.PlayTrack(track)
		return
	}
	go func() {
		videoInfo, err := fetchVideoInfo(result, false)
		if err != nil {
			utils.HandleError(err, "Cannot fetch video info")
		}

		var trackInfo *providers.TrackInfo
		if videoInfo.Artist == "" || videoInfo.Track == "" {
			artistInfo := providers.ArtistInfo{
				Name: videoInfo.Uploader,
				MBID: "!YT:" + videoInfo.UploaderID,
			}
			albumInfo := providers.AlbumInfo{
				Title:    "YouTube videos by " + videoInfo.Uploader,
				MBID:     "!YT:" + videoInfo.UploaderID,
				ImageURL: videoInfo.GetThumbnail(),
			}
			trackInfo = &providers.TrackInfo{
				Artist: &artistInfo,
				Album:  &albumInfo,
				Title:  videoInfo.Title,
				MBID:   "!YT:" + videoInfo.ID,
			}
		} else {
			trackInfo, err = providers.FetchTrackInfo(videoInfo)
			utils.HandleError(err, "Cannot fetch track info")
			if trackInfo.Album.ImageURL == "" {
				trackInfo.Album.ImageURL = videoInfo.GetThumbnail()
			}
		}

		track, err := db.StoreTrack(videoInfo, trackInfo, false)
		if err != nil {
			utils.HandleError(err, "Cannot store track to db")
		}

		// create album folder
		os.MkdirAll(track.Album.GetAlbumPath(), 0744)

		// download album art
		utils.DownloadFile(trackInfo.Album.ImageURL, track.Album.GetAlbumArtPath())

		// start
		player.PlayTrack(track)
	}()
}
