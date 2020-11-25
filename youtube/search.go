package youtube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/buger/jsonparser"
)

type YoutubeEntry struct {
	Title, Uploader, Duration, ID string
	Live                          bool
}

func (result *YoutubeEntry) URL() string {
	return fmt.Sprintf("https://youtube.com/watch?v=%s", result.ID)
}

func getContent(data []byte, index int) []byte {
	id := fmt.Sprintf("[%d]", index)
	contents, _, _, _ := jsonparser.Get(data, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents", id, "itemSectionRenderer", "contents")
	return contents
}

func SearchFor(searchTerm string, limit int) ([]*YoutubeEntry, error) {
	url := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(searchTerm))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Language", "en")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code %d, 200 expected", res.StatusCode)
	}

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := string(buffer)
	splittedScript := strings.Split(body, `window["ytInitialData"] = `)
	if len(splittedScript) != 2 {
		splittedScript = strings.Split(body, `var ytInitialData = `)
	}

	if len(splittedScript) != 2 {
		return nil, fmt.Errorf("Too much splitted scripts")
	}
	splittedScript = strings.Split(splittedScript[1], `window["ytInitialPlayerResponse"] = null;`)
	jsonData := []byte(splittedScript[0])

	index := 0
	var contents []byte

	for {
		contents = getContent(jsonData, index)
		_, _, _, err = jsonparser.Get(contents, "[0]", "carouselAdRenderer")

		if err == nil {
			index++
		} else {
			break
		}
	}

	results := []*YoutubeEntry{}

	_, err = jsonparser.ArrayEach(contents, func(value []byte, t jsonparser.ValueType, i int, err error) {
		if limit > 0 && len(results) >= limit {
			return
		}

		id, err := jsonparser.GetString(value, "videoRenderer", "videoId")
		if err != nil {
			return
		}

		title, err := jsonparser.GetString(value, "videoRenderer", "title", "runs", "[0]", "text")
		if err != nil {
			return
		}

		uploader, err := jsonparser.GetString(value, "videoRenderer", "ownerText", "runs", "[0]", "text")
		if err != nil {
			return
		}

		live := false
		duration, err := jsonparser.GetString(value, "videoRenderer", "lengthText", "simpleText")

		if err != nil {
			duration = ""
			live = true
		}

		results = append(results, &YoutubeEntry{
			Title:    title,
			Uploader: uploader,
			Duration: duration,
			ID:       id,
			Live:     live,
		})
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
