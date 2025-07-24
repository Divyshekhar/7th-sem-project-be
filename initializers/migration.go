package initializers

import (
	"fmt"
	"log"

	"github.com/Divyshekhar/7th-sem-project-be/models"
)

func init() {
	LoadEnv()
	ConnectDb()
}
func Migrate() {
	err := Db.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.UserSubject{},
		&models.Question{},
	)
	if err != nil {
		log.Fatal("Error migrating the models", err)
		return
	}
	fmt.Println("Database migrated successfully")
}
