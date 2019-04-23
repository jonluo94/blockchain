/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
	"encoding/json"
)

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Printf("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func main() {
	sdk, err := fabsdk.New(config.FromFile("./first-network.yaml"))
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to create SDK"))
		os.Exit(-1)
	}
	defer sdk.Close()

	user := "User1"
	org := "Org1"
	channelName := "mychannel"
	ccid := "mycc"
	eventId := "fnevent"

	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	// client for interacting with a chaincode
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	// 注册事件
	reg, notifier := regitserEvent(client, ccid, eventId)
	defer client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID: ccid,
		Fcn:         "invoke",
		Args:        [][]byte{[]byte("a"), []byte("b"), []byte("10"), []byte(eventId)},
	}
	resp, err := client.Execute(req)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	err = eventResult(notifier, eventId)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	res, err := json.Marshal(resp)
	fmt.Println(string(res))

}
