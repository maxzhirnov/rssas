package service

import (
	"time"

	"github.com/mmcdole/gofeed"

	"rssas/internal/log"
)

type repo interface {
	SaveItems(feed *gofeed.Feed) error
	SaveFeed(feed *gofeed.Feed, feedURL string) error
	LoadFeeds() ([]string, error)
}

type parser interface {
	ParseAll() error
	ParsedFeeds() []*gofeed.Feed
	ParseFeed(string) (*gofeed.Feed, error)
	AddFeeds([]string)
}

type App struct {
	repo   repo
	parser parser
	logger *log.Logger
}

func NewApp(repo repo, parser parser, logger *log.Logger) *App {
	return &App{
		repo:   repo,
		parser: parser,
		logger: logger,
	}
}

func (app App) ParseAllFeeds() error {
	app.logger.Log.Info("parsing all feeds")
	err := app.parser.ParseAll()
	if err != nil {
		app.logger.Log.Error(err)
		return err
	}

	feeds := app.parser.ParsedFeeds()

	for _, feed := range feeds {
		err := app.repo.SaveItems(feed)
		if err != nil {
			app.logger.Log.Error(err)
			return err
		}
	}
	return nil
}

func (app App) AddNewFeed(feedURL string) error {
	feed, err := app.parser.ParseFeed(feedURL)
	if err != nil {
		app.logger.Log.Error(err)
		return err
	}

	if err := app.repo.SaveFeed(feed, feedURL); err != nil {
		app.logger.Log.Error(err)
		return err
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
					app.logger.Log.Error(err)
				} else {
					app.logger.Log.Info("Parsed existing feed on schedule")
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

func (app App) StartFeedListUpdater(minutes int) (stopFunc func()) {
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Update feeds in parser
				feeds, err := app.repo.LoadFeeds()
				app.parser.AddFeeds(feeds)
				if err != nil {
					app.logger.Log.Error(err)
				} else {
					app.logger.Log.Infof("Added new feeds to parser: %s", feeds)
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
