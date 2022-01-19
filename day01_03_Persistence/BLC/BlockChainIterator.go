package BLC

import (
	"log"

	"github.com/boltdb/bolt"
)

// 新增一个结构体
type BlockChainIterator struct {
	// 当前区块的Hash
	CurrentHash []byte

	// 数据库对象
	DB *bolt.DB
}

// 获取区块
func (bcIterator *BlockChainIterator) Next() *Block {
	block := new(Block)

	// 1. 打开数据库并读取
	err := bcIterator.DB.View(func(t *bolt.Tx) error {
		// 2. 打开数据表
		bk := t.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			// 3. 根据当前Hash获取并反序列化
			blockBytes := bk.Get(bcIterator.CurrentHash)
			block = DeserializeBlock(blockBytes)

			// 4. 更新当前的Hash
			bcIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic("数据库打开失败")
	}

	return block
}
