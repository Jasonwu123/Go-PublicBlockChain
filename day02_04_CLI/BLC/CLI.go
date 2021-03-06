package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 定义CLI结构体
type CLI struct {}

// 添加Run方法
func (cli *CLI) Run()  {
	// 判断命令行参数的长度
	isValiArgs()

	// 创建flagset标签对象
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	// 设置标签后的参数
	flagAddBlockData := addBlockCmd.String("data", "helloworld", "交易数据")
	flagCreateBlockChainData := createBlockChainCmd.String("data", "Genesis block data", "创世块交易数据")

	// 解析
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChains()
	}

	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockChainData)
	}
}

func isValiArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage()  {
	fmt.Println("Usage: ")
	fmt.Println("\tcreateblockchain -data DATA -- 创建创世块")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块链信息")
}

func (cli *CLI) addBlock(data string)  {
	bc := GetBlockchainObject()
	if bc == nil {
		fmt.Println("没有创始块，无法添加。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.AddBlockToBlockChain(data)
}

func (cli *CLI) printChains()  {
	bc := GetBlockchainObject()
	if bc == nil {
		fmt.Println("没有区块可以打印。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.PrintChains()
}

func (cli *CLI) createGenesisBlockchain(data string)  {
	CreateBlcokChainWithGensisBlcok(data)
}


