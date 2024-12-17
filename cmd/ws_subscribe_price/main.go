package main

import (
	"log"
	"uniswapv2_go_client/bindings"
	"uniswapv2_go_client/helper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// https://github.com/pymongo/uniswap_client/blob/master/cmd/uniswap_get_price_demo/main.go
func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	info := helper.ChainInfo["ARB"]
	client, err := ethclient.Dial(info.WsUrl)
	noerr(err)

	pair, err := bindings.NewPair(info.V2pair, client)
	noerr(err)
	syncCh := make(chan *bindings.PairSync, 2)
	syncSub, err := pair.WatchSync(nil, syncCh) // client.SubscribeFilterLogs
	noerr(err)

	pool, err := bindings.NewV3Pool(info.V3pool, client)
	noerr(err)
	swapCh := make(chan *bindings.V3PoolSwap, 2)
	swapSub, err := pool.WatchSwap(nil, swapCh, []common.Address{}, []common.Address{})
	noerr(err)

	log.Println("before for select")
	for {
		select {
		case err = <-syncSub.Err():
			log.Fatalln(err)
		case err = <-swapSub.Err():
			log.Fatalln(err)
		case syncEvt := <-syncCh:
			reserve0_, _ := syncEvt.Reserve0.Float64()
			reserve0 := reserve0_ / 1e18
			reserve1_, _ := syncEvt.Reserve1.Float64()
			reserve1 := reserve1_ / 1e6
			log.Println(syncEvt.Raw.BlockNumber, reserve1/reserve0)
		case swapEvt := <-swapCh:
			log.Println(swapEvt.Raw.BlockNumber, helper.EthUsdcPrice(swapEvt.SqrtPriceX96))
		}
	}
}

func noerr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
