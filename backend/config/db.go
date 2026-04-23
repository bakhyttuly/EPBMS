package config

import (
	"epbms/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=akzharkyn dbname=epbms port=5432 sslmode=disable TimeZone=Asia/Almaty"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	DB = db
	fmt.Println("Database connected successfully")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Performer{},
		&models.Booking{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database migrated successfully")
}
