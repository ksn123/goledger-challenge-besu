package services

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecContract() *error {

	abi, err := parseABI(os.Getenv("ABI_JSON"))

	if err != nil {
		return &err
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout())
	defer cancel()

	client, err := dialContext(ctx, os.Getenv("BESU_RPC_URL"))
	if err != nil {
		return &err
	}
	defer client.Close()

	slog.Info("querying chain id")

	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("error querying chain id: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	priv, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("error loading private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainId)
	if err != nil {
		log.Fatalf("error creating transactor: %v", err)
	}

	tx, err := boundContract.Transact(auth, "get")
	if err != nil {
		log.Fatalf("error transacting: %v", err)
	}

	fmt.Println("waiting until transaction is mined",
		"tx", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(
		context.Background(),
		client,
		tx,
	)
	if err != nil {
		log.Fatalf("error waiting for transaction to be mined: %v", err)
	}

	fmt.Printf("transaction mined: %v\n", receipt)

	return nil
}

func CallContract() *error {
	var result interface{}

	abi, err := parseABI(os.Getenv("ABI_JSON"))

	if err != nil {
		return &err
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout())
	defer cancel()

	client, err := dialContext(ctx, os.Getenv("BESU_RPC_URL"))
	if err != nil {
		return &err
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
	err = boundContract.Call(&caller, &output, "get")
	if err != nil {
		log.Fatalf("error calling contract: %v", err)
	}
	result = output

	fmt.Println("Successfully called contract!", result)
	return nil
}

func parseABI(abiJson string) (abi.ABI, error) {
	abi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		log.Printf("error parsing abi: %v", err)

	}
	return abi, err

}

func dialContext(ctx context.Context, besuUrl string) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, besuUrl)
	if err != nil {
		log.Printf("error dialing node: %v", err)
	}
	return client, err

}

func contextTimeout() time.Duration {
	return 10 * time.Second
}
