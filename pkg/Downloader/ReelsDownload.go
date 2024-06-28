package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type ReelsResponse struct {
	DownloadLink string `json:"video_url"`
}

func GetReels(urlShare, ApiKey string) (string, error) {
	shortCode, err := ShortCodeFinder(urlShare)
	if err != nil {
		return "", err
	}

	videoUrl, err := GetReelsUrl(shortCode, ApiKey)
	if err != nil {
		return "", err
	}

	fileName, err := GetVideo(videoUrl, shortCode)
	if err != nil {
		return "", err
	}
	//получили наш видос, пока что не знаю как
	// и нужно его отправить назад боту, чтобы он переслал его пользователю
	return fileName, nil
}

func GetReelsUrl(shortCode, ApiKey string) (string, error) {
	log.Println("Получение прямой ссылки на Instagram Reels видео.")
	errPath := "ReelsDownload.go - GetReelsURL."

	url := "https://instagram-scraper-2022.p.rapidapi.com/ig/post_info/?shortcode=" + shortCode

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", "instagram-scraper-2022.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("\n%s\nОшибка запроса на API по URL: %s.\n Error: %w", errPath, url, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("\n%s\nОшибка записи видео в файл: %w", errPath, err)
	}

	// Declare a variable of type TikTokResponse
	var response ReelsResponse

	// Unmarshal the JSON into the variable
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", fmt.Errorf("\n%s\nОшибка анмаршалинга json: %s.\n Error: %w", errPath, body, err)
	}

	// Print the no_watermark_link
	return response.DownloadLink, nil
}

func ShortCodeFinder(url string) (string, error) {
	log.Println("Получение ShortCode на Reels видео.")
	errPath := "ReelsDownload.go - GetShortCodeFinder."

	re := regexp.MustCompile(`reel/([^/?]+)`)

	// Поиск совпадений
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("\n%s\nНе найдено прямой ссылки для скачивания Reels", errPath)
	}

	// Возвращаем первый найденный идентификатор
	return (matches[1]), nil
}
