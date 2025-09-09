package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type URL struct {
	LongURL   string `gorm:"primaryKey"`
	ShortCode string
	CreatedAt time.Time
}

var DB *gorm.DB

func Connect(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres : %v", err)
	}

	if err := DB.AutoMigrate(&URL{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
}

func GetLongURL(shortCode string) (string, error) {
	var url URL
	if err := DB.First(&url, "short_code = ?", shortCode).Error; err != nil {
		return "", err
	}
	return url.LongURL, nil
}

func InsertLongURL(longURL string, shortCode string) error {
	url := URL{
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}
	return DB.Create(&url).Error
}
