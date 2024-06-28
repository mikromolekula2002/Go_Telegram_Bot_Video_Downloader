package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Функция для удаления файла
func DeleteVideo(fileName string) error {
	err := os.Remove(fileName)
	log.Println("Видео успешно удалено с хранилища.\n.")
	return fmt.Errorf("\nVideoManager.go - DeleteVideo: %w", err)
}

func GetVideo(videourl, videoID string) (string, error) {
	// Создаем файл для записи
	file, err := os.Create("./Videos/" + videoID + ".mp4")
	if err != nil {
		return "", fmt.Errorf("\nVideoManager.go - GetVideo: Ошибка создания файла.\n Error: %w", err)
	}
	defer file.Close()

	// Отправляем GET запрос по URL видео
	resp, err := http.Get(videourl)
	if err != nil {
		return "", fmt.Errorf("\nVideoManager.go - GetVideo: Ошибка GET запроса на источник: %s.\n Error: %w", videourl, err)
	}
	defer resp.Body.Close()

	// Проверяем статус код ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("\nVideoManager.go - GetVideo: Неверный статус код: %d", resp.StatusCode)
	}

	// Копируем содержимое ответа в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("\nVideoManager.go - GetVideo: Ошибка записи видео в файл.\n Error: %w", err)
	}

	log.Println("Видео успешно скачано в хранилище.")

	return file.Name(), nil
}
