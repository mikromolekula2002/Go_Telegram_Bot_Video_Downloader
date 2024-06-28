package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ShortsLink struct {
	Quality string `json:"quality"`
	Link    string `json:"link"`
}

type ShortsResponse struct {
	Links    []ShortsLink `json:"links"`
	ShortsId string       `json:"id"`
}

func GetShorts(urlShare, ApiKey string) (string, error) {

	videoUrl, videoID, err := GetShortsUrl(urlShare, ApiKey)
	if err != nil {
		return "", err
	}

	fileName, err := GetVideo(videoUrl, videoID)
	if err != nil {
		return "", err
	}
	//получили наш видос, пока что не знаю как
	// и нужно его отправить назад боту, чтобы он переслал его пользователю
	return fileName, nil
}

func GetShortsUrl(urlShare, ApiKey string) (string, string, error) {
	log.Println("Получение прямой ссылки на Youtube Shorts видео.")
	errPath := "ShortsDownload.go - GetShortsURL."

	url := "https://youtube-video-and-shorts-downloader1.p.rapidapi.com/api/getYTVideo?url=" + urlShare

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", "youtube-video-and-shorts-downloader1.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("\n%s\nОшибка запроса на API по URL: %s.\n Error: %w", errPath, url, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("\n%s\nОшибка записи видео в файл: %w", errPath, err)
	}

	// Declare a variable of type TikTokResponse
	var response ShortsResponse

	// Unmarshal the JSON into the variable
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", "", fmt.Errorf("\n%s\nОшибка анмаршалинга json: %s.\n Error: %w", errPath, body, err)
	}

	// Ищем ссылку с качеством "render_1080p", в крайнем случае "render_720p"
	var selectedLink string

	for _, link := range response.Links {
		if link.Quality == "render_1080p" {
			selectedLink = link.Link
			break
		} else if link.Quality == "render_720p" {
			selectedLink = link.Link
			break
		}
	}

	if selectedLink == "" {
		return "", "", fmt.Errorf("\n%s\n Ошибка: не найдено ссылки на скачивание видео", errPath)
	}
	// Print the no_watermark_link
	return selectedLink, response.ShortsId, nil
}
