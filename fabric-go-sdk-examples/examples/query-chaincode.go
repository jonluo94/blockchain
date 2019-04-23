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
)

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

	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	// client for interacting with a chaincode
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	// this query is specific to fabric-chaincode-evm
	resp, err := client.Query(channel.Request{
		ChaincodeID: ccid,
		Fcn:         "query",
		Args:        [][]byte{[]byte("a")},
	})

	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	fmt.Println(string(resp.Payload))
}