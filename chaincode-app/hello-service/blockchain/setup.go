package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"time"
)

// FabricSetup 实现
type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          chclient.ChannelClient
	admin           resmgmt.ResourceMgmtClient
	sdk             *fabsdk.FabricSDK
}

// 读取配置文件初始化 SDK 客户端、链和事件
func (setup *FabricSetup) Initialize() error {

	// 检查 SDK 是否已经初始化
	if setup.initialized {
		return fmt.Errorf("SDK已经初始化")
	}

	// 通过配置文件初始化 SDK
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("创建SDK失败: %v", err)
	}
	setup.sdk = sdk

	// 通道管理客户端负责管理通道（创建/更新通道）
	// 提供具有创建通道的特权的用户（一般是订购者管理）
	chMgmtClient, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).ChannelMgmt()
	if err != nil {
		return fmt.Errorf("添加管理员用户到SDK失败: %v", err)
	}

	// 组织管理员用户为正在创建的通道签名用户
	// 会话方法是目前获取用户标识的唯一方法。
	session, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).Session()
	if err != nil {
		return fmt.Errorf("获得会话失败 %s, %s: %s", setup.OrgName, setup.OrgAdmin, err)
	}
	orgAdminUser := session

	// 创建通道 jonluo-chain
	// 通道由ITS（组织、锚节点，共享的分类帐，链码应用和订购服务）定义
	// 网络上的每个交易都在一个通道上执行
	req := chmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfig: setup.ChannelConfig, SigningIdentity: orgAdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		return fmt.Errorf("创建通道失败: %v", err)
	}

	// 允许订购者处理通道创建，等待5秒
	time.Sleep(time.Second * 5)

	// 资源管理客户端是管理系统资源的客户端API
	// 它将允许我们直接与区块链进行交互。它可以与管理员状态相关联。
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("创建资源管理客户端失败: %v", err)
	}

	// 组织节点加入通道
	if err = setup.admin.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("组织节点加入通道失败: %v", err)
	}

	fmt.Println("通道初始化成功")
	setup.initialized = true
	return nil
}

// 安装和实例化chaincode
func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// 打包go链码
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("打包链码失败: %v", err)
	}

	// 在组织节点上安装链码
	// 资源管理客户端将链码发送到通道中的所有节点，以便它们存储链码并在以后与它交互
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("在组织对等节点上安装链码失败：%v", err)
	}

	// 建立链码策略
    // 如果你的交易必须遵循特定的规则，则链码策略是必须的
    // 如果你不提供任何政策，每笔交易都会被认可，这可能不是你想要的。
    // 在这种情况下，我们将规则设置为：如果交易已由组织id为“org1.jonluo.com”的成员签名，则认可交易
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.jonluo.com"})

	// 在组织节点上实例化链码
	// 资源管理客户端告诉其通道中的所有对等节点实例化前先安装的链码
	err = setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Args: [][]byte{[]byte("hello"),[]byte("world")}, Policy: ccPolicy})
	if err != nil {
		return fmt.Errorf("实例化链码失败: %v", err)
	}

	// 通道客户端用于查询和执行交易
	setup.client, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName)).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("通道客户端创建失败: %v", err)
	}

	fmt.Println("链码安装和实例化成功")
	return nil
}
