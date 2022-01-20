package main

import "BLC"

func main() {

	//// Open the my.db data file in your current directory.
	//// It will be created if it doesn't exist.
	//// db, err := bolt.Open("my.db", 0600, nil)
	//// if err != nil {
	//// 	log.Fatal(err)
	//// }
	//// defer db.Close()
	//
	///*
	//	Update(), 读写
	//	View(), 只读
	//*/
	//
	///*
	//	// 1. 创建表
	//	err = db.Update(func(t *bolt.Tx) error {
	//		// 1. 创建bucket
	//		bk, err := t.CreateBucket([]byte("MyBucket"))
	//		if err != nil {
	//			return fmt.Errorf("create bucket: %s", err)
	//		}
	//		// 2. 向bucket中存储数据
	//		if bk != nil {
	//			err = bk.Put([]byte("1"), []byte("send 1000 BTC to 王二狗"))
	//			if err != nil {
	//				log.Panic("数据存储错误吧。")
	//			}
	//		}
	//
	//		return nil
	//	})
	//
	//	if err != nil {
	//		log.Panic("创建表失败。")
	//	}
	//
	//	// 读取数据
	//	err = db.View(func(t *bolt.Tx) error {
	//		bk := t.Bucket([]byte("MyBucket"))
	//		if bk != nil {
	//			// 根据key查看数据
	//			data := bk.Get([]byte("1"))
	//			fmt.Println(data)
	//			fmt.Printf("%s\n", data)
	//
	//			// 获取不存在的key
	//			data2 := bk.Get([]byte("2"))
	//			fmt.Println(data2)
	//			fmt.Printf("%s\n", data2)
	//		}
	//		return nil
	//	})
	//
	//	if err != nil {
	//		log.Panic("读取表数据失败。")
	//	}
	//*/
	//
	//// err = db.View(func(t *bolt.Tx) error {
	//// 	bk := t.Bucket([]byte("MyBucket"))
	//// 	bkc := bk.Cursor()
	//// 	for k, v := bkc.First(); k != nil; k, v = bkc.Next() {
	//// 		fmt.Printf("key = %s, value = %s\n", k, v)
	//// 		}
	//// 	return nil
	//// })
	//
	///*
	//	// 创建区块，存入数据库
	//	// 打开数据库
	//	block := BLC.NewBlock("helloworld", make([]byte, 32, 32), 0)
	//
	//	db, err := bolt.Open("my.db", 0600, nil)
	//	if err != nil {
	//		log.Fatal("数据库打开失败")
	//	}
	//
	//	defer db.Close()
	//
	//	// 存储一个block区块
	//	err = db.Update(func(t *bolt.Tx) error {
	//		bk := t.Bucket([]byte("blocks"))
	//		if bk == nil {
	//			bk, err = t.CreateBucket([]byte("blocks"))
	//			if err != nil {
	//				log.Panic("创建bucket失败")
	//			}
	//		}
	//
	//		// 添加数据
	//		err = bk.Put([]byte("1"), block.Serialize())
	//		if err != nil {
	//			log.Panic("添加数据失败")
	//		}
	//
	//		return nil
	//	})
	//
	//	if err != nil {
	//		log.Panic("打开数据库失败")
	//	}
	//
	//	// 从数据库中读取该区块数据
	//	err = db.View(func(t *bolt.Tx) error {
	//		bk := t.Bucket([]byte("blocks"))
	//		if bk != nil {
	//			data := bk.Get([]byte("1"))
	//
	//			// 反序列化
	//			blockData := BLC.DeserializeBlock(data)
	//			// fmt.Println(blockData)
	//			fmt.Printf("%v\n", blockData)
	//		}
	//		return nil
	//	})
	//
	//	if err != nil {
	//		log.Panic("数据库查询失败")
	//	}
	//*/
	//
	//// 测创世块存入数据库
	//blockchain := BLC.CreateBlcokChainWithGensisBlcok("Genesis Block")
	//fmt.Println(blockchain)
	//defer blockchain.DB.Close()
	//
	////8.测试新添加的区块
	//blockchain.AddBlockToBlockChain("Send 100RMB to wangergou")
	//blockchain.AddBlockToBlockChain("Send 100RMB to lixiaohua")
	//blockchain.AddBlockToBlockChain("Send 100RMB to rose")
	//fmt.Println(blockchain)
	//blockchain.PrintChains()

	// 9. CLI操作
	cli := BLC.CLI{}
	cli.Run()
}
