package database

import (
	"log"

	"task-tracker/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Task{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db
}