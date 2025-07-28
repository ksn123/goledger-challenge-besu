package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func SaveSyncedValue(value string) error {
	db, err := OpenConn()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO contract_state(value, synced_at) VALUES($1, $2)", value, time.Now())
	if err != nil {
		return fmt.Errorf("db exec failed: %w", err)
	}

	return nil
}

func FetchLatestValue() (string, error) {
	db, err := OpenConn()
	if err != nil {
		return "", err
	}
	defer db.Close()
	var value string
	err = db.QueryRow("SELECT value FROM contract_state ORDER BY synced_at DESC LIMIT 1").Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return value, nil
}

func OpenConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		return nil, fmt.Errorf("db open failed: %w", err)
	}

	return db, nil
}
