package utils

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config stores application configurations
type Config struct {
	LogLevel        string `envconfig:"LOGLEVEL" default:"debug"`
	LogFormat       string `envconfig:"LOGFORMAT" default:"text"`
	BitbucketPAT    string `envconfig:"BITBUCKET_PAT" required:"true"`
	BitbucketServer string `envconfig:"BITBUCKET_SERVER" required:"true"`
	SlackAppToken   string `envconfig:"SLACK_APP_TOKEN" required:"true"`
	SLackBotToken   string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
	ChannelRegex    string `envconfig:"CHANNEL_REGEX" default:"^devops-([a-zA-Z0-9_]+)$"`
}

// ConfigFromEnvironment loads config from env variables and .env file
func ConfigFromEnvironment(path string) (Config, error) {
	// we do not care if there is no .env file.
	_ = godotenv.Overload(path)

	var s Config
	err := envconfig.Process("", &s)
	if err != nil {
		return s, err
	}

	return s, nil
}
