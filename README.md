```
abigen --abi bindings/router.abi --out bindings/router.go --pkg bindings --type Router
abigen --abi bindings/pair.abi --out bindings/pair.go --pkg bindings --type Pair
abigen --abi bindings/erc20.abi --out bindings/erc20.go --pkg bindings --type Erc20
abigen --abi bindings/UniswapV3Pool.abi --out bindings/UniswapV3Pool.go --pkg bindings --type V3Pool
```

> go run main.go
