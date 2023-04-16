package app

import (
	"context"
	"encoding/json"
	"ethereum-scrapper/internal/config"
	"ethereum-scrapper/internal/services/eth"
	"ethereum-scrapper/internal/utils/graph"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func Start() {
	ctx := context.Background()

	cfg, err := config.ParseConfig("./config/config.yml")
	if err != nil {
		log.Fatalf("parse config error: %v", err)
	}

	ethClient, err := ethclient.Dial(cfg.EthereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to make a connection to the Ethereum Node: %v", err)
	}

	ethSvc := eth.New(ethClient)

	in, out, err := ethSvc.GetTransactionsByAddress(ctx, common.HexToAddress(cfg.Address), cfg.MinBlock, cfg.MaxBlock)
	if err != nil {
		log.Fatalf("GetTransactionsByAddress error: %v", err)
	}

	fj, err := getFormatedJSON(in)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("incoming txs summary:\n%v\n", string(fj))

	fj, err = getFormatedJSON(out)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("outgoing txs summary:\n%v\n", string(fj))

	err = graph.SaveScheme("scheme.gv", cfg.Address, in, out)
	if err != nil {
		log.Fatalf("SaveScheme error: %v", err)
	}
}

func getFormatedJSON(i interface{}) ([]byte, error) {
	f, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return f, nil
}
