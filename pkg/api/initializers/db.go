package initializers

import (
	"fmt"
	"log"

	"github.com/emarifer/go-fiber-jwt/pkg/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Madrid",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassoword,
		config.DBName,
		config.DBPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("🔥 failed to connect to the Database!\n", err.Error())
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("🔥 migration failed:\n", err.Error())
	}

	log.Println("🚀 Connected Successfully to the Database")
}
