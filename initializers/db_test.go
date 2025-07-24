package intializers

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestConnection(t *testing.T) {
	godotenv.Load(".test.env")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	ConnectDb()

	if Db == nil {
		t.Fatal("Expected Database connection, got nil")
	}

}
