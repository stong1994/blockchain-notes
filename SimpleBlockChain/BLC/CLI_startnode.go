package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) startNode(nodeID string, minerAddr string) {
	// 启动服务器
	if minerAddr == "" || IsValidForAddress([]byte(minerAddr)) {
		// 启动服务器
		fmt.Println("启动服务器:localhost:%s\n", nodeID)
		startServer(nodeID, minerAddr)
	} else {
		fmt.Println("指定的地址无效")
		os.Exit(0)
	}
}
