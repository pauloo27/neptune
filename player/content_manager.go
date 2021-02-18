package player

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func PlayTracks(tracks []*db.Track) {
	ClearQueue()

	if len(tracks) == 0 {
		return
	}

	AddToTopOfQueue(tracks[0])
	LoadFile(tracks[0].GetPath())

	for _, track := range tracks[1:] {
		AddToQueue(track)
		AppendFile(track.GetPath())
	}
}

func PlayTrack(track *db.Track) {
	ClearQueue()

	filePath := track.GetPath()
	err := db.PlayTrack(track)
	utils.HandleError(err, "Cannot play track")

	AddToTopOfQueue(track)
	LoadFile(filePath)
}

func PlayResult(result *youtube.YoutubeEntry) {
	ClearQueue()

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
				var err error
				trackInfo, err = providers.FetchTrackInfo(videoInfo)
				utils.HandleError(err, "Cannot fetch track info")
				if trackInfo.Album.ImageURL == "" {
					trackInfo.Album.ImageURL = videoInfo.GetThumbnail()
				}
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

			AddToTopOfQueue(track)
			LoadFile(filePath)
		}()
	} else {
		AddToTopOfQueue(track)
		LoadFile(track.GetPath())
	}
	State.Fetching = nil
}
