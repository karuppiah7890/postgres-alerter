package config

import (
	"fmt"
	"os"
)

// All configuration is through environment variables

const POSTGRES_NAME_ENV_VAR = "POSTGRES_NAME"
const DEFAULT_POSTGRES_NAME = "Postgres"
const POSTGRES_URI_ENV_VAR = "POSTGRES_URI"
const ENVIRONMENT_NAME_ENV_VAR = "ENVIRONMENT_NAME"
const DEFAULT_ENVIRONMENT_NAME = "Production"
const SLACK_TOKEN_ENV_VAR = "SLACK_TOKEN"
const SLACK_CHANNEL_ENV_VAR = "SLACK_CHANNEL"

type Config struct {
	postgresName    string
	postgresUri     string
	environmentName string
	slackToken      string
	slackChannel    string
}

func NewConfigFromEnvVars() (*Config, error) {
	postgresName := getPostgresName()

	postgresUri, err := getPostgresUri()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting postgres uri: %v", err)
	}

	environmentName := getEnvironmentName()

	slackToken, err := getSlackToken()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting slack token: %v", err)
	}

	slackChannel, err := getSlackChannel()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting slack channel: %v", err)
	}

	return &Config{
		postgresName:    postgresName,
		postgresUri:     postgresUri,
		environmentName: environmentName,
		slackToken:      slackToken,
		slackChannel:    slackChannel,
	}, nil
}

// Get optional name for the Postgres instance. Default is "Postgres".
// This will be used in the alert messages
func getPostgresName() string {
	postgresName, ok := os.LookupEnv(POSTGRES_NAME_ENV_VAR)
	if !ok {
		return DEFAULT_POSTGRES_NAME
	}

	return fmt.Sprintf("%s (Postgres)", postgresName)
}

// Get postgres uri
func getPostgresUri() (string, error) {
	postgresUri, ok := os.LookupEnv(POSTGRES_URI_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable is not defined and is required. Please define it", POSTGRES_URI_ENV_VAR)
	}

	return postgresUri, nil
}

// Get optional environment name for the environment where
// the services are running. Default is "Production". This name will
// be used in the alert messages
func getEnvironmentName() string {
	environmentName, ok := os.LookupEnv(ENVIRONMENT_NAME_ENV_VAR)
	if !ok {
		return DEFAULT_ENVIRONMENT_NAME
	}

	return environmentName
}

func getSlackToken() (string, error) {
	slackToken, ok := os.LookupEnv(SLACK_TOKEN_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable is not defined and is required. Please define it", SLACK_TOKEN_ENV_VAR)
	}
	return slackToken, nil
}

func getSlackChannel() (string, error) {
	slackChannel, ok := os.LookupEnv(SLACK_CHANNEL_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable is not defined and is required. Please define it", SLACK_CHANNEL_ENV_VAR)
	}
	return slackChannel, nil
}

func (c *Config) GetPostgresName() string {
	return c.postgresName
}

func (c *Config) GetPostgresUri() string {
	return c.postgresUri
}

func (c *Config) GetEnvironmentName() string {
	return c.environmentName
}

func (c *Config) GetSlackToken() string {
	return c.slackToken
}

func (c *Config) GetSlackChanel() string {
	return c.slackChannel
}
