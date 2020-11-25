package main

import (
	"github.com/valentincraciun99/blockchain-poc/blockchain"
)

func main() {
	chain:= blockchain.InitBlockChain()

	chain.AddBlock("First")

}
