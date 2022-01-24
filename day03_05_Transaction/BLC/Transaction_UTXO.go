package BLC

// 创建UTXO结构体，表示所有未花费的
type UTXO struct {
	TxID   []byte // 当前Transaction的交易ID
	Index  int    // 下标索引
	Output *TXOuput
}
