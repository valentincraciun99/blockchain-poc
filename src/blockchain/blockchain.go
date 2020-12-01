package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}

	merkleTree := CreateMerkleTree(txHashes)
	txHash := merkleTree.root.Data

	return txHash[:]
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
	}
	hash := block.HashTransactions()

	block.Hash = hash
	return block
}

func (chain *BlockChain) AddBlock(txs []*Transaction) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(txs, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{1})
}

func InitBlockChain(address string) *BlockChain {
	cbtx := CoinbaseTx(address, "")
	genesis := Genesis(cbtx)

	blockchain := BlockChain{[]*Block{genesis}}

	return &blockchain
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	for _, block := range chain.Blocks {

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.Id)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutIndex)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	spentOuts := chain.GetSpentOutputs()
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount && !Find(spentOuts[txID], outIdx) {

				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (chain *BlockChain) GetSpentOutputs() map[string][]int {
	spentOuts := make(map[string][]int)

	for _, block := range chain.Blocks {
		for _, tx := range block.Transactions {
			for _, in := range tx.Inputs {
				txId := hex.EncodeToString(in.ID)
				spentOuts[txId] = append(spentOuts[txId], in.OutIndex)
			}
		}
	}
	return spentOuts
}

func Find(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (c *BlockChain) Print(addresses []string) {

	fmt.Print("Blocks:\n")
	for i := len(c.Blocks) - 1; i >= 0; i-- {
		b := c.Blocks[i]
		fmt.Printf("Block #%d\n", len(c.Blocks)-i)
		fmt.Printf("Hash: %s\n", hex.EncodeToString(b.Hash))
		fmt.Printf("Prev hash: %s\n", hex.EncodeToString(b.PrevHash))
		fmt.Printf("Transactions:\n")
		for _, tx := range b.Transactions {
			fmt.Printf("\tTx hash: %s\n", hex.EncodeToString(tx.Id))
		}
		fmt.Print("\n")
	}

	fmt.Print("Balances:\n")
	for _, address := range addresses {
		balance, _ := c.FindSpendableOutputs(address, math.MaxInt64)
		fmt.Printf("Address %s has a balance of %d\n", address, balance)
	}
	fmt.Print("\n")
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
