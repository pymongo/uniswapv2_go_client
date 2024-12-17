package helper

import "math/big"

var q96 = new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

// formula `reserve1/10**decimal1 / reserve0/10**decimal0` = `reserve1/reserve0*10**(decimal0-decimal1)`
var ethusdcDecimal = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18-6), nil))

func EthUsdcPrice(sqrtPriceX96 *big.Int) *big.Float {
	sqrtPrice := new(big.Float).SetInt(sqrtPriceX96)
	sqrtPrice = sqrtPrice.Quo(sqrtPrice, q96)
	price := sqrtPrice.Mul(sqrtPrice, sqrtPrice)
	return price.Mul(price, ethusdcDecimal)
}
