package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	tools "github.com/sajimenezher_meli/meli-frescos-8/pkg"

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

func SelectOne(db *sql.DB, tablename string, fields []string, condition string, values ...any) *sql.Row {
	columns := tools.SliceToString(fields, ",")
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", columns, tablename)
	if condition != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, condition)
	}

	return db.QueryRow(sqlStatement, values...)
}

func Select(db *sql.DB, tablename string, fields []string, condition string, values ...any) (*sql.Rows, error) {
	columns := tools.SliceToString(fields, ",")
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", columns, tablename)
	if condition != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, condition)
	}

	return db.Query(sqlStatement, values...)
}

func Insert(db *sql.DB, tablename string, data map[any]any) (sql.Result, error) {
	keys, values := tools.GetSlicesOfKeyAndValuesFromMap(data)
	columns := tools.SliceToString(keys, ",")
	placeholders := tools.SliceToString(tools.FillNewSlice(len(data), "?"), ",")

	sqlStatement := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", tablename, columns, placeholders)

	return db.Exec(sqlStatement, values...)
}

func Delete(db *sql.DB, tablename string, condition string, values ...any) (sql.Result, error) {
	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE %s;", tablename, condition)

	return db.Exec(sqlStatement, values...)
}
