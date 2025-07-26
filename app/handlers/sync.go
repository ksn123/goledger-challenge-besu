package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	client, _ := ethclient.Dial(os.Getenv("BESU_RPC_URL"))
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Println("❌ Failed to open DB connection:", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	abiJson, _ := ioutil.ReadFile("../besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json")
	var abiMap map[string]interface{}
	_ = json.Unmarshal(abiJson, &abiMap)
	abiBytes, _ := json.Marshal(abiMap["abi"])
	parsedAbi, _ := abi.JSON(strings.NewReader(string(abiBytes)))
	data, _ := parsedAbi.Pack("get")
	addr := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	result, _ := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}, nil)

	var value *big.Int
	var parserErr = parsedAbi.UnpackIntoInterface(&value, "get", result)
	if parserErr != nil {
		log.Println("UnpackIntoInterface error:", parserErr)
		http.Error(w, "ABI unpack failed", 500)
		return
	}

	res := value.String()

	_, insertErr := db.Exec("INSERT INTO contract_state(value, synced_at) VALUES($1, $2)", res, time.Now())
	if insertErr != nil {
		log.Println(" DB insert failed:", insertErr)
		http.Error(w, "Database insert failed", 500)
		return
	}
	fmt.Fprintf(w, "Synced value %s to DB", res)
}
