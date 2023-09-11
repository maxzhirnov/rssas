package service

import (
	"time"

	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

type repo interface {
	SaveItems(feed *gofeed.Feed) error
	SaveFeed(feed *gofeed.Feed) error
}

type parser interface {
	ParseAll() error
	ParsedFeeds() []*gofeed.Feed
	ParseFeed(string) (*gofeed.Feed, error)
}

type App struct {
	repo   repo
	parser parser
}

func NewApp(repo repo, parser parser) *App {
	return &App{
		repo:   repo,
		parser: parser,
	}
}

func (app App) ParseAllFeeds() error {
	err := app.parser.ParseAll()
	if err != nil {
		return err
	}

	feeds := app.parser.ParsedFeeds()

	for _, feed := range feeds {
		err := app.repo.SaveItems(feed)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app App) StartFeedParserWorker(hours int) (stopFunc func()) {
	ticker := time.NewTicker(time.Duration(hours) * time.Hour)
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := app.ParseAllFeeds()
				if err != nil {
					log.Error(err)
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(stop)
	}
}

func (app App) AddNewFeed(feedURL string) error {
	feed, err := app.parser.ParseFeed(feedURL)
	if err != nil {
		return err
	}

	if err := app.repo.SaveFeed(feed); err != nil {
		return err
	}
	return nil
}
