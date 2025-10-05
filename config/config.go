package config

import "os"

func DatabaseConnectionString() string {
	return mustEnv("DATABASE_CONNECTIONSTRING")
}

func HmacSecret() []byte {
	return []byte(mustEnv("JWT_SECRET"))
}

func AdminPassword() string {
	return mustEnv("ADMIN_PASSWORD")
}

func mustEnv(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		panic(envName + " MUST BE SET")
	}

	return value
}
