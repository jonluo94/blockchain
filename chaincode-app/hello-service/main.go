package main

import (
	"fmt"
	"gitee.com/chaincode-app/hello-service/blockchain"
	"gitee.com/chaincode-app/hello-service/web"
	"os"
	"time"
)

func main() {
	// Fabric SDK 属性定义
	fSetup := blockchain.FabricSetup{
		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		ConfigFile:      "config.yaml",

		// 通道参数
		ChannelID:       "jonluo-chain",
		ChannelConfig:   os.Getenv("GOPATH") + "/src/gitee.com/chaincode-app/hlf-network/channel-artifacts/channel.tx",

		// Chaincode parameters
		ChainCodeID:     "hellocc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "gitee.com/chaincode-app/chaincode",

		// 用户参数
		UserName: "User1",
	}

	// 从先前设置的属性初始化SDK
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("初始化 Fabric SDK 失败: %v\n", err)
	}

	// 安装和实例化链码
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("安装和实例化链码失败: %v\n", err)
	}

	//测试
	//test(&fSetup)

	web.Run(&fSetup)
}

//测试方法
func test(fSetup *blockchain.FabricSetup)  {
	time.Sleep(time.Second * 5)
	// 查询链码
	response, err := fSetup.Query("hello")
	if err != nil {
		fmt.Printf("查询 hello 失败: %v\n", err)
	} else {
		fmt.Printf("查询 hello 的值: %s\n", response)
	}

	// 调用链码
	txId, err := fSetup.Invoke("set","hello","jonluo")
	if err != nil {
		fmt.Printf("调用链码失败: %v\n", err)
	} else {
		fmt.Printf("调用成功，交易ID: %s\n", txId)
	}

	// Query again the chaincode
	response, err = fSetup.Query("hello")
	if err != nil {
		fmt.Printf("查询 hello 失败: %v\n", err)
	} else {
		fmt.Printf("查询 hello 的值: %s\n", response)
	}
}