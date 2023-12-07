package settings

import (
	"os"
	"path"

	"strings"

	"github.com/joho/godotenv"
)

const (
	Development = "development"
	Test        = "test"
)

func LoadEnvFile(environment string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	switch {
	case environment == Development:
		loadEnvFileFromFolder(currentDir, ".env.dev")
	case environment == Test:
		loadEnvFileFromFolder(currentDir, ".env.test")
	}

	return nil
}

func loadEnvFileFromFolder(dir string, envFile string) {
	err := loadEnv(path.Join(dir, envFile))
	if err != nil {
		runes := []rune(dir)
		lastSlash := strings.LastIndex(dir, "/")
		if lastSlash == 0 {
			return
		}

		newDir := string(runes[0:lastSlash])
		loadEnvFileFromFolder(newDir, envFile)
	}
}

func loadEnv(file string) error {
	localFile := file + ".local"

	if _, statErr := os.Stat(localFile); statErr == nil {
		if err := godotenv.Load(localFile); err != nil {
			return err
		}
	}

	return godotenv.Load(file)
}
