package settings_test

import (
	"os"
	"testing"

	"github.com/okulik/gwook/internal/settings"
)

func TestEnvLoadEnvFile(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("unable to get current working directory")
	}
	createEnvFile(dir, t)
	defer deleteEnvFile(dir, t)

	if err := settings.LoadEnvFile("test"); err != nil {
		t.Error("unable to load env file")
	}

	result := os.Getenv("FOO")
	if result != "BAR" {
		t.Error("expected value does not exist.")
	}
}

func createEnvFile(dir string, t *testing.T) {
	file := dir + "/.env.test"
	_, err := os.Stat(file)
	if err != nil {
		d1 := []byte("FOO=BAR\n")
		err = os.WriteFile(file, d1, 0644)
		if err != nil {
			t.Error("unable to write .env.test file")
		}
	}
}

func deleteEnvFile(dir string, t *testing.T) {
	file := dir + "/.env.test"
	_, err := os.Stat(file)
	if err == nil {
		err = os.Remove(file)
		if err != nil {
			t.Error("unable to delete .env.test file")
		}
	}
}
