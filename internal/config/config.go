package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	openaiToken string
	mongoConn   string
	RSSFeeds    []string `yaml:"rss_feeds"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Error(err)
		return err
	}

	if err := yaml.Unmarshal(configFile, c); err != nil {
		log.Error(err)
		return err
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Warn(err)
	}

	openaiToken, ok := os.LookupEnv("OPENAI_SECRET")
	if ok {
		c.openaiToken = openaiToken
	} else {
		err := fmt.Errorf("openai token haven't been found in env")
		log.Error(err)
	}

	mongoConn, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		err := fmt.Errorf("mongo conn string haven't been found in env")
		log.Error(err)
		return err
	}

	c.mongoConn = mongoConn

	return nil
}

func (c *Config) MongoConn() string {
	return c.mongoConn
}

func (c *Config) AddFeeds(newFeeds []string) {
	existingItemsMap := make(map[string]struct{})

	// Заполняем карту существующими элементами.
	for _, item := range c.RSSFeeds {
		existingItemsMap[item] = struct{}{}
	}

	// Проверяем каждый новый элемент на наличие в карте.
	for _, newItem := range newFeeds {
		if _, exists := existingItemsMap[newItem]; !exists {
			c.RSSFeeds = append(c.RSSFeeds, newItem)
			existingItemsMap[newItem] = struct{}{}
		}
	}
}
