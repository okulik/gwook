package settings_test

import (
	"os"
	"testing"

	"github.com/okulik/gwook/internal/settings"
)

func TestSettingsLoad(t *testing.T) {
	os.Setenv("SVIX_AUTH_TOKEN", "tok-123")
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("HTTP_SERVER_PORT", "8080")

	settings, err := settings.Load()
	if err != nil {
		t.Error("unable to load settings")
	}

	if settings.Svix.AuthToken != "tok-123" {
		t.Error("expected value for SVIX_AUTH_TOKEN does not exist")
	}
	if settings.Auth.Username != "admin" {
		t.Error("expected value for AUTH_USERNAME does not exist")
	}
	if settings.Http.Port != 8080 {
		t.Error("expected value for HTTP_SERVER_PORT does not exist")
	}
}
