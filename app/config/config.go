package config

import (
	"os"
	"regexp"
	"strings"
	"time"
)

var PixelHackPlaceholderRx = regexp.MustCompile(`:([a-z0-9_]+):`)

func Timezone() *time.Location {
	timezone := getEnvOr("TIMEKEEPER_TIMEZONE", "Europe/Berlin")
	l, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return l
}

func TelemetryEnabled() bool {
	return strings.ToLower(getEnvOr("TIMEKEEPER_TELEMETRY_ENABLED", "false")) == "true"
}

func MetricsEndpointToken() string {
	return getEnvOr("TIMEKEEPER_METRICS_TOKEN", "")
}

func DatabaseConnectionString() string {
	return mustEnv("DATABASE_CONNECTIONSTRING")
}

func HmacSecret() []byte {
	return []byte(mustEnv("JWT_SECRET"))
}

func AdminPassword() string {
	return getEnvOr("TIMEKEEPER_ADMIN_PASSWORD", "")
}

func BaseUrl() string {
	return getEnvOr("TIMEKEEPER_BASE_URL", "https://zeit.haeck.se")
}

func Port() string {
	return getEnvOr("PORT", "80")
}

func getEnvOr(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

func mustEnv(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		panic(envName + " MUST BE SET")
	}

	return value
}
