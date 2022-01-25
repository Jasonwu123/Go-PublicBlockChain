package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// step5: 创建区块链
type BlockChain struct {
	// 最后一个区块的Hash值
	Tip []byte

	// 数据库对象
	DB *bolt.DB
}

//step6：创建区块链，带有创世块
/*
修改该方法：
    仅仅用来创建区块链
    如果数据库存在，证明区块链存在，直接结束该方法
    否则进程创建创世区块，并存入数据库中
*/
func CreateBlcokChainWithGensisBlcok(address string) {
	// a. 首先判断数据库是否存在，如果存在，则直接读取
	if dbExists() {
		fmt.Println("数据库已经存在。")
		return
	}

	fmt.Println("创建创世块：")

	// 2. 数据库不存在，说明第一次创建，然后存入到数据库中
	fmt.Println("数据库不存在。")

	// 创建创世块
	// 先创建coinbase
	txCoinBase := NewCoinBaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})

	// 打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 存入数据表
	err = db.Update(func(t *bolt.Tx) error {
		bk, err := t.CreateBucket([]byte(BLOCKTABLENAME))
		if err != nil {
			log.Panic(err)
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
		log.Panic(err)
	}
}

//step7：添加一个新的区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	// 1. 更新数据库
	err := bc.DB.Update(func(t *bolt.Tx) error {
		bk := t.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			// 2. 根据最新快的Hash读取数据，并反序列化最后一个块
			blockBytes := bk.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)

			// 3. 创建新的区块
			newBlock := NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)

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
		fmt.Println("\t交易：")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID: %x\n", tx.TxID)
			fmt.Println("\t\tVins: ")
			for _, in := range tx.Vins {
				fmt.Printf("\t\t\tTxID: %x\n", in.TxID)
				fmt.Printf("\t\t\tVout: %d\n", in.Vout)
				fmt.Printf("\t\t\tScriptSiq: %s\n", in.ScriptSiq)
			}
			fmt.Println("\t\tVouts: ")
			for _, out := range tx.Vouts {
				fmt.Printf("\t\t\tvalue: %d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey: %s\n", out.ScriptPubKey)
			}
		}
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

// 新增方法，用于获取区块链
func GetBlockchainObject() *BlockChain {
	/*
		1. 如果数据库不存在，直接返回nil
		2. 读取数据库
	*/

	if !dbExists() {
		fmt.Println("数据库不存在，无法获取区块链。")
		return nil
	}

	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockchain *BlockChain

	// 读取数据库
	err = db.View(func(tx *bolt.Tx) error {
		// 打开表
		bk := tx.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			// 读取最后一个Hash
			hash := bk.Get([]byte("1"))
			// 创建blockchain
			blockchain = &BlockChain{hash, db}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return blockchain
}

// 挖掘新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	/*
		1.新建交易
		2.新建区块
		3.将区块存入到数据库
	*/

	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		amountInt, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], amountInt, bc, txs)
		txs = append(txs, tx)
	}

	var (
		block    *Block // 数据库中的最后一个block
		newBlock *Block // 要创建的新block
	)
	bc.DB.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			hash := bk.Get([]byte("1"))
			blockBytes := bk.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	newBlock = NewBlock(txs, block.Hash, block.Height+1)

	bc.DB.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BLOCKTABLENAME))
		if bk != nil {
			bk.Put(newBlock.Hash, newBlock.Serialize())
			bk.Put([]byte("1"), newBlock.Hash)
			bc.Tip = newBlock.Hash
		}
		return nil
	})

}

// 找到所有未花费的交易输出
func (bc BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	/*
		1.先遍历未打包的交易(参数txs)，找出未花费的Output
		2.遍历数据库，获取每个块中的Transaction，找出未花费的Outpu
	*/
	var unUTXOs []*UTXO                      // 未花费
	spentTxOutputs := make(map[string][]int) // 存储已经花费

	// 添加先从txs遍历，查找未花费
	for i := len(txs) - 1; i >= 0; i-- {
		unUTXOs = caculate(txs[i], address, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		/*
			统计未花费
			获取block中的每个Transa
		*/
		for i := len(block.Txs) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.Txs[i], address, spentTxOutputs, unUTXOs)
		}

		// 结束迭代
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

func (bc BlockChain) GetBalance(address string, txs []*Transaction) int64 {
	unUTXOs := bc.UnUTXOs(address, txs)
	var amount int64
	for _, utxo := range unUTXOs {
		amount += utxo.Output.Value
	}
	return amount
}

// 转账时查获在可用的UTXO
func (bc BlockChain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
		1.获取所有的UTXO
		2.遍历UTXO
		返回值：map[hash]{index}
	*/
	var balance int64
	utxos := bc.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxID)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if balance >= amount {
			break
		}
	}
	if balance < amount {
		fmt.Printf("%s 余额不足。总额：%d, 需要：%d\n", from, balance, amount)
		os.Exit(1)
	}
	return balance, spendableUTXO
}

func caculate(tx *Transaction, address string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	// 先遍历TxInputs,表示花费
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			// 如果解锁
			if in.UnLockWithAddress(address) {
				key := hex.EncodeToString(in.TxID)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}
	// 遍历TxOutputs
outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(address) {
			// 如果对应的花费容器中长度不为0
			if len(spentTxOutputs) != 0 {
				var isSpentUTXO bool
				for txID, indexArry := range spentTxOutputs {
					for _, i := range indexArry {
						if i == index && txID == hex.EncodeToString(tx.TxID) {
							isSpentUTXO = true
							continue outputs
						}
					}
				}
				if !isSpentUTXO {
					utxo := &UTXO{tx.TxID, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			} else {
				utxo := &UTXO{tx.TxID, index, out}
				unUTXOs = append(unUTXOs, utxo)
			}
		}
	}
	return unUTXOs
}
