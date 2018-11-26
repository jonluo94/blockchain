package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

// 查询链码key的值
func (setup *FabricSetup) Query(key string) (string, error) {

	// 准备参数
	var args []string
	args = append(args, "get")
	args = append(args, key)

	response, err:= setup.client.Query(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn:args[0], Args: [][]byte{[]byte(args[1])}})

	if err != nil {
		return "", fmt.Errorf("查询失败: %v", err)
	}

	return string(response.Payload), nil
}