package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecContract(method string, setValue int64) (string, error) {

	abi, err := ParseABI(os.Getenv("ABI_JSON"))
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout())
	defer cancel()

	client, err := DialContext(ctx, os.Getenv("BESU_RPC_URL"))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Besu RPC: %w", err)
	}
	defer client.Close()

	slog.Info("Querying chain ID...")
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	boundContract := bind.NewBoundContract(contractAddress, abi, client, client, client)

	priv, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainId)
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %w", err)
	}

	tx, err := boundContract.Transact(auth, method, big.NewInt(setValue))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	slog.Info("Waiting for transaction to be mined", "tx", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return "", fmt.Errorf("transaction mining failed: %w", err)
	}

	slog.Info("Transaction mined", "receipt", receipt)
	return receipt.TxHash.Hex(), nil
}

func CallContract(method string) (string, error) {
	abi, err := ParseABI(os.Getenv("ABI_JSON"))
	if err != nil {
		return "", fmt.Errorf("parseABI failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout())
	defer cancel()

	client, err := DialContext(ctx, os.Getenv("BESU_RPC_URL"))
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	var output []interface{}
	err = boundContract.Call(&caller, &output, method)
	if err != nil {
		return "", fmt.Errorf("contract call failed: %w", err)
	}

	if len(output) == 0 {
		return "", fmt.Errorf("no output returned from method %q", method)
	}

	value, ok := output[0].(*big.Int)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected *big.Int, got %T", output[0])
	}

	return value.String(), nil
}

func ParseABI(abiJsonPath string) (abi.ABI, error) {

	abiJson, err := os.ReadFile(abiJsonPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI file: %w", err)
	}

	var contractMeta map[string]interface{}
	if err := json.Unmarshal(abiJson, &contractMeta); err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	abiPart, err := json.Marshal(contractMeta["abi"])
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to extract ABI definition: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiPart)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI string: %w", err)
	}

	return parsedABI, nil

}

func DialContext(ctx context.Context, besuUrl string) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, besuUrl)
	if err != nil {
		log.Printf("error dialing node: %v", err)
	}
	return client, err

}

func ContextTimeout() time.Duration {
	return 10 * time.Second
}
