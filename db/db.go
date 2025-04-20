package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(autoMigrate bool) {
	// Set up the database connection
	// Use environment variable for DSN
	// Example: POSTGRES_DSN='host=localhost user=postgres password=password dbname=authservice port=5432 sslmode=disable TimeZone=America/Los_Angeles'
	postgresDSN := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
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
