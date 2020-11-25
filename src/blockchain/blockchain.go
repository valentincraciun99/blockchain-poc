package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type BlockCahin struct{
	Blocks []*Block
}

type Block struct{
	Hash []byte
	Data []byte
	PrevHash []byte
}

func (b *Block) CreateHash(){
	data:= bytes.Join([][]byte{b.Data, b.PrevHash},[]byte{})
	hash:= sha256.Sum256(data)
	b.Hash = hash[:]
}

func CreateBlock(data string,prevHash []byte) *Block{
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	block.CreateHash()

	return block
}

func (chain *BlockCahin) AddBlock(data string){
	prevBlock:=chain.Blocks[len(chain.Blocks)-1]
	newBlock:= CreateBlock(data,prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

func Genesis() *Block{
	return CreateBlock("Genesis",[]byte{})
}

func InitBlockChain() *BlockCahin{
	return &BlockCahin{[]*Block{Genesis()}}
}
