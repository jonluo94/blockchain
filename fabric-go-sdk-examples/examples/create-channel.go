package main


import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

func main() {

	sdk, err := fabsdk.New(config.FromFile("./first-network.yaml"))
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to create SDK"))
		os.Exit(-1)
	}
	defer sdk.Close()

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resmgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		fmt.Printf("Failed to create channel management client: %s", err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("Org1"))
	if err != nil {
		fmt.Println(err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		fmt.Println(err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: "mychannel",
		ChannelConfigPath: "/home/jonluo/gopath/src/github.com/hyperledger/fabric-samples/first-network/channel-artifacts/channel.tx",
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resmgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	fmt.Println(txID)

}