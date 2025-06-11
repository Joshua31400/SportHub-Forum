package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitDB(dsn string) error {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}

	DB = sqlDB
	log.Println("✅ Connexion à la base réussie.")
	return nil
}
