package repro_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	repro "github.com/abayer/testcontainers-gotest-log-missing"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAddNumbers(t *testing.T) {
	_, closeFunc := DBForTest(t)
	defer closeFunc()

	t.Log("Sleeping for 15 seconds before adding numbers")
	time.Sleep(15 * time.Second)
	t.Log("Done sleeping, let's add")
	result := repro.AddNumbers(1, 2)

	if result != 3 {
		t.Errorf("%d should have been 3", result)
	}
}

// DBForTest spins up a postgres container, creates the test database on it, migrates it, and returns the db and a close function
func DBForTest(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	// container and database
	container, db, err := CreateTestContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	closeFunc := func() {
		_ = db.Close()
		_ = container.Terminate(ctx)
	}

	return db, closeFunc
}

// CreateTestContainer spins up a Postgres database container
func CreateTestContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       "some-db",
	}
	dockerPort := "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", port.Port(), "some-db")
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:12",
			ExposedPorts: []string{dockerPort},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(nat.Port(dockerPort), "postgres", dbURL).Timeout(time.Second * 30),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(dockerPort))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at dockerPort: ", mappedPort)

	url := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), "some-db")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}
