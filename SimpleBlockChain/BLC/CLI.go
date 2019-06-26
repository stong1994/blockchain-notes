package BLC

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

type CLI struct {
}

func (cli *CLI) Run() {
	// 判断命令行参数
	isValidArgs()
	/*
		获取节点ID,
		返回当前进程的环境变量varname的值,若变量没有定义时返回nil
		export NODE_ID=8888
		每次打开一个终端,都需要设置NODE_ID的值
	*/
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID环境变量没有设置..\n")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID: %s\n", nodeID)

	// 1.创建flagset标签对象
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	addressListsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	// 2.设置标签后的参数
	flagFromData := sendBlockCmd.String("from", "", "转账源地址")
	flagToData := sendBlockCmd.String("to", "", "转账目标地址")
	flagAmountData := sendBlockCmd.String("amount", "", "转账金额")
	flagCreateBlockChainData := createBlockCmd.String("address", "", "创世区块交易地址")
	//flagAddBlockData := addBlockCmd.String("data", "helloworld...", "交易数据")
	flagGetBalanceData := getBalanceCmd.String("address", "", "要查看的账户的余额")

	flagMiner := startNodeCmd.String("miner", "", "定义挖矿奖励的地址....")
	flagMine := sendBlockCmd.Bool("mine", false, "是否在当前节点中立即验证...")

	// 3.解析
	switch os.Args[1] {
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
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
		err := createBlockCmd.Parse(os.Args[2:])
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
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
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
		from := JSONToArray(*flagFromData)
		to := JSONToArray(*flagToData)
		amount := JSONToArray(*flagAmountData)

		for i := 0; i < len(from); i++ {
			if !IsValidForAddress([]byte(from[i])) || !IsValidForAddress([]byte(to[i])) {
				fmt.Println("钱包地址无效")
				printUsage()
				os.Exit(1)
			}
		}

		cli.send(from, to, amount, nodeID, *flagMine)
	}
	if printChainCmd.Parsed() {
		cli.printChains(nodeID)
	}
	if createBlockCmd.Parsed() {
		if !IsValidForAddress([]byte(*flagCreateBlockChainData)) {
			fmt.Println("创建地址无效")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockChainData, nodeID)
	}

	if getBalanceCmd.Parsed() {
		if *flagGetBalanceData == "" {
			fmt.Println("查询地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceData, nodeID)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}
	if addressListsCmd.Parsed() {
		cli.addresslists(nodeID)
	}
	if testCmd.Parsed() {
		cli.TestMethod(nodeID)
	}
	if startNodeCmd.Parsed() {
		cli.startNode(nodeID, *flagMiner)
	}
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\taddresslists -- 输出所有钱包地址")
	fmt.Println("\tcreateblockchain -address DATA -- 创建创世区块")
	fmt.Println("\tsend -from From -to To -amount Amount - 交易数据")
	fmt.Println("\tprintchain -- 输出信息")
	fmt.Println("\tgetbalance -address DATA -- 查询账户余额")
	fmt.Println("\ttest -- 测试")
	fmt.Println("\tstartnode -miner ADDRESS -- 启动节点服务器,并且指定挖矿奖励的地址.")
}
