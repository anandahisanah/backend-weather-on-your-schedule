package database

import (
	"backend-see-weather-on-your-schedule/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

var ginMode = os.Getenv("GIN_MODE")

func StartDB() {
	if ginMode != "release" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln(err.Error())
		}
	}

	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, dbPort)
	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	fmt.Println("success connecting to database")
	db.Debug().AutoMigrate(models.User{}, models.Province{}, models.City{}, models.Forecast{}, models.Event{})
}

func GetDB() *gorm.DB {
	return db
}
