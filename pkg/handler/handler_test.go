package handler

import (
	"avito/pkg/repository"
	"avito/pkg/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestInitRoutes(t *testing.T) {
	db, err := repository.NewDB(repository.Config{
		Driver:   "postgres",
		Host:     host,
		Port:     port,
		Username: "user_name",
		Password: "secret",
		Name:     "testdb",
		SSLMode:  "disable",
	})
	assert.NoError(t, err)

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := NewHandler(services)
	router := handler.InitRoutes()

	testRequest(t, router, "GET", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "POST", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "PATCH", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "DELETE", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "GET", "/api/users/history/", http.StatusBadRequest)
	testRequest(t, router, "DELETE", "/api/users/expired-segments/", http.StatusInternalServerError)
}

func testRequest(t *testing.T, router http.Handler, method, url string, expectedStatusCode int) {
	req, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedStatusCode {
		t.Errorf("expected status code %d for %s request to %s, got %d", expectedStatusCode, method, url, w.Code)
	}
}
