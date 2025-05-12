package db

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DialectorPostgres(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func DialectorSQLite() gorm.Dialector {
	return sqlite.Open("test.db")
}

func Connect(dialector gorm.Dialector, config *gorm.Config, autoMigrate bool) {
	// Set up the database connection
	// Use environment variable for DSN

	db, err := gorm.Open(dialector, config)
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
	DB.AutoMigrate(&Session{})
	slog.Info("Migrated Session schema")

	DB.AutoMigrate(&User{})
	slog.Info("Migrated User schema")

	// Notify that migrations are complete
	slog.Info("Migrations completed")
}
