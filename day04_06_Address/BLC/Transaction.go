package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"time"
)

// 创建Transaction结构体
type Transaction struct {
	TxID  []byte     // 交易ID
	Vins  []*TXInput // 输入
	Vouts []*TXOuput // 输出
}

/*
 Transaction 创建分两种情况
 1.创世块创建时的Transaction
 2.转账时产生的Transaction
*/
func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOuput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOuput{txOutput}}
	txCoinbase.SetTxID()
	return txCoinbase
}

func (tx Transaction) SetTxID() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln(err)
	}
	buffBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), buff.Bytes()}, []byte{})
	hash := sha256.Sum256(buffBytes)
	tx.TxID = hash[:]
}

func NewSimpleTransaction(from, to string, amount int64, bc *BlockChain, txs []*Transaction) *Transaction {
	var (
		txInputs  []*TXInput
		txOutputs []*TXOuput
	)
	balance, spendableUTXO := bc.FindSpendableUTXOs(from, amount, txs)

	// 代表消费
	for txID, indexArry := range spendableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArry {
			txInput := &TXInput{txIDBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	// 转账
	txOutput1 := &TXOuput{amount, to}
	txOutputs = append(txOutputs, txOutput1)

	// 找零
	txOutput2 := &TXOuput{balance - amount, from}
	txOutputs = append(txOutputs, txOutput2)
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.SetTxID()
	return tx
}

// 判断当前交易是否是coinbase交易
func (tx Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}
