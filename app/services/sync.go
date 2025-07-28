package services

import (
	"fmt"

	"app/db"
)

func SyncValue() (string, error) {
	value, err := CallContract("get")
	if err != nil {
		return "", fmt.Errorf("contract read error: %w", err)
	}

	err = db.SaveSyncedValue(value)
	if err != nil {
		return "", fmt.Errorf("db insert error: %w", err)
	}

	return value, nil
}
