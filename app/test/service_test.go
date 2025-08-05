//go:build unit

package test

import (
	"math/big"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"app/services"
)

func TestCallContract(t *testing.T) {
	_ = godotenv.Load("../.env")
	res, err := services.CallContract("get")
	assert.NoError(t, err)
	_, ok := new(big.Int).SetString(res, 10)
	assert.True(t, ok)
}

func TestExecContract(t *testing.T) {
	_ = godotenv.Load("../.env")
	res, err := services.ExecContract("set", 999)
	assert.NoError(t, err)
	assert.Regexp(t, "^0x[0-9a-fA-F]{64}$", res)
}

func TestSyncValue(t *testing.T) {
	_ = godotenv.Load("../.env")
	_, err := services.SyncValue()
	assert.NoError(t, err)
}

func TestCompareContractWithDB(t *testing.T) {
	_ = godotenv.Load("../.env")
	match, err := services.CompareContractWithDB()
	assert.NoError(t, err)
	assert.IsType(t, true, match)
}
