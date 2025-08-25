package crawler

import (
	"io"
	"log"
	"strings"
	"sync"

	"wb_tech/l2_16/internal/downloader"
	"wb_tech/l2_16/internal/parser"
	"wb_tech/l2_16/internal/storage"
)

type Task struct {
	URL   string
	Depth int
}

type Crawler struct {
	startURL   string
	maxDepth   int
	baseHost   string
	maxWorkers int

	tasks   chan Task
	wg      *sync.WaitGroup
	visited *sync.Map
}

func New(startURL, baseHost string, maxDepth, maxWorkers int) *Crawler {
	return &Crawler{
		startURL:   startURL,
		maxDepth:   maxDepth,
		baseHost:   baseHost,
		maxWorkers: maxWorkers,
		tasks:      make(chan Task, 100),
		wg:         &sync.WaitGroup{},
		visited:    &sync.Map{},
	}
}

func (c *Crawler) Start() {
	log.Printf("Запускаем %d воркеров...", c.maxWorkers)
	for i := 0; i < c.maxWorkers; i++ {
		go c.worker()
	}

	c.wg.Add(1)
	c.tasks <- Task{URL: c.startURL, Depth: c.maxDepth}

	c.wg.Wait()
	close(c.tasks)
	log.Println("Зеркалирование сайта завершено.")
}

func (c *Crawler) worker() {
	for task := range c.tasks {
		c.crawl(task.URL, task.Depth)
	}
}

func (c *Crawler) crawl(rawURL string, depth int) {
	defer c.wg.Done()

	if depth < 0 {
		return
	}

	if _, loaded := c.visited.LoadOrStore(rawURL, true); loaded {
		return
	}

	log.Printf("Глубина: %d | Скачиваем: %s", depth, rawURL)

	resp, err := downloader.Fetch(rawURL)
	if err != nil {
		log.Printf("Ошибка скачивания %s: %v", rawURL, err)
		return
	}
	defer resp.Body.Close()

	localPath, err := storage.URLToPath(rawURL)
	if err != nil {
		log.Printf("Некорректный URL для сохранения %s: %v", rawURL, err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа от %s: %v", rawURL, err)
		return
	}

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		result, err := parser.ProcessHTML(strings.NewReader(string(body)), rawURL, c.baseHost)
		if err != nil {
			log.Printf("Ошибка парсинга HTML с %s: %v", rawURL, err)
		} else {
			body = result.ModifiedHTML
			for _, link := range result.NewLinks {
				if _, loaded := c.visited.Load(link); !loaded {
					c.wg.Add(1)
					c.tasks <- Task{URL: link, Depth: depth - 1}
				}
			}
		}
	}

	if err := storage.Save(localPath, body); err != nil {
		log.Printf("Ошибка сохранения файла %s: %v", localPath, err)
	}
}
