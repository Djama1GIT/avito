package repository_test

import (
	"avito/pkg/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var host string
var port string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=testdb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/testdb?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120)
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	host = resource.GetBoundIP("5432/tcp")
	port = resource.GetPort("5432/tcp")
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
func TestPostgres(t *testing.T) {
	tests := []struct {
		name     string
		config   repository.Config
		expected error
	}{
		{
			name: "ValidConfig_Success",
			config: repository.Config{
				Driver:   "postgres",
				Host:     host,
				Port:     port,
				Username: "user_name",
				Password: "secret",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expected: nil,
		},
		{
			name: "InvalidConfig",
			config: repository.Config{
				Driver:   "postgres",
				Host:     host,
				Port:     port,
				Username: "testuserinvalid123",
				Password: "testpassword",
				SSLMode:  "enable",
			},
			expected: errors.New("pq: SSL is not enabled on the server"),
		},
		{
			name: "ConnectionError",
			config: repository.Config{
				Driver:   "postgres",
				Host:     "invalidhost",
				Port:     "5432",
				Username: "testuser",
				Password: "testpassword",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expected: errors.New("dial tcp: lookup invalidhost: no such host"),
		},
		{
			name: "OpenError",
			config: repository.Config{
				Driver:   "postgresanet",
				Host:     "invalidhost",
				Port:     "5432",
				Username: "testuser",
				Password: "testpassword",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expected: errors.New("sql: unknown driver \"postgresanet\" (forgotten import?)"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := repository.NewDB(test.config)
			if test.expected != nil {
				assert.Error(t, err, "Expected an error")
				assert.EqualError(t, err, test.expected.Error(), "Unexpected error message")
			} else {
				assert.NoError(t, err, "Failed to create a new Postgres DB connection")
			}

			if db != nil {
				db.Close()
			}
		})
	}
}
