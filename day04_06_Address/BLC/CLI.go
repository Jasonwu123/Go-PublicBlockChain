package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 定义CLI结构体
type CLI struct{}

// 添加Run方法
func (cli *CLI) Run() {
	// 判断命令行参数的长度
	isValiArgs()

	// 创建flagset标签对象
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	addressListsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	// 设置标签后的参数
	flagFromData := sendBlockCmd.String("from", "", "转账源地址")
	flagToData := sendBlockCmd.String("to", "", "转账目标地址")
	flagAmountData := sendBlockCmd.String("amount", "", "转账金额")
	flagCreateBlockChainData := createBlockChainCmd.String("address", "", "创世块交易地址")
	flagGetBalanceData := getBalanceCmd.String("address", "", "要查询的某个账户的余额")

	// 解析
	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
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

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "addresslists":
		err := addressListsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFromData == "" || *flagToData == "" || *flagAmountData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagFromData)
		fmt.Println(*flagToData)
		fmt.Println(*flagAmountData)

		from := JSONToArray(*flagFromData)
		to := JSONToArray(*flagToData)
		amount := JSONToArray(*flagAmountData)

		for i := 0; i < len(from); i++ {
			if !IsValidForAddress([]byte(from[i])) || !IsValidForAddress([]byte(to[i])) {
				fmt.Println("钱包地址无效。")
				printUsage()
				os.Exit(1)
			}
		}

		cli.send(from, to, amount)
	}

	if printChainCmd.Parsed() {
		cli.printChains()
	}

	if createBlockChainCmd.Parsed() {
		if !IsValidForAddress([]byte(*flagCreateBlockChainData)) {
			fmt.Println("创建地址无效。")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockChainData)
	}

	if getBalanceCmd.Parsed() {
		if !IsValidForAddress([]byte(*flagGetBalanceData)) {
			fmt.Println("查询地址无效。")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceData)
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet()
	}

	if addressListsCmd.Parsed() {
		// 获取所有的钱包地址
		cli.addressLists()
	}
}

func isValiArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\taddresslists -- 列出所有钱包地址")
	fmt.Println("\tcreateblockchain -data DATA -- 创建创世块")
	fmt.Println("\tsend -from From -to To -amount Amount -data -- 交易数据")
	fmt.Println("\tprintchain -- 输出信息")
	fmt.Println("\tgetbalance -address DATA -- 查询账户余额")
}
