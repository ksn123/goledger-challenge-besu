package services

import (
	"fmt"

	"app/db"
)

func CompareContractWithDB() (bool, error) {
	chainValue, err := CallContract("get")
	if err != nil {
		return false, fmt.Errorf("failed to read from contract: %w", err)
	}

	dbValue, err := db.FetchLatestValue()
	if err != nil {
		return false, fmt.Errorf("failed to read from db: %w", err)
	}

	return dbValue == chainValue, nil
}
