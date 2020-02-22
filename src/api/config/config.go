package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	apiGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
	LogLevel             = "LOG_LEVEL"
	goEnvironment        = "GO_ENVIRONMENT"
	production           = "production"
)

var (
	githubAccessToken string
	logLevel          string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	githubAccessToken = os.Getenv(apiGithubAccessToken)
	logLevel = os.Getenv(LogLevel)
}

func GetGithubAccessToken() string {
	return githubAccessToken
}

func GetLogLevel() string {
	return logLevel
}

func IsProduction() bool {
	return os.Getenv(goEnvironment) == production
}
