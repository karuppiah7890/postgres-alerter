package main

import (
	"fmt"
	"log"

	"github.com/karuppiah7890/postgres-alerter/pkg/config"
	"github.com/karuppiah7890/postgres-alerter/pkg/postgres"
	"github.com/karuppiah7890/postgres-alerter/pkg/slack"
)

func main() {
	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("error occurred while getting configuration from environment variables: %v", err)
	}

	// TODO: Use Mocks to test the integration with ease for different cases with unit tests
	postgresStatus := postgres.GetPostgresStatus(c.GetPostgresUri())

	if !postgresStatus.IsUp {
		message := fmt.Sprintf("Critical alert :rotating_light:! %s is down in %s environment :rotating_light:", c.GetPostgresName(), c.GetEnvironmentName())
		// TODO: Use Mocks to test the integration with ease for different cases with unit tests
		err := slack.SendMessage(c.GetSlackToken(), c.GetSlackChanel(), message)
		if err != nil {
			log.Fatalf("error occurred while sending slack alert message: %v", err)
		}
	}
}
