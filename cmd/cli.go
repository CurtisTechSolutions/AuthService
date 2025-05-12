package cmd

import (
	"flag"
	"log/slog"
	"os"

	"github.com/CTS/AuthService/db"
	"github.com/CTS/AuthService/internal"
	"github.com/CTS/AuthService/server"
	"gorm.io/gorm"
)

type Config struct {
	DevMode     bool
	UseSQLite   bool
	AutoMigrate bool
	DBURL       string
	ServerPort  int
	// Add more fields as needed
}

func ParseFlags() *Config {
	// Define command-line flags
	// Make sure the defaults for these flags assume were running in production.
	isDevMode := flag.Bool("dev", false, "Run in development mode")
	useSQLite := flag.Bool("sqlite", false, "Use SQLite instead of Postgres. DBUrl flag will be ignored")
	autoMigrate := flag.Bool("automigrate", false, "Auto migrate the database")
	dbUrl := flag.String("db", "", "Database URL. Example: postgres://user:password@localhost:5432/dbname")
	serverPort := flag.Int("port", 9090, "Server port")
	// Add more flags as needed
	flag.Parse()

	return &Config{
		DevMode:     *isDevMode,
		UseSQLite:   *useSQLite,
		AutoMigrate: *autoMigrate,
		DBURL:       *dbUrl,
		ServerPort:  *serverPort,
	}
}

func Start() {
	cfg := ParseFlags()
	if cfg.DevMode {
		// Start in development mode
		slog.Info("Starting server in development mode")
		// Initialize the logger
		internal.InitializeLogger(&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		// Start in production mode
		slog.Info("Starting server in production mode")
		// Initialize the logger
		internal.InitializeLogger(&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	}

	// Initialize DB connection
	// Use SQLite if the flag is set to true
	if cfg.UseSQLite {
		// Start in SQLite mode
		slog.Info("Using SQLite database")
		// Initialize the database connection
		db.Connect(db.DialectorSQLite(), &gorm.Config{}, cfg.AutoMigrate)
	} else {
		// Use postgres if the DBURL is provided and SQLite is false
		// Check if the DBURL is provided
		if cfg.DBURL == "" {
			slog.Error("DBURL is required to use Postgres")
			os.Exit(1)
			return
		}
		// Start in Postgres mode
		slog.Info("Initializing Postgres connection")
		// Initialize the database connection
		db.Connect(db.DialectorPostgres(cfg.DBURL), &gorm.Config{}, cfg.AutoMigrate)
	}

	// Start the server
	err := server.Start(cfg.ServerPort)
	if err != nil {
		slog.Error("Server stopped", "error", err.Error())
		os.Exit(1)
	}
}
