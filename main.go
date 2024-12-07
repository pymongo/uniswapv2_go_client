package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"uniswapv2_go_client/bindings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Constants
const (
	rpcUrl     = "https://mainnet.base.org"
	wethAddr   = "0x4200000000000000000000000000000000000006"
	usdcAddr   = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	routerAddr = "0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24"
)

// PrefixLogger is a logger with a prefix
type PrefixLogger struct {
	*log.Logger
	prefix string
}

// Printf overwrites the Printf method, adding the prefix
func (l *PrefixLogger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(l.prefix+format, v...)
}

// NewPrefixLogger creates a new PrefixLogger
func NewPrefixLogger(prefix string) *PrefixLogger {
	return &PrefixLogger{
		Logger: log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile),
		prefix: prefix,
	}
}

func main() {
	client, err := ethclient.Dial(rpcUrl)
	assertNoErr(err)
	getPriceFromStorage(client)
	getPriceFromRouter(client)
	getPriceFromPair(client)
}

// forge inspect
// https://explorer.sim.io/base/23375405/0x88a43bbdf9d098eec7bceda4e2494615dfd9bb9c
func getPriceFromStorage(client *ethclient.Client) {
	pairAddress := common.HexToAddress("0x88A43bbDF9D098eEC7bCEda4e2494615dfD9bB9C")
	storage, err := client.StorageAt(context.Background(), pairAddress, common.BigToHash(big.NewInt(8)), nil) // Get storage at slot 8
	if err != nil {
		panic(err)
	}

	// Convert storage to big.Int for bitwise operations
	value := new(big.Int).SetBytes(storage)
	// Create bit masks (1n << 112n) - 1n
	mask112 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 112), big.NewInt(1))
	// Extract reserve0 (least significant 112 bits)
	reserve0 := new(big.Int).And(value, mask112)
	// Extract reserve1 (middle 112 bits)
	reserve1 := new(big.Int).And(
		new(big.Int).Rsh(value, 112),
		mask112,
	)

	// Extract blockTimestampLast (most significant 32 bits)
	// blockTimestampLast := uint32(new(big.Int).Rsh(value, 224).Uint64())

	// log.Println(reserve0, reserve1, blockTimestampLast)
	r0, _ := reserve0.Float64()
	r1, _ := reserve1.Float64()
	fmt.Println("using StorageAt API, ETH price", (r1/1e6)/(r0/1e18))
}

// Get price through Router contract
func getPriceFromRouter(client *ethclient.Client) {
	logger := NewPrefixLogger("Method 1 (Router): ")

	routerClient, err := bindings.NewRouterCaller(common.HexToAddress(routerAddr), client)
	assertNoErr(err)

	out, err := routerClient.GetAmountsOut(nil, big.NewInt(1e18), []common.Address{
		common.HexToAddress(wethAddr),
		common.HexToAddress(usdcAddr),
	})
	assertNoErr(err)

	sellEthGetUsdc, _ := out[1].Float64()
	price := sellEthGetUsdc / 1e6
	logger.Printf("1 ETH = %.6f USDC\n", price)
}

// Get price and other relevant information through Pair contract
func getPriceFromPair(client *ethclient.Client) {
	logger := NewPrefixLogger("Method 2 (Pair): ")

	// Get Factory contract
	routerClient, err := bindings.NewRouterCaller(common.HexToAddress(routerAddr), client)
	assertNoErr(err)
	factoryAddr, err := routerClient.Factory(nil)
	assertNoErr(err)

	// Get Pair address
	factoryClient, err := bindings.NewFactoryCaller(factoryAddr, client)
	assertNoErr(err)
	pairAddr, err := factoryClient.GetPair(nil, common.HexToAddress(wethAddr), common.HexToAddress(usdcAddr))
	assertNoErr(err)
	log.Println("pairAddr", pairAddr)

	// Get reserves from Pair contract
	pairClient, err := bindings.NewPairCaller(pairAddr, client)
	assertNoErr(err)
	reserves, err := pairClient.GetReserves(nil)
	assertNoErr(err)

	// Calculate price
	ethReserve := new(big.Int).Set(reserves.Reserve0)
	usdcReserve := new(big.Int).Set(reserves.Reserve1)

	// Calculate price without considering trading fees and slippage
	rawPrice := new(big.Float).Quo(
		new(big.Float).SetInt(usdcReserve),
		new(big.Float).SetInt(ethReserve),
	)
	rawPriceFloat, _ := rawPrice.Float64()
	rawPriceAdjusted := rawPriceFloat * 1e12 // Adjust decimal places

	logger.Printf("1 ETH = %.6f USDC (without trading fees and slippage)", rawPriceAdjusted)

	// Calculate price using logic similar to UniswapV2Library.getAmountOut
	amountIn := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1 ETH
	amountOut := getAmountOut(amountIn, ethReserve, usdcReserve)

	// Adjust decimal places
	amountOutFloat := new(big.Float).SetInt(amountOut)
	amountOutAdjusted := new(big.Float).Quo(amountOutFloat, big.NewFloat(1e6))

	priceFloat, _ := amountOutAdjusted.Float64()

	logger.Printf("1 ETH = %.6f USDC (with trading fees and slippage)", priceFloat)

	// Calculate estimated slippage rate
	slippageRate := (rawPriceAdjusted - priceFloat) / rawPriceAdjusted * 100
	logger.Printf("Estimated slippage rate: %.2f%%", slippageRate)

	// Print other useful information
	ethReserveFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(ethReserve), big.NewFloat(1e18)).Float64()
	usdcReserveFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(usdcReserve), big.NewFloat(1e6)).Float64()
	logger.Printf("ETH reserve: %.6f", ethReserveFloat)
	logger.Printf("USDC reserve: %.6f", usdcReserveFloat)
	logger.Printf("Last update time: %d", reserves.BlockTimestampLast)

	// Get other information from Factory
	feeTo, err := factoryClient.FeeTo(nil)
	assertNoErr(err)
	feeToSetter, err := factoryClient.FeeToSetter(nil)
	assertNoErr(err)
	logger.Printf("FeeTo address: %s", feeTo.Hex())
	logger.Printf("FeeToSetter address: %s", feeToSetter.Hex())
}

// Error checking helper function
func assertNoErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Implement a function similar to UniswapV2Library.getAmountOut
func getAmountOut(amountIn, reserveIn, reserveOut *big.Int) *big.Int {
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(
		new(big.Int).Mul(reserveIn, big.NewInt(1000)),
		amountInWithFee,
	)
	return new(big.Int).Div(numerator, denominator)
}
