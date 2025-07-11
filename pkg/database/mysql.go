package database

import (
	"database/sql"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// InitDB initializes and returns a MySQL database connection
func InitDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to MySQL!")
	return db
}
