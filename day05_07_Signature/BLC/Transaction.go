package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
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
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTXOuput(10, address)
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

	// 获取钱包
	wallets := NewWallets()
	wallet := wallets.WalletsMap[from]

	// 代表消费
	for txID, indexArry := range spendableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArry {
			txInput := &TXInput{txIDBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}

	// 转账
	txOutput1 := NewTXOuput(amount, to)
	txOutputs = append(txOutputs, txOutput1)

	// 找零
	txOutput2 := NewTXOuput(balance-amount, from)
	txOutputs = append(txOutputs, txOutput2)
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.SetTxID()
	return tx
}

// 判断当前交易是否是coinbase交易
func (tx Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

/*
签名
正如上面提到的，为了对一笔交易进行签名，我们需要获取交易输入所引用输出
因为我们需要存储这些输出的交易
*/
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 如果是coinbase交易，无需签名
	if tx.IsCoinbaseTransaction() {
		return
	}

	// input没有对应的transaction，无法签名
	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxID)].TxID == nil {
			log.Panicln("当前的input没有对应的transa")
		}
	}

	// 获取transaction的部分数据副本
	txCopy := tx.TrimmedCopy()

	for index, input := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(input.TxID)]
		input.Signature = nil
		input.PublicKey = prevTx.Vouts[input.Vout].PubKeyHash
		data := txCopy.getData()
		input.PublicKey = nil

		/*
				签名
				通过privkey对txcopy.id进行签名
			    一个 ECDSA 签名就是一对数字，我们对这对数字连接起来，并存储在输入的 Signature 字段。
		*/
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, data)
		if err != nil {
			log.Panicln(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[index].Signature = signature
	}
}

/*
获取签名所需要的Transaction副本，需要剪裁
这个副本包含了所有的输入和输出，但是 TXInput.Signature 和 TXIput.PubKey 被设置为 nil。
*/
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOuput
	for _, input := range tx.Vins {
		inputs = append(inputs, &TXInput{input.TxID, input.Vout, nil, nil})
	}

	for _, output := range tx.Vouts {
		outputs = append(outputs, &TXOuput{output.Value, output.PubKeyHash})
	}

	txCopy := Transaction{tx.TxID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Serialize() []byte {
	jsonByte, err := json.Marshal(tx)
	if err != nil {
		log.Panicln(err)
	}
	return jsonByte
}

func (tx Transaction) getData() []byte {
	txCopy := tx
	txCopy.TxID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

/*
验证数字签名
*/
func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	// input没有对应的transaction，无法签名
	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxID)].TxID == nil {
			log.Panicln("当前input没有对应transaction，无法验证。")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()
	for index, input := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(input.TxID)]
		txCopy.Vins[index].Signature = nil
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash
		data := txCopy.getData()
		txCopy.Vins[index].PublicKey = nil

		// 签名中s和r
		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:sigLen/2])
		s.SetBytes(input.Signature[sigLen/2:])

		// 通过公钥，产生新的s和r，与原来的进行对比
		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])

		// 根据椭圆曲线算法，以及x，y获取公钥
		// 我们使用输入提取公钥创建了一个ecdsa.PublicKey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		//这里我们解包存储在 TXInput.Signature 和 TXInput.PubKey 中的值，
		// 因为一个签名就是一对数字，一个公钥就是一对坐标。
		// 我们之前为了存储将它们连接在一起，现在我们需要对它们进行解包在 crypto/ecdsa 函数中使用。

		//验证
		//在这里：我们使用从输入提取的公钥创建了一个 ecdsa.PublicKey，通过传入输入中提取的签名执行了 ecdsa.Verify。
		// 如果所有的输入都被验证，返回 true；如果有任何一个验证失败，返回 false.
		if ecdsa.Verify(&rawPubKey, data, &r, &s) == false {
			return false
		}
	}
	return true
}
