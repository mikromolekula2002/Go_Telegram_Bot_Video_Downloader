package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TikTokResponse struct {
	NoWatermarkLink string `json:"no_watermark_link"`
	AwemeID         string `json:"aweme_id"`
}

func GetTiktok(urlShare, ApiKey string) (string, error) {

	videoUrl, videoID, err := GetTikTokUrl(urlShare, ApiKey)
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

func GetTikTokUrl(urlShare, ApiKey string) (string, string, error) {
	log.Println("Получение прямой ссылки на тикток видео.")

	errPath := "TikTokDownload.go - GetTikTokURL."
	url := "https://scraptik.p.rapidapi.com/video-without-watermark?url=" + urlShare
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", "scraptik.p.rapidapi.com")

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
	var response TikTokResponse

	// Unmarshal the JSON into the variable
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", "", fmt.Errorf("\n%s\nОшибка анмаршалинга json: %s.\n Error: %w", errPath, body, err)
	}

	// Print the no_watermark_link
	return response.NoWatermarkLink, response.AwemeID, nil
}
