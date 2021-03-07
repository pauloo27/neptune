package player

import (
	"fmt"
	"os"
	"path"

	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/providers"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func PlayTracks(tracks []*db.Track) {
	clearQueue()

	if len(tracks) == 0 {
		return
	}

	addToTopOfQueue(tracks[0])
	loadFile(tracks[0].GetPath())

	for _, track := range tracks[1:] {
		addToQueue(track)
		appendFile(track.GetPath())
	}
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func PlayTrack(track *db.Track) {
	clearQueue()

	go func() {
		err := db.PlayTrack(track)
		utils.HandleError(err, "Cannot play track")
	}()

	addToTopOfQueue(track)
	loadFile(track.GetPath())
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func AddTrackToQueue(track *db.Track) {
	addToQueue(track)
	appendFile(track.GetPath())
	hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
}

func PlayResult(result *youtube.YoutubeEntry) {
	clearQueue()

	State.Fetching = result
	hook.CallHooks(hook.HOOK_RESULT_FETCH_STARTED, nil)
	track, err := db.PlayEntry(result)
	utils.HandleError(err, "Cannot find track")
	if track == nil {
		hook.CallHooks(hook.HOOK_RESULT_DOWNLOAD_STARTED, nil)
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

			addToTopOfQueue(track)
			loadFile(filePath)
			hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
		}()
	} else {
		addToTopOfQueue(track)
		loadFile(track.GetPath())
		hook.CallHooks(hook.HOOK_QUEUE_UPDATE_FINISHED)
	}
	State.Fetching = nil
}
