package models

import (
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Title string
	Link  string
}

func NewFeed(goFeed *gofeed.Feed) *Feed {
	return &Feed{
		Title: goFeed.Title,
		Link:  goFeed.Link,
	}
}
