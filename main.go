package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/karuppiah7890/postgres-alerter/pkg/config"
	"github.com/karuppiah7890/postgres-alerter/pkg/postgres"
	"github.com/karuppiah7890/postgres-alerter/pkg/slack"
)

func main() {
	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	go checkSignal(signals, done)

	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("error occurred while getting configuration from environment variables: %v", err)
	}

	client := postgres.NewClient(c.GetPostgresUri())

	ticker := time.NewTicker(time.Second)

	keepGoing := true

	for keepGoing {
		select {
		case <-ticker.C:
			postgresStatus := client.GetPostgresStatus()
			alertAboutPgStatus(postgresStatus, c)

		case <-done:
			keepGoing = false
		}
	}

}

func checkSignal(signals chan os.Signal, done chan bool) {
	<-signals
	done <- true
	os.Exit(0)
}

func alertAboutPgStatus(postgresStatus postgres.PostgresStatus, c *config.Config) {
	// TODO: Use Mocks to test the integration with ease for different cases with unit tests

	if !postgresStatus.IsUp {
		message := fmt.Sprintf("Critical alert :rotating_light:! %s is down in %s environment :rotating_light:", c.GetPostgresName(), c.GetEnvironmentName())
		// TODO: Use Mocks to test the integration with ease for different cases with unit tests
		err := slack.SendMessage(c.GetSlackToken(), c.GetSlackChanel(), message)
		if err != nil {
			log.Fatalf("error occurred while sending slack alert message: %v", err)
		}
	}
}
