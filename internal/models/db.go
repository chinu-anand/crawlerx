package models

import (
	"fmt"
	"log"

	"github.com/chinu-anand/crawlerx/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := configs.App.DatabaseURL

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&CrawlJob{})
	if err != nil {
		log.Fatal("❌ Failed to auto migrate:", err)
	}

	fmt.Println("✅ Connected to Database")
}
