package postgres

import (
	"FIO_App/pkg/storage/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type Storage struct {
	DB *gorm.DB
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	log.Println("Connected!")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations..")
	err = db.AutoMigrate(&models.Person{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
