package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

/*
将一个int64的整数：转换为二进制后，再转为[]byte
*/
func IntToHex(num int64) []byte {
	buf := new(bytes.Buffer)

	// 将二进制数据写入io.Writer
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}
	
	// 转为[]byte并返回
	return buf.Bytes()
}
