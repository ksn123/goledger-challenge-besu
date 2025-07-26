package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SetHandler(w http.ResponseWriter, r *http.Request) {
	client, err := ethclient.Dial(os.Getenv("BESU_RPC_URL"))
	if err != nil {
		http.Error(w, "Failed to connect to RPC", 500)
		return
	}

	abiJson, _ := ioutil.ReadFile("../besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json")
	var abiMap map[string]interface{}
	_ = json.Unmarshal(abiJson, &abiMap)
	abiBytes, _ := json.Marshal(abiMap["abi"])
	parsedAbi, _ := abi.JSON(strings.NewReader(string(abiBytes)))

	var payload struct {
		Value int64 `json:"value"`
	}
	_ = json.NewDecoder(r.Body).Decode(&payload)

	privateKey, _ := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	chainID, _ := client.ChainID(context.Background())

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	input, _ := parsedAbi.Pack("set", big.NewInt(payload.Value))
	tx := types.NewTransaction(nonce, common.HexToAddress(os.Getenv("CONTRACT_ADDRESS")), big.NewInt(0), auth.GasLimit, gasPrice, input)

	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		http.Error(w, "Failed to send tx", 500)
		return
	}

	fmt.Fprintf(w, "TX Hash: %s", signedTx.Hash().Hex())
}
