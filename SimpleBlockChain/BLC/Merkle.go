package BLC

import (
	"crypto/sha256"
	"math"
)

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNOde *MerkleNode
	Data      []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

// 给一个左右节点,生成一个新的节点
func NewMerkleNode(leftNode, rightNode *MerkleNode, txHash []byte) *MerkleNode {
	mNode := &MerkleNode{}
	if leftNode == nil && rightNode == nil {
		hash := sha256.Sum256(txHash)
		mNode.Data = hash[:]
	} else {
		prevHash := append(leftNode.Data, rightNode.Data...)
		hash := sha256.Sum256(prevHash)
		mNode.Data = hash[:]
	}
	mNode.LeftNode = leftNode
	mNode.RightNOde = rightNode
	return mNode
}

func NewMerkleTree(txHashData [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	// 判断交易的奇偶性
	if len(txHashData)%2 != 0 {
		// 奇数,复制最后一个
		txHashData = append(txHashData, txHashData[len(txHashData)-1])
	}
	// 创建一排的叶子节点
	for _, datum := range txHashData {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}
	// 生成树其他的节点
	count := GetCircleCount(len(nodes))

	for i := 0; i < count; i++ {
		var newLevel []*MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}
		// 判断newLevel的长度的奇偶性
		if len(newLevel)%2 != 0 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}
		nodes = newLevel
	}
	mTree := &MerkleTree{nodes[0]}
	return mTree
}

// 获取产生Merkle树根需要的循环的次数
func GetCircleCount(len int) int {
	count := 0
	for {
		if int(math.Pow(2, float64(count))) >= len {
			return count
		}
		count++
	}

}
