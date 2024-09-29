package main

import (
	"log"
	"math/big"
	"uniswapv2_go_client/bindings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	rpcUrl := "https://mainnet.base.org"
	weth := common.HexToAddress("0x4200000000000000000000000000000000000006")
	// usdc := common.HexToAddress("0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85") // OP
	usdc := common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913")
	routerAddr := common.HexToAddress("0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24") // https://docs.uniswap.org/contracts/v2/reference/smart-contracts/v2-deployments

	client, err := ethclient.Dial(rpcUrl)
	assertNoErr(err)
	routerClient, err := bindings.NewRouterCaller(routerAddr, client)
	assertNoErr(err)
	// factory, err := routerClient.Factory(nil)
	out, err := routerClient.GetAmountsOut(nil, big.NewInt(1e18), []common.Address{weth, usdc})
	assertNoErr(err)
	sellEthGetUsdc, _ := out[1].Float64() // out[0] 是支付的ETH数量
	log.Printf("1ETH=%fUSDC\n", sellEthGetUsdc/1e6)
}

func assertNoErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
