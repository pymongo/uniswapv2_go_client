package main

import (
	"log"
	"math/big"
	"uniswapv2_go_client/bindings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var q96 = new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

// https://suivision.xyz/object/0xb8d7d9e66a60c239e7a60110efcf8de6c705580ed924d0dde141f4a0e2c90105
// cetus USDC/SUI是2**64而uniswapV3是2**96
// 1 / ((270886523197904091236/2**64)**2 * 10**(6-9)) = 4.63
//
// https://basescan.org/address/0xd0b53D9277642d899DF5C87A3966A349A798F224#readContract
// (5000604271374031598349728/2**96)**2 * 10**(18-6) = 3983.69
func main() {
	client, err := ethclient.Dial("https://mainnet.base.org")
	noerr(err)
	pooladdr := common.HexToAddress("0xd0b53D9277642d899DF5C87A3966A349A798F224")
	pool, err := bindings.NewV3Pool(pooladdr, client)
	noerr(err)
	slot0, err := pool.Slot0(nil)
	noerr(err)
	sqrtPrice := new(big.Float).SetInt(slot0.SqrtPriceX96)
	sqrtPrice = sqrtPrice.Quo(sqrtPrice, q96)
	price := sqrtPrice.Mul(sqrtPrice, sqrtPrice)
	decimal := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18-6), nil))
	price = price.Mul(price, decimal)
	log.Println(price)
}

func noerr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
