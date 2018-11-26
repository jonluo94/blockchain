package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"fmt"
	"time"
)

// 更改状态
func (setup *FabricSetup) Invoke(fnName string,key string,value string) (string, error) {

	// 准备参数
	var args []string
	args = append(args, fnName)
	args = append(args, key)
	args = append(args, value)

    //链码中发出的事件id
	eventID := "eventInvokeSet"
	// 在客户端上注册通知处理程序
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("注册链码事件失败: %v", err)
	}

	// 添加交易数据描述，例如调用请求的描述备注
	transientDataMap := make(map[string][]byte)
	transientDataMap["remark"] = []byte("set方法调用")

	// 创建请求并发送
	response, err := setup.client.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("交易失败: %v", err)
	}

	// 等待提交的结果
	select {
	    case ccEvent := <-notifier:
		     fmt.Printf("接受链码事件: %s\n", ccEvent)
	    case <-time.After(time.Second * 20):
	    	 return "", fmt.Errorf("没有接收到链码事件(%s)", eventID)
	}

	// 注销以前在客户端创建的通知处理程序
	err = setup.client.UnregisterChaincodeEvent(rce)
    //返回交易id
	return response.TransactionID.ID, nil
}