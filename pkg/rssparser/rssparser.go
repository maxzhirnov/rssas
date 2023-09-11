package rssparser

import (
	"sync"

	"github.com/mmcdole/gofeed"
)

type RSSParser struct {
	feedURLs    []string
	parser      *gofeed.Parser
	parsedFeeds []*gofeed.Feed
	mu          sync.Mutex
}

func NewRSSParser(feedURLs []string) *RSSParser {
	return &RSSParser{
		parser:   gofeed.NewParser(),
		feedURLs: feedURLs,
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
			if err != nil {
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
		return nil, err
	}
	return feed, nil
}
