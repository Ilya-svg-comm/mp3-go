package store

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "host=192.168.145.66 dbname=player sslmode=disable password=postgres user=postgres"
	}

	os.Exit(m.Run())
}
