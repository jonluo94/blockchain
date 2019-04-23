/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
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

	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	// client for interacting directly with the ledger
	ledger, err := ledger.New(clientChannelContext)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	bci, err := ledger.QueryInfo()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	fmt.Println(bci)
}