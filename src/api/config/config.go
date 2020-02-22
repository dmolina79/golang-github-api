package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	apiGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
)

var (
	githubAccessToken string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	githubAccessToken = os.Getenv(apiGithubAccessToken)
}

func GetGithubAccessToken() string {
	return githubAccessToken
}
