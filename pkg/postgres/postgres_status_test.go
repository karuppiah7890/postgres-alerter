package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/karuppiah7890/postgres-alerter/pkg/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetPostgresStatus(t *testing.T) {
	t.Run("Postgres is up", func(t *testing.T) {
		ctx := context.Background()
		username := "root"
		password := "password"
		dbName := "test"
		req := testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
			Env: map[string]string{
				"POSTGRES_USER":     username,
				"POSTGRES_PASSWORD": password,
				"POSTGRES_DB":       dbName,
			},
		}
		postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			t.Error(err)
		}
		defer func() {
			if terminateErr := postgresContainer.Terminate(ctx); terminateErr != nil {
				t.Fatalf("failed to terminate container: %s", terminateErr.Error())
			}
		}()

		host, err := postgresContainer.Host(ctx)
		if err != nil {
			t.Errorf("error occurred: expected no errors while getting test db host but got one: %v", err)
		}

		port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
		if err != nil {
			t.Errorf("error occurred: expected no errors while getting test db port but got one: %v", err)
		}

		uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port.Int(), dbName)

		postgresStatus := postgres.GetPostgresStatus(uri)

		if postgresStatus.IsUp == false {
			t.Errorf("error occurred: expected Postgres to be up and running, but Postgres is not up: %v", postgresStatus.Errors)
		}
	})

	t.Run("Postgres is down", func(t *testing.T) {
		postgresStatus := postgres.GetPostgresStatus("postgres://user:password@localhost:5432/db?sslmode=disable")

		if postgresStatus.IsUp == true {
			t.Error("error occurred: expected Postgres to be down, but the Postgres is up and running")
		}

		if postgresStatus.Errors == nil {
			t.Error("error occurred: expected atleast one error but got none")
		}
	})
}
