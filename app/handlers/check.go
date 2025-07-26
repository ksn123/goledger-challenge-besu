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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))

	var dbValue string
	err := db.QueryRow("SELECT value FROM contract_state ORDER BY synced_at DESC LIMIT 1").Scan(&dbValue)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("⚠️ No value found in database")
			http.Error(w, "No value in database", 404)
			return
		}
		log.Printf("❌ SQL error: %v\n", err)
		log.Println("❌ DB query error:", err)
		http.Error(w, "Database query failed", 500)
		return
	}

	client, _ := ethclient.Dial(os.Getenv("BESU_RPC_URL"))
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

	log.Println("check value", res)
	log.Println("check  dbvalue", dbValue)

	if dbValue == res {
		fmt.Fprint(w, `{"equal": true}`)
	} else {
		fmt.Fprint(w, `{"equal": false}`)
	}
}
