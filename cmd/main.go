package main

import (
	log "github.com/sirupsen/logrus"

	"rssas/internal/api"
	"rssas/internal/config"
	"rssas/internal/repo"
	"rssas/internal/service"
	"rssas/internal/storage"
	"rssas/pkg/rssparser"
)

func main() {
	conf := config.New()
	if err := conf.Load(); err != nil {
		log.Fatal(err)
	}

	mongoStorage, err := storage.NewMongoStorage(conf.MongoConn(), "rss2")
	if err != nil {
		log.Fatal(err)
	}

	repository := repo.NewRepo(mongoStorage)
	defer repository.Close()

	if err := repository.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	feeds, err := repository.LoadFeeds()
	if err != nil {
		log.Error(err)
	}

	conf.AddFeeds(feeds)
	log.Infof("feeds: %s'n", conf.RSSFeeds)

	rssParser := rssparser.NewRSSParser(conf.RSSFeeds)
	app := service.NewApp(repository, rssParser)
	_ = app.StartFeedParserWorker(1)

	server := api.NewServer(app)
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	block := make(chan struct{})
	<-block
}
