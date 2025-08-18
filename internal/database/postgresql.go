package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgreSQLDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),os.Getenv("DB_PASSWORD"),os.Getenv("DB_NAME"),os.Getenv("DB_HOST"),os.Getenv("DB_PORT"))
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Print("Databse connected successfuly\n")
	return db, nil
}
