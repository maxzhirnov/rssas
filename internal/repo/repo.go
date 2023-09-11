package repo

import (
	"github.com/mmcdole/gofeed"

	"rssas/internal/models"
)

const (
	collNameItems = "items"
	collNameFeeds = "feeds"
)

type storage interface {
	InsertMany(document []interface{}, collection string) error
	InsertOne(document interface{}, collection string) error
	GetFeedsLinks() ([]string, error)
	Bootstrap() error
	Close() error
}

type Repo struct {
	storage storage
}

func NewRepo(storage storage) *Repo {
	return &Repo{
		storage: storage,
	}
}

func (r Repo) SaveItems(feed *gofeed.Feed) error {
	items := make([]interface{}, len(feed.Items))
	for i, p := range feed.Items {
		items[i] = models.NewFeedItem(feed.Title, p)
	}
	return r.storage.InsertMany(items, collNameItems)
}

func (r Repo) SaveFeed(feed *gofeed.Feed, feedURL string) error {
	if err := r.storage.InsertOne(models.NewFeed(feed.Title, feedURL), collNameFeeds); err != nil {
		return err
	}
	return nil
}

func (r Repo) LoadFeeds() ([]string, error) {
	return r.storage.GetFeedsLinks()
}

func (r Repo) Bootstrap() error {
	return r.storage.Bootstrap()
}

func (r Repo) Close() error {
	return r.storage.Close()
}
