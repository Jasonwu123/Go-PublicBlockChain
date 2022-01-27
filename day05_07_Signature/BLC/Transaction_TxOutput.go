package BLC

import "bytes"

type TXOuput struct {
	Value      int64
	PubKeyHash []byte // 公钥
}

// 判断当前TXOutpu消费，和指定的address是否一致
func (txOutput *TXOuput) UnLockWithAddress(address string) bool {
	fullPayloadHash := Base58Decode([]byte(address))
	pubKeyHash := fullPayloadHash[1 : len(fullPayloadHash)-4]
	return bytes.Compare(txOutput.PubKeyHash, pubKeyHash) == 0
}

func NewTXOuput(value int64, address string) *TXOuput {
	txOuput := &TXOuput{value, nil}
	txOuput.Lock(address)
	return txOuput
}

func (txOutput TXOuput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	txOutput.PubKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
}
