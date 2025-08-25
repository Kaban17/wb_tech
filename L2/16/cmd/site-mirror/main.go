package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"wb_tech/l2_16/internal/crawler"
)

const maxWorkers = 5

func main() {
	startURL := flag.String("url", "", "Стартовый URL для скачивания (обязательный)")
	depth := flag.Int("depth", 2, "Глубина рекурсивного скачивания")
	flag.Parse()

	if *startURL == "" {
		log.Fatal("Необходимо указать URL с помощью флага -url")
	}

	parsedURL, err := url.Parse(*startURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https" ) {
		log.Fatalf("Некорректный URL: %s", *startURL)
	}
	baseHost := parsedURL.Host

	// Создаем корневую директорию для сохранения сайта.
	if err := os.MkdirAll(baseHost, 0755); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", baseHost, err)
	}

	// Создаем и запускаем новый кр��улер.
	c := crawler.New(*startURL, baseHost, *depth, maxWorkers)
	c.Start()
}
