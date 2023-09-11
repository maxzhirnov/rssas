package models

type Feed struct {
	Title string
	Link  string
}

func NewFeed(feedTitle, feedURL string) *Feed {
	return &Feed{
		Title: feedTitle,
		Link:  feedURL,
	}
}
