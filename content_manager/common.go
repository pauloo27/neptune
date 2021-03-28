package content_manager

import (
	"path"

	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/providers/youtube"
	"github.com/Pauloo27/neptune/utils"
)

func fetchVideoInfo(result *youtube.YoutubeEntry, download bool) (*youtube.VideoInfo, error) {
	player.State.Fetching = result
	var info *youtube.VideoInfo
	var err error
	if download {
		tmpFile := path.Join(utils.GetTmpFolder(), result.ID)
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
