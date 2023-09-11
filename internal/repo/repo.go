package repo

import (
	"github.com/mmcdole/gofeed"

	"rssas/internal/models"
)

type storage interface {
	InsertMany(document []interface{}) error
	Bootstrap() error
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
	return r.storage.InsertMany(items)
}
