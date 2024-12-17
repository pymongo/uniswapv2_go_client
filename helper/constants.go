package helper

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// https://docs.chain.link/data-feeds/l2-sequencer-feeds
// https://docs.arbitrum.io/run-arbitrum-node/sequencer/read-sequencer-feed
type EvmChainInfo struct {
	// ChainName string
	// GasCoin string
	ChainId      int64
	BlockTime    time.Duration // https://megaeth.systems/research
	RpcUrl       string
	WsUrl        string
	SequencerUrl string // Sequencer endpoints only support eth_sendRawTransaction
	Weth         common.Address
	Usdc         common.Address
	// Usdt common.Address
	V2router common.Address
	V2pair   common.Address
	V3pool   common.Address
}

var ChainInfo map[string]EvmChainInfo = map[string]EvmChainInfo{
	"BASE": {
		ChainId:   8453,
		BlockTime: 2000 * time.Millisecond,
		// RpcUrl: "https://mainnet.base.org",
		RpcUrl: "https://rpc.ankr.com/base/" + ankr,
		// WsUrl: "wss://base-rpc.publicnode.com",
		WsUrl:        "wss://rpc.ankr.com/base/ws/" + ankr,
		SequencerUrl: "https://mainnet.base.org",
		Weth:         common.HexToAddress("0x4200000000000000000000000000000000000006"),
		Usdc:         common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
		V2router:     common.HexToAddress("0x4752ba5dbc23f44d87826276bf6fd6b1c372ad24"),
		V2pair:       common.HexToAddress("0x88A43bbDF9D098eEC7bCEda4e2494615dfD9bB9C"),
		V3pool:       common.HexToAddress("0xd0b53D9277642d899DF5C87A3966A349A798F224"),
	},
	// https://docs.arbitrum.io/build-decentralized-apps/reference/node-providers
	"ARB": {
		ChainId:   42161,
		BlockTime: 250 * time.Millisecond,
		// RpcUrl: "https://arb1.arbitrum.io/rpc",
		RpcUrl:       "https://rpc.ankr.com/arbitrum/" + ankr,
		SequencerUrl: "https://arb1-sequencer.arbitrum.io/rpc",
		// WsUrl: "wss://arbitrum-one-rpc.publicnode.com",
		WsUrl:  "wss://rpc.ankr.com/arbitrum/ws/" + ankr,
		Weth:   common.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
		Usdc:   common.HexToAddress("0xaf88d065e77c8cC2239327C5EDb3A432268e5831"),
		V2pair: common.HexToAddress("0xF64Dfe17C8b87F012FCf50FbDA1D62bfA148366a"),
		V3pool: common.HexToAddress("0xC6962004f452bE9203591991D15f6b388e09E8D0"),
	},
	"OP":     {},
	"MANTLE": {},
}
