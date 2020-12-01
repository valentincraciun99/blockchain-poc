package main

import (
	"github.com/valentincraciun99/blockchain-poc/blockchain"
)

const (
	Alice = "17de9d40ed1797efa51739e5b3db7ab97aa8aadeec71d3f1a3ba5f6760caf847"
	Bob   = "29123b23c65f86f9392fc71e206b1f020533d3322914c68325167d16c9940e5a"
)

func main() {
	Addresses := []string{Alice, Bob}
	chain := blockchain.InitBlockChain(Alice)

	tx := blockchain.NewTransaction(Alice, Bob, 10, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	chain.Print(Addresses)

	tx1 := blockchain.NewTransaction(Alice, Bob, 10, chain)
	chain.AddBlock([]*blockchain.Transaction{tx1})
	chain.Print(Addresses)

}
