package main

import (
	"context"
	"log"
	"math/big"
	"uniswapv2_go_client/bindings"
	"uniswapv2_go_client/helper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// https://suivision.xyz/object/0xb8d7d9e66a60c239e7a60110efcf8de6c705580ed924d0dde141f4a0e2c90105
// cetus USDC/SUI是2**64而uniswapV3是2**96
// 1 / ((270886523197904091236/2**64)**2 * 10**(6-9)) = 4.63
//
// https://basescan.org/address/0xd0b53D9277642d899DF5C87A3966A349A798F224#readContract
// (5000604271374031598349728/2**96)**2 * 10**(18-6) = 3983.69
func main() {
	info := helper.ChainInfo["BASE"]
	client, err := ethclient.Dial(info.RpcUrl)
	noerr(err)
	pooladdr := info.V3pool
	pool, err := bindings.NewV3Pool(pooladdr, client)
	noerr(err)
	slot0, err := pool.Slot0(nil)
	noerr(err)
	storage, err := client.StorageAt(context.Background(), pooladdr, common.BigToHash(big.NewInt(0)), nil)
	noerr(err)
	sqrtPriceX96_0 := new(big.Int).SetBytes(storage[20:])
	if sqrtPriceX96_0.Cmp(slot0.SqrtPriceX96) != 0 {
		log.Fatalln(sqrtPriceX96_0, slot0.SqrtPriceX96)
	}

	log.Println(helper.EthUsdcPrice(slot0.SqrtPriceX96))
}

func noerr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
