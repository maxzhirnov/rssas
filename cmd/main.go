package main

import (
	"os"

	"rssas/internal/api"
	"rssas/internal/config"
	"rssas/internal/log"
	"rssas/internal/repo"
	"rssas/internal/service"
	"rssas/internal/storage"
	"rssas/pkg/rssparser"
)

func main() {
	logger, err := log.NewLogger()
	if err != nil {
		os.Exit(0)
	}
	defer logger.LogFile.Close()
	logger.Log.Info("Starting...")

	conf := config.New(logger)
	if err := conf.Load(); err != nil {
		logger.Log.Fatal(err)
	}

	mongoStorage, err := storage.NewMongoStorage(conf.MongoConn(), "rss2", logger)
	if err != nil {
		logger.Log.Fatal(err)
	}

	repository := repo.NewRepo(mongoStorage, logger)
	defer repository.Close()

	if err := repository.Bootstrap(); err != nil {
		logger.Log.Fatal(err)
	}

	feeds, err := repository.LoadFeeds()
	if err != nil {
		logger.Log.Error(err)
	}

	rssParser := rssparser.NewRSSParser(conf.RSSFeeds, logger)
	rssParser.AddFeeds(feeds)

	app := service.NewApp(repository, rssParser, logger)
	_ = app.StartFeedParserWorker(1)
	_ = app.StartFeedListUpdater(15)

	server := api.NewServer(app)
	if err := server.Run(":8080"); err != nil {
		logger.Log.Fatal(err)
	}

	block := make(chan struct{})
	<-block
}
