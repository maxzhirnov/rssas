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

	rssParser := rssparser.NewRSSParser(conf.RSSFeeds)

	mongoStorage, _ := storage.NewMongoStorage(conf.MongoConn(), "rss2", "items")
	defer mongoStorage.Close()
	if err := mongoStorage.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	repository := repo.NewRepo(mongoStorage)

	app := service.NewApp(repository, rssParser)
	_ = app.StartFeedParserWorker(3)

	server := api.NewServer(app)
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	block := make(chan struct{})
	<-block
}
