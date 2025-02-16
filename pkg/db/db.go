package db

import (
	"fmt"
	"go/avito-test/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *config.Config) (*Db, func()) {
	gormInstance, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbClose := func() {
		db, _ := gormInstance.DB()
		if err := db.Close(); err != nil {
			fmt.Printf("Error closing database: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
	return &Db{gormInstance}, dbClose
}
