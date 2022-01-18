package BLC

// step5: 创建区块链
type BlockChain struct {
	Blocks []*Block
}

//step6：创建区块链，带有创世块
func CreateBlcokChainWithGensisBlcok(data string) *BlockChain {
	// 创建创世块
	genesisBlock := CreateGenesisBlock(data)

	// 返回区块链对象
	return &BlockChain{[]*Block{genesisBlock}}
}

//step7：添加一个新的区块到区块链中
func (bc *BlockChain) AddBlockToBlockChain(data string, prevHash []byte, height int64) {
	// 创建新区块
	newBlock := NewBlock(data, prevHash, height)

	// 添加到slice中
	bc.Blocks = append(bc.Blocks, newBlock)
}
