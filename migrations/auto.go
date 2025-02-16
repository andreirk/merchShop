package main

import (
	"go/avito-test/internal/models"
	"go/avito-test/pkg/db"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error connecting to database while auto migration", err)
	}
	err = database.AutoMigrate(&models.CoinTransaction{}, &models.User{}, &models.Order{}, &models.MerchItem{})
	if err != nil {
		log.Fatal(err)
	}

	seedMerch(&db.Db{DB: database})

	log.Println("Database migrated")
}

func seedMerch(db *db.Db) {
	merchItems := []models.MerchItem{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
		{Name: "powerbank", Price: 200},
		{Name: "hoody", Price: 300},
		{Name: "umbrella", Price: 200},
		{Name: "socks", Price: 10},
		{Name: "wallet", Price: 50},
		{Name: "pink-hoody", Price: 500},
	}

	for _, item := range merchItems {
		db.Create(&item)
	}

	log.Println("Merch items seeded successfully!")
}
