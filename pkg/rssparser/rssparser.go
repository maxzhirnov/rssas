package rssparser

import (
	"sync"

	"github.com/mmcdole/gofeed"

	"rssas/internal/log"
)

type RSSParser struct {
	feedURLs    []string
	parser      *gofeed.Parser
	parsedFeeds []*gofeed.Feed
	logger      *log.Logger
	mu          sync.Mutex
}

func NewRSSParser(feedURLs []string, logger *log.Logger) *RSSParser {
	return &RSSParser{
		parser:   gofeed.NewParser(),
		feedURLs: feedURLs,
		logger:   logger,
	}
}

func (p *RSSParser) ParseAll() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(p.feedURLs))

	for _, u := range p.feedURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			feed, err := p.parser.ParseURL(url)
			p.logger.Log.Infof("Parsed feed with title: %s", feed.Title)
			if err != nil {
				p.logger.Log.Error(err)
				errChan <- err
				return
			}

			p.mu.Lock()
			p.parsedFeeds = append(p.parsedFeeds, feed)
			p.mu.Unlock()
		}(u)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			p.logger.Log.Error(err)
			return err
		}
	}
	return nil
}

func (p *RSSParser) ParsedFeeds() []*gofeed.Feed {
	return p.parsedFeeds
}

func (p *RSSParser) ParseFeed(feedURL string) (*gofeed.Feed, error) {
	feed, err := p.parser.ParseURL(feedURL)
	if err != nil {
		p.logger.Log.Error(err)
		return nil, err
	}
	return feed, nil
}

func (p *RSSParser) AddFeeds(newFeeds []string) {
	existingItemsMap := make(map[string]struct{})

	// Заполняем карту существующими элементами.
	for _, item := range p.feedURLs {
		existingItemsMap[item] = struct{}{}
	}

	// Проверяем каждый новый элемент на наличие в карте.
	for _, newItem := range newFeeds {
		if _, exists := existingItemsMap[newItem]; !exists {
			p.feedURLs = append(p.feedURLs, newItem)
			existingItemsMap[newItem] = struct{}{}
		}
	}
}
