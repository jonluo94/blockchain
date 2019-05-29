package main


import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
)

//const ccPath  = "github.com/example_cc"
const ccPath  = "github.com/jonluo94/fabric-go-sdk-examples/chaincode"

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

	//ccPkg, err := gopackager.NewCCPackage(ccPath, integration.GetDeployPath())

	ccPkg, err := gopackager.NewCCPackage(ccPath, os.Getenv("GOPATH"))
	if err != nil {
		fmt.Println(err)
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: "mycc", Path: ccPath, Version: "0", Package: ccPkg}
	res, err := resmgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := resmgmtClient.InstantiateCC(
		"mychannel",
		resmgmt.InstantiateCCRequest{
			Name: "mycc",
			Path: ccPath,
			Version: "0",
			Args: [][]byte{[]byte("init"),[]byte("a"), []byte("100"), []byte("b"), []byte("200")},
			Policy: ccPolicy,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	fmt.Println(resp)

}


