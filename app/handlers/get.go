package handlers

import (
	"context"
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
)

func GetHandler(w http.ResponseWriter, r *http.Request) {

	client, err := ethclient.Dial(os.Getenv("BESU_RPC_URL"))
	if err != nil {
		http.Error(w, "RPC connection failed", 500)
		return
	}

	abiJson, err := ioutil.ReadFile("../besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json")
	if err != nil {
		http.Error(w, "Failed to read ABI", 500)
		return
	}

	var contractMeta map[string]interface{}
	_ = json.Unmarshal(abiJson, &contractMeta)
	abiString, err := json.Marshal(contractMeta["abi"])
	if err != nil {
		http.Error(w, "Failed to marshal ABI", 500)
		return
	}

	parsedAbi, err := abi.JSON(strings.NewReader(string(abiString)))
	if err != nil {
		http.Error(w, "Failed to parse ABI", 500)
		return
	}

	data, err := parsedAbi.Pack("get")
	if err != nil {
		http.Error(w, "Failed to encode ABI call", 500)
		return
	}

	toAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	callMsg := ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		http.Error(w, "CallContract failed", 500)
		return
	}

	var value *big.Int
	err = parsedAbi.UnpackIntoInterface(&value, "get", result)
	if err != nil {
		log.Println("Unpack error:", err)
		http.Error(w, "Failed to unpack result", 500)
		return
	}

	fmt.Fprintf(w, "Value on chain: %s", value.String())
}
