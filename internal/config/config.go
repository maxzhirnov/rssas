package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	"rssas/internal/log"
)

type Config struct {
	openaiToken string
	mongoConn   string
	RSSFeeds    []string `yaml:"rss_feeds"`
	logger      *log.Logger
}

func New(log *log.Logger) *Config {
	return &Config{
		logger: log,
	}
}

func (c *Config) Load() error {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		c.logger.Log.Error(err)
		return err
	}

	if err := yaml.Unmarshal(configFile, c); err != nil {
		c.logger.Log.Error(err)
		return err
	}

	if err := godotenv.Load(".env"); err != nil {
		c.logger.Log.Warn(err)
	}

	openaiToken, ok := os.LookupEnv("OPENAI_SECRET")
	if ok {
		c.openaiToken = openaiToken
	} else {
		err := fmt.Errorf("openai token haven't been found in env")
		c.logger.Log.Error(err)
	}

	mongoConn, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		err := fmt.Errorf("mongo conn string haven't been found in env")
		c.logger.Log.Error(err)
		return err
	}

	c.mongoConn = mongoConn

	return nil
}

func (c *Config) MongoConn() string {
	return c.mongoConn
}
