package db

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(autoMigrate bool) {
	dsn := "host=localhost user=postgres password=password dbname=authservice port=5432 sslmode=disable TimeZone=America/Los_Angeles"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	slog.Info("Connected to database")

	if autoMigrate {
		slog.Info("Auto-migrating database")
		// Migrate the schema
		migrate()
	}
}

func migrate() {
	DB.AutoMigrate(&User{})
	slog.Info("Migrated User schema")

	DB.AutoMigrate(&Session{})
	slog.Info("Migrated Session schema")

	// Notify that migrations are complete
	slog.Info("Migrations completed")
}
