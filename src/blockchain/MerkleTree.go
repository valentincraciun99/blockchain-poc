package blockchain

import "crypto/sha256"

type MerkleTree struct {
	root *MerkleNode
}

type MerkleNode struct {
	Data  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		childrenHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(childrenHashes)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = right

	return &node
}

func CreateMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, item := range data {
		node := CreateMerkleNode(nil, nil, item)
		nodes = append(nodes, *node)
	}
	treeHeight := len(data) / 2

	for index := 0; index < treeHeight; index++ {
		var nextLevelNodes []MerkleNode

		for nodeIndex := 0; nodeIndex < len(nodes); nodeIndex = nodeIndex + 2 {
			node := CreateMerkleNode(&nodes[nodeIndex], &nodes[nodeIndex+1], nil)

			nextLevelNodes = append(nextLevelNodes, *node)

		}
		nodes = nextLevelNodes
	}

	treeRoot := MerkleTree{&nodes[0]}

	return &treeRoot
}
