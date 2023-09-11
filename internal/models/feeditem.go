package models

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Feed      string       `bson:"feed"`
	Title     string       `bson:"title"`
	Guid      string       `bson:"guid"`
	Link      string       `bson:"link"`
	Published time.Time    `bson:"published"`
	Parsed    time.Time    `bson:"parsed"`
	Data      *gofeed.Item `bson:"data"`
}

func NewFeedItem(feedTitle string, goFeedItem *gofeed.Item) *FeedItem {
	return &FeedItem{
		Feed:      feedTitle,
		Title:     goFeedItem.Title,
		Guid:      goFeedItem.GUID,
		Link:      goFeedItem.Link,
		Published: *goFeedItem.PublishedParsed,
		Parsed:    time.Now(),
		Data:      goFeedItem,
	}
}
