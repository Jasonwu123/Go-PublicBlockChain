package BLC

type TXOuput struct {
	Value        int64
	ScriptPubKey string // 公钥：可理解为用户名
}

// 判断当前TXOutpu消费，和指定的address是否一致
func (txOutput *TXOuput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
