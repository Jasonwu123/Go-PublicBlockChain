package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// step1: 创建Blcok结构体
type Block struct {
	// 高度Height：区块编号，第一个区块叫创世块，高度为0
	Height int64

	// 前区块哈希值PrevHash
	PrevBlockHash []byte

	// 交易数据Data: 目前先设计为[]byte,后期是Transaction
	Data []byte

	// 时间戳TimeStamp
	TimeStamp int64

	// 当前区块哈希值Hash：32个字节，64个16进制
	Hash []byte

	Nonce int64
}

// step2: 创建新的区块
func NewBlock(data string, prevBlockHash []byte, height int64) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
	}

	// 调用方法设置哈希值
	//block.SetHash()

	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()

	block.Hash = hash

	block.Nonce = nonce

	return block
}

// step3: 设置区块的hash
func (block *Block) SetHash() {
	// 1. 将高度转为字节数组
	heightBytes := IntToHex(block.Height)

	// 2. 将时间戳转为字节数组
	timeString := strconv.FormatInt(block.TimeStamp, 2)
	timeBytes := []byte(timeString)
	//timeBytes := IntToHex(block.TimeStamp)

	// 3. 拼接所有的属性
	blockBytes := bytes.Join([][]byte{
		heightBytes,
		block.PrevBlockHash,
		block.Data,
		timeBytes,
	}, []byte{})

	// 4. 生成哈希值
	hash := sha256.Sum256(blockBytes)
	block.Hash = hash[:]
}

// step4: 创建创世块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, make([]byte, 32, 32), 0)
}
