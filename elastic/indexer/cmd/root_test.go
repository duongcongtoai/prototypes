package cmd

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnv(t *testing.T) {

	os.Setenv("DATABASE_HOSTNAME", "whatup")

	initConfig()

	if config.Database.Hostname != "whatup" {
		t.Fatalf("Host must be '5kbps.io', got %s", config.Database.Hostname)
	}
}
