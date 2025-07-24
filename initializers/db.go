package intializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var TestDb *gorm.DB

func ConnectDb() {

	var err error
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error Connecting to the database")
		return
	}
}
