package integration

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	_ "github.com/lib/pq" //postgres driver
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbname = "timekeeper"
	dbuser = "postgres"
	dbpass = "password1!"
)

//go:embed schema.sql
var schema string

func Test(t *testing.T) {
	ctx := context.Background()

	// 1. Start the postgres databaseC and run any migrations on it
	t.Log("starting postgres")
	databaseC, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		testcontainers.WithLogConsumers(logger{}),
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpass),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("postgres"),
	)
	testcontainers.CleanupContainer(t, databaseC)
	require.NoError(t, err)

	// Run any migrations on the database
	t.Log("migrating schema")
	connstr, err := databaseC.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	_, _, err = databaseC.Exec(ctx, []string{"psql", "-U", dbuser, "-d", dbname, "-c", schema})
	require.NoError(t, err)

	t.Run("Ping Database", func(t *testing.T) {
		db, err := sql.Open("postgres", connstr)
		require.NoError(t, err)

		err = db.Ping()
		t.Log("pinged database:", connstr)
		require.NoError(t, err)

		_, err = db.Exec(`INSERT INTO timekeeper.users(login_name, password) VALUES ('test', 'pw')`)
		require.NoError(t, err)
		row := db.QueryRow(`SELECT login_name FROM timekeeper.users WHERE id = 1`)
		require.NoError(t, row.Err())
		var loginName string
		err = row.Scan(&loginName)
		require.NoError(t, err)
		t.Log(loginName)
	})

	// 2. Start timekeeper
	timekeeperC, err := testcontainers.Run(
		ctx, "ghcr.io/m4schini/timekeeper:"+getEnvOr("TIMEKEEPER_CONTAINER_TAG", "latest"),
		testcontainers.WithEnv(map[string]string{
			"PORT":                      "8080",
			"DATABASE_CONNECTIONSTRING": connstr,
			"TIMEKEEPER_ADMIN_PASSWORD": "password1!",
			"JWT_SECRET":                "secret",
		}),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.NetworkMode = "host"
		}),
		testcontainers.WithLogConsumers(logger{}),
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("8080/tcp"),
			wait.ForLog("serving timekeeper"),
		),
	)
	testcontainers.CleanupContainer(t, timekeeperC)
	require.NoError(t, err)

	endpoint, err := timekeeperC.PortEndpoint(ctx, "8080/tcp", "http")
	require.NoError(t, err)
	// http client
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client := http.Client{Jar: jar}

	t.Run("Landing Page (GUEST)", func(t *testing.T) {
		resp, err := client.Get(endpoint)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Login", func(t *testing.T) {
		resp, err := client.PostForm(endpoint+"/_/login", url.Values{
			"username": []string{"admin"},
			"password": []string{"password1!"},
		})
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		for i, cookie := range resp.Cookies() {
			t.Log(i, cookie.Name, cookie.Value)
		}
	})

	t.Run("Landing Page (ADMIN)", func(t *testing.T) {
		resp, err := client.Get(endpoint)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// assert landing page shows admin actions
		require.Contains(t, string(body), `<a class="button" href="/event/new">New Event</a>`)
	})

	t.Run("Create Event (Admin)", func(t *testing.T) {
		resp, err := client.PostForm(endpoint+"/_/event", url.Values{
			"name":  []string{"Test Event 1"},
			"start": []string{"02.01.2026"},
			"slug":  []string{"test-event-1"},
		})
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

	})

	t.Run("Create Event (Guest)", func(t *testing.T) {
		resp, err := http.PostForm(endpoint+"/_/event", url.Values{})
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Landing Page with Events (Admin)", func(t *testing.T) {
		resp, err := client.Get(endpoint)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// assert landing page shows admin actions
		require.Contains(t, string(body), `Test Event 1`)
	})

	t.Run("Event Page (Guest)", func(t *testing.T) {
		resp, err := http.Get(endpoint + "/event/1")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// assert landing page shows admin actions
		require.Contains(t, string(body), `Test Event 1`)
	})

	t.Run("Event Page (Admin)", func(t *testing.T) {
		resp, err := client.Get(endpoint + "/event/1")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// assert landing page shows admin actions
		require.Contains(t, string(body), `Test Event 1`)
	})

	time.Sleep(1 * time.Second)
}

func getEnvOr(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

type logger struct {
}

func (l logger) Accept(log testcontainers.Log) {
	fmt.Print(string(log.Content))
}
