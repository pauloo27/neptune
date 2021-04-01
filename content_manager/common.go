package content_manager

import (
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func fetchVideoInfo(result *youtube.YoutubeEntry, download bool) (*youtube.VideoInfo, error) {
	player.State.Fetching = result
	var info *youtube.VideoInfo
	var err error
	if download {
		tmpFile := path.Join(utils.GetTmpFolder(), result.ID+".m4a")
		info, err = youtube.FetchInfoAndDownload(result, tmpFile)
		if err != nil {
			return nil, err
		}
	} else {
		info, err = youtube.FetchInfo(result)
		if err != nil {
			return nil, err
		}
	}
	return info, nil
}

func store(result *youtube.YoutubeEntry, download bool) {
	player.ClearQueue()

	track, _ := db.FindTrackByEntry(result)
	if track != nil {
		player.PlayTrack(track)
		return
	}
	go func() {
		videoInfo, err := fetchVideoInfo(result, download)
		utils.HandleError(err, "Cannot fetch video info")

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

		track, err := db.StoreTrack(videoInfo, trackInfo, download)
		utils.HandleError(err, "Cannot store track to db")

		// create album folder
		os.MkdirAll(track.Album.GetAlbumPath(), 0744)

		// download album art
		utils.DownloadFile(trackInfo.Album.ImageURL, track.Album.GetAlbumArtPath())

		// move file (download to the temp folder) to the album folder
		err = os.Rename(path.Join(utils.GetTmpFolder(), result.ID+".m4a"), track.GetLocalPath())
		utils.HandleError(err, "Cannot move tmp download file to cache file")

		// start
		player.PlayTrack(track)
	}()
}
