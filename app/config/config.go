package config

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var PixelHackPlaceholderRx = regexp.MustCompile(`:([a-z0-9_]+):`)

func Timezone() *time.Location {
	timezone := viper.GetString("timezone")
	l, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return l
}

func TelemetryEnabled() bool {
	return viper.GetBool("telemetry.enabled")
}

func MetricsEndpointToken() string {
	return viper.GetString("telemetry.endpoint.token")
}

func DatabaseConnectionString() string {
	str := viper.GetString("database.connectionstring")
	if str == "" {
		panic("RZA_DATABASE_CONNECTIONSTRING is required")
	}
	return str
}

func HmacSecret() []byte {
	secret := viper.GetString("jwt.secret.key")
	if secret != "" {
		return []byte(secret)
	}
	secretFileLocation := viper.GetString("jwt.secret.file")
	s, err := os.ReadFile(secretFileLocation)
	if err != nil {
		panic(err)
	}
	return s
}

func AdminPassword() string {
	return viper.GetString("admin.password")
}

func BaseUrl() string {
	return viper.GetString("baseUrl")
}

func Port() string {
	return viper.GetString("port")
}

func Load() error {
	viper.SetDefault("timezone", "Europe/Berlin")
	viper.SetDefault("telemetry.enabled", false)
	viper.SetDefault("baseUrl", "https://zeit.haeck.se")
	viper.SetDefault("port", "80")
	viper.SetDefault("jwt.secret.file", "/etc/raumzeitalpaka/jwt.secret")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("rza")
	viper.AutomaticEnv()

	// Add search paths to find the file
	viper.SetConfigName("raumzeitalpaka")
	viper.AddConfigPath("/etc/raumzeitalpaka/")
	viper.AddConfigPath("$HOME/.raumzeitalpaka")
	viper.AddConfigPath(".")

	// Find and read the config file
	return viper.ReadInConfig()
}
