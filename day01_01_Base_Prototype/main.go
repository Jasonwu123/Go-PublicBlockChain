package main

import (
	"BLC"
	"fmt"
)

func main() {
	// 1. 测试Block
	// block := BLC.NewBlock("I am a blcok", make([]byte, 32, 32), 1)
	// fmt.Printf("区块高度为：%x\n", block.Height)
	// fmt.Printf("交易数据为：%s\n", block.Data)

	// 2. 测试创世块
	// genesisBlcok := BLC.CreateGenesisBlock("创世块")
	// fmt.Printf("区块高度为：%x\n", genesisBlcok.Height)
	// fmt.Printf("前区块哈希为：%x\n", genesisBlcok.PrevBlockHash)
	// fmt.Printf("交易数据为：%s\n", genesisBlcok.Data)

	//// 3. 测试区块链
	// genesisBlcokChain := BLC.CreateBlcokChainWithGensisBlcok("genesisBlockChain")
	// fmt.Println(genesisBlcokChain)
	// fmt.Println(genesisBlcokChain.Blocks)
	// fmt.Printf("区块高度为：%x\n", genesisBlcokChain.Blocks[0].Height)
	// fmt.Printf("前区块哈希为：%x\n", genesisBlcokChain.Blocks[0].PrevBlockHash)
	// fmt.Printf("交易数据为：%s\n", genesisBlcokChain.Blocks[0].Data)

	// 4. 测试添加新区块
	blockChain := BLC.CreateBlcokChainWithGensisBlcok("Genesis Block...")
	blockChain.AddBlockToBlockChain("Send 1 btc to Wangergou", blockChain.Blocks[len(blockChain.Blocks)-1].Hash, blockChain.Blocks[len(blockChain.Blocks)-1].Height+1)
	blockChain.AddBlockToBlockChain("Send 3 btc to Lixiaohua", blockChain.Blocks[len(blockChain.Blocks)-1].Hash, blockChain.Blocks[len(blockChain.Blocks)-1].Height+1)
	blockChain.AddBlockToBlockChain("Send 10 btc to Rose", blockChain.Blocks[len(blockChain.Blocks)-1].Hash, blockChain.Blocks[len(blockChain.Blocks)-1].Height+1)

	for _, block := range blockChain.Blocks {
		fmt.Printf("Block height: %x\n", block.Height)
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
