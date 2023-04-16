package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"sync"
)

const maxGoroutinesPerTask = 50

type Service struct {
	client *ethclient.Client
}

func New(client *ethclient.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetTransactionsByAddress(ctx context.Context, addr common.Address, minBlock, maxBlock int64) (
	in map[string]*big.Int,
	out map[string]*big.Int,
	err error,
) {
	in = make(map[string]*big.Int)
	inMx := &sync.Mutex{}

	out = make(map[string]*big.Int)
	outMx := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	activeGoroutines := make(chan struct{}, maxGoroutinesPerTask)

	for curBlock := minBlock; curBlock <= maxBlock; curBlock++ {
		log.Println("processing block: ", curBlock)

		wg.Add(1)
		activeGoroutines <- struct{}{}

		go func(curBlock int64) {
			defer func() {
				<-activeGoroutines
				wg.Done()
			}()

			block, err := s.client.BlockByNumber(ctx, big.NewInt(curBlock))
			if err != nil {
				log.Println("BlockByNumber error: ", err)
				return
			}

			for _, tx := range block.Transactions() {
				if tx.To() == nil {
					continue
				}

				from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
				if err != nil {
					log.Println("Sender func error: ", err)
					continue
				}

				if tx.To().Hex() == addr.Hex() {
					inMx.Lock()

					if in[from.Hex()] == nil {
						in[from.Hex()] = big.NewInt(0)
					}

					in[from.Hex()].Add(in[from.Hex()], tx.Value())

					inMx.Unlock()
				}

				if from.Hex() == addr.Hex() {
					outMx.Lock()

					if out[tx.To().Hex()] == nil {
						out[tx.To().Hex()] = big.NewInt(0)
					}

					out[tx.To().Hex()].Add(out[tx.To().Hex()], tx.Value())

					outMx.Unlock()
				}
			}
		}(curBlock)
	}

	wg.Wait()

	return in, out, nil
}
