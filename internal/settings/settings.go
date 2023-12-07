package settings

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type HttpSettings struct {
	GracefulShutdownTimeout time.Duration `envconfig:"HTTP_SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"10s"`
	Port                    int           `envconfig:"HTTP_SERVER_PORT" default:"4000"`
	IdleTimeout             time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`
	ReadTimeout             time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout            time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"20s"`
}

type AuthSettings struct {
	Username string `envconfig:"AUTH_USERNAME" required:"true"`
	Password string `envconfig:"AUTH_PASSWORD" required:"true"`
	Realm    string `envconfig:"AUTH_REALM" default:"localhost"`
}

type SvixSettings struct {
	AuthToken     string `envconfig:"SVIX_AUTH_TOKEN"`
	Debug         bool   `envconfig:"SVIX_DEBUG" default:"false"`
	ServerUrl     string `envconfig:"SVIX_SERVER_URL"`
	RetryMax      int    `envconfig:"SVIX_RETRY_MAX" default:"10"`
	ApplicationID string `envconfig:"SVIX_APPLICATION_ID" default:"app_0000000000000000"`
}

type Settings struct {
	Http *HttpSettings
	Auth *AuthSettings
	Svix *SvixSettings
}

func Load() (*Settings, error) {
	appEnv := getAppEnv()
	if rawErr := LoadEnvFile(appEnv); rawErr != nil {
		return nil, fmt.Errorf("failed to load env file for '%s': %s", appEnv, rawErr)
	}

	settings := &Settings{}
	if err := envconfig.Process("", settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func getAppEnv() string {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return "development"
	}
	return appEnv
}
