package BLC

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// step5: 创建区块链
type BlockChain struct {
	// 最后一个区块的Hash值
	Tip []byte

	// 数据库对象
	DB *bolt.DB
}

//step6：创建区块链，带有创世块
func CreateBlcokChainWithGensisBlcok(data string) *BlockChain {
	// a. 首先判断数据库是否存在，如果存在，则直接读取
	if dbExists() {
		fmt.Println("数据库已经存在。")
		// 打开数据库
		db, err := bolt.Open(DBNAME, 0600, nil)
		if err != nil {
			log.Panic("数据库打开失败")
		}

		var blockchain *BlockChain

		// b. 读取数据库
		err = db.View(func(t *bolt.Tx) error {
			// 打开表
			bk := t.Bucket([]byte(BLOCKTABLENAME))
			if bk != nil {
				// 读取最后一个Hash
				hash := bk.Get([]byte("1"))
				blockchain = &BlockChain{hash, db}
			}
			return nil
		})

		if err != nil {
			log.Panic("数据库查询失败")
		}
		return blockchain
	}

	// 2. 数据库不存在，说明第一次创建，然后存入到数据库中
	fmt.Println("数据库不存在。")

	// 创建创世块
	genesisBlock := CreateGenesisBlock(data)

	// 打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Panic("数据库打开失败。")
	}

	// 存入数据表
	err = db.Update(func(t *bolt.Tx) error {
		bk, err := t.CreateBucket([]byte(BLOCKTABLENAME))
		if err != nil {
			log.Panic("数据表打开失败。")
		}
		if bk != nil {
			err = bk.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic("创世块存储有误。")
			}
			// 更新最新区块的Hash
			bk.Put([]byte("1"), genesisBlock.Hash)
		}
		return nil
	})

	if err != nil {
		log.Panic("数据存入失败。")
	}

	// 返回区块链对象
	return &BlockChain{genesisBlock.Hash, db}
}

//step7：添加一个新的区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(data string) {
	// 1. 更新数据库
	err := bc.DB.Update(func(t *bolt.Tx) error {
		bk := t.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			// 2. 根据最新快的Hash读取数据，并反序列化最后一个块
			blockBytes := bk.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)

			// 3. 创建新的区块
			newBlock := NewBlock(data, lastBlock.Hash, lastBlock.Height+1)

			// 4. 将新的区块序列化并存储
			err := bk.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 5. 更新最后一个Hash，以及blockchain的tip
			bk.Put([]byte("1"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic("区块添加失败")
	}
}

// 判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取一个迭代器方法
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

// 历数据库，打印输出所有的区块信息
func (bc *BlockChain) PrintChains() {
	// 获取迭代器对象
	bcIterator := bc.Iterator()

	var count = 0

	// 循环迭代
	for {
		block := bcIterator.Next()
		count++
		fmt.Printf("第%d个区块的信息：\n", count)
		//获取当前hash对应的数据，并进行反序列化
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("\t当前的hash：%x\n", block.Hash)
		fmt.Printf("\t数据：%s\n", block.Data)
		//fmt.Printf("\t时间：%v\n", block.TimeStamp)
		fmt.Printf("\t时间：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\t次数：%d\n", block.Nonce)

		// 值到创世块Hash值为0
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
}
