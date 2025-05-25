package main

import (
	"context"
	"log"
	"time"

	"dumb-api/config"
	"dumb-api/internal/services"

	"github.com/gobuffalo/pop/v6"
)

func main() {
	avalanche_handler := services.NewAvalancheHandler(config.AVALANCHE_RPC_URL)
	log.Printf("Listening for events")
	db, err := pop.Connect(config.ENV)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	go blockListener(avalanche_handler, db)

	select {}
}

func blockListener(handler services.EVMHandler, db *pop.Connection) {
	lastBlock, _ := handler.LastBlock(context.Background(), db)

	if lastBlock == 0 {
		block, err := handler.Block(context.Background(), nil)
		if err != nil {
			log.Printf("Failed to get latest block - if this isn't the first deploy check the error: %v", err)
			return
		}
		lastBlock = int64(block.NumberU64())
	}

	for {
		log.Printf("Listening for events")
		latestBlock, err := handler.Block(context.Background(), nil)
		if err != nil {
			log.Printf("Failed to get latest block: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		latestBlockNumber := int64(latestBlock.NumberU64())

		if latestBlockNumber <= lastBlock {
			time.Sleep(1 * time.Second)
			continue
		}

		nextBlockNumber := lastBlock + 1
		block, err := handler.BlockResult(context.Background(), &nextBlockNumber)
		if err != nil {
			log.Printf("Failed to get block %d: %v", nextBlockNumber, err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = handler.Listen(db, block)
		if err != nil {
			log.Printf("Failed to listen for %s events in block %d: %v", handler.GetChainName(), nextBlockNumber, err)
		}

		err = handler.UpdateLastBlock(db, nextBlockNumber)
		if err != nil {
			log.Printf("Failed to update last block to %d: %v", nextBlockNumber, err)
		} else {
			lastBlock = nextBlockNumber
			log.Printf("Processed block %d", nextBlockNumber)
		}

		blocksBehind := latestBlockNumber - lastBlock
		if blocksBehind > 10 {
			time.Sleep(1 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
