package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

/*
Json字符串转[]sting数组
*/
func JSONToArray(jsonString string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panicln(err)
	}
	return sArr
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
