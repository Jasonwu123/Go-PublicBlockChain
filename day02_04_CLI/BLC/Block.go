package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
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

	// step3：调用工作量证明的方法，并返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()

	block.Hash = hash

	block.Nonce = nonce

	return block
}

// step4: 创建创世块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, make([]byte, 32, 32), 0)
}

// 将区块序列化，得到一个字节数组---区块的行为，设计为方法
func (block *Block) Serialize() []byte {
	// 1. 创建一个buffer
	var result bytes.Buffer

	// 2. 创建一个编码器
	encoder := gob.NewEncoder(&result)

	// 3. 编码---> 打包
	err := encoder.Encode(block)
	if err != nil {
		log.Panic("编码失败")
	}
	return result.Bytes()
}

// 反序列化，得到一个区块---设计为函数
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	reader := bytes.NewReader(blockBytes)

	// 1. 创建一个解码器
	decoder := gob.NewDecoder(reader)

	// 2. 解包
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码失败")
	}
	return &block
}
