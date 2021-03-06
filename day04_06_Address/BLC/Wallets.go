package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 创建钱包
type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建一个钱包集合: 文件中存在就从文件中读取，否则新建一个
const walletFile = "Wallets.dat"

func NewWallets() *Wallets {
	// 判断钱包文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		fmt.Println("钱包文件不存在")
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets
	}

	// 否则读取文件中的数据
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicln(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicln(err)
	}

	return &wallets
}

// 创建一个新钱包
func (ws *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Printf("创建钱包地址：%s\n", wallet.GetAddress())
	ws.WalletsMap[string(wallet.GetAddress())] = wallet

	// 保存钱包
	ws.SaveWallets()
}

/*
要让数据对象在网络上传输或存储，我们需要进行编码和解码。
现在比较流行的编码方式有：JSON,XML等。然而，Go在gob包中提供另一种方式，编码解码效率高于JSON
gob包是Go自带的一个数据结构序列化的编码/解码工具包
*/
func (ws *Wallets) SaveWallets() {
	var content bytes.Buffer

	// 注册的目的是为了可以序列化任何类型，wallet结构体中有接口类型，将接口进行注册
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panicln(err)
	}

	// 将序列化后的数据写入文件，源文件内容会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panicln(err)
	}
}
