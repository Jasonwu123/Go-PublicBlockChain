package BLC

type TXInput struct {
	TxID      []byte // 交易ID
	Vout      int    // 存储TxOutput的vout里面的索引
	ScriptSiq string // 用户名
}

// 判断当前TxInput消费，和指定的address是否一致
func (txInput *TXInput) UnLockWithAddress(address string) bool {
	return txInput.ScriptSiq == address
}
