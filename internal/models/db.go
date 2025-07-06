package models

import (
	"log"
	"time"

	"github.com/chinu-anand/crawlerx/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := configs.App.DatabaseURL
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is empty")
	}

	var db *gorm.DB
	var err error

	for attempts := 1; attempts <= 5; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("⏳ Attempt %d: failed to connect to DB: %v", attempts, err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to Postgres")
	DB = db

	if err := DB.AutoMigrate(&CrawlJob{}); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
}
