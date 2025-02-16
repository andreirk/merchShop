package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username                 string            `gorm:"unique;not null"`
	Password                 string            `gorm:"not null"`
	CoinBalance              int               `gorm:"default:1000"`
	SentCoinTransactions     []CoinTransaction `gorm:"foreignKey:SenderID"`
	ReceivedCoinTransactions []CoinTransaction `gorm:"foreignKey:ReceiverID"`
	Orders                   []Order           `gorm:"foreignKey:UserID"`
}

type MerchItem struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Price int    `gorm:"not null"`
}

type CoinTransaction struct {
	gorm.Model
	SenderID   uint `gorm:"not null"`
	ReceiverID uint `gorm:"not null"`
	Amount     int  `gorm:"not null"`
	CreatedAt  time.Time
	Sender     User `gorm:"foreignKey:SenderID"`
	Receiver   User `gorm:"foreignKey:ReceiverID"`
}

type Order struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	ItemID    uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	CreatedAt time.Time
	User      User      `gorm:"foreignKey:UserID"`
	Item      MerchItem `gorm:"foreignKey:ItemID"`
}
