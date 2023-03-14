package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/karuppiah7890/postgres-alerter/pkg/config"
	"github.com/karuppiah7890/postgres-alerter/pkg/postgres"
	"github.com/karuppiah7890/postgres-alerter/pkg/slack"
	"github.com/karuppiah7890/postgres-alerter/state"
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

	oldState, err := state.New(c.GetStateFilePath())
	if err != nil {
		log.Fatalf("error occurred while initializing state from state file at %s: %v", c.GetStateFilePath(), err)
	}

	lastThreadTimestamp := oldState.LastThreadTimestamp

	// TODO: Send message only if the time diff between last sent message and current message is greater than min time T
	if oldState.SendAlert(postgresStatus.IsUp) {
		var message string

		if postgresStatus.IsUp {
			message = fmt.Sprintf("%s is back up in %s environment!", c.GetPostgresName(), c.GetEnvironmentName())
		} else {
			message = fmt.Sprintf("Critical alert :rotating_light:! %s is down in %s environment :rotating_light:", c.GetPostgresName(), c.GetEnvironmentName())
		}

		if createNewThread(lastThreadTimestamp, time.Now(), c.GetNewThreadMinInterval()) {
			// TODO: Use Mocks to test the integration with ease for different cases with unit tests
			lastThreadTimestamp, err = slack.SendMessage(c.GetSlackToken(), c.GetSlackChanel(), message)
			if err != nil {
				log.Fatalf("error occurred while sending slack alert message: %v", err)
			}
		} else {
			// ignore the existing thread's new message's timestamp
			_, err = slack.SendMessageToThread(c.GetSlackToken(), c.GetSlackChanel(), message, lastThreadTimestamp)
			if err != nil {
				log.Fatalf("error occurred while sending slack alert message: %v", err)
			}
		}
	}

	// store current state
	newState := state.State{
		PostgresIsUp:        postgresStatus.IsUp,
		LastThreadTimestamp: lastThreadTimestamp,
	}

	err = newState.StoreToFile(c.GetStateFilePath())
	if err != nil {
		log.Fatalf("error occurred while storing new state to state file at %s: %v", c.GetStateFilePath(), err)
	}
}

func createNewThread(lastThreadTimestampStr string, now time.Time, newThreadMinInterval time.Duration) bool {
	if lastThreadTimestampStr == "" {
		return true
	}

	lastThreadTimestamp, err := strconv.ParseFloat(lastThreadTimestampStr, 64)
	if err != nil {
		log.Printf("error occurred while parsing last thread timestamp string value (%s) to float: %v", lastThreadTimestampStr, err)
		return true
	}

	lastThreadTime := time.Unix(int64(lastThreadTimestamp), 0)
	duration := now.Sub(lastThreadTime)

	return duration > newThreadMinInterval
}
