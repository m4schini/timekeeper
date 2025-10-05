package config

import (
	"os"
	"time"
)

func Timezone() *time.Location {
	l, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	return l
}

func DatabaseConnectionString() string {
	return mustEnv("DATABASE_CONNECTIONSTRING")
}

func HmacSecret() []byte {
	return []byte(mustEnv("JWT_SECRET"))
}

func AdminPassword() string {
	return mustEnv("ADMIN_PASSWORD")
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
