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
		log.Error(err)
		return err
	}

	openaiToken, ok := os.LookupEnv("OPENAI_SECRET")
	if !ok {
		log.Warn(fmt.Errorf("openai token haven't been found in env"))
	}

	mongoConn, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		err := fmt.Errorf("mongo conn string haven't been found in env")
		log.Warn(err)
		return err
	}

	c.openaiToken = openaiToken
	c.mongoConn = mongoConn

	return nil
}

func (c *Config) MongoConn() string {
	return c.mongoConn
}
