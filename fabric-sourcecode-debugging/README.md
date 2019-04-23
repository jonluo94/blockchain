#### fabric-sourcecode-debugging
借助开发网络调试 fabric 源码本地调试

* IDE Goland
* Go 1.9.7
* fabric-samples 模块 chaincode-docker-devmode
* fabric 源码

#### 步骤
* 添加本地域名  
  127.0.0.1       peer  
  127.0.0.1       orderer  
* 用 ide 打开 gopath 下的fabric源码目录
* 在源码目录下添加 dev-network
* 把 sampleconfig 下的所有文件复制到 dev-network 
  * 修改 core.yaml 中 fileSystemPath = fabric源码目录下dev-network/production/peer(绝对路径)
  * 修改 orderer.yaml 中 Location = fabric源码目录下dev-network/production/orderer(绝对路径)
* 在 dev-network 新建 config，并复制 fabric-samples 模块 chaincode-docker-devmode下的 myc.tx 和 orderer.block
* 1，接着调试网络，用debug模式运行 order 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行

  Environment：  
  ```
  ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
  ORDERER_GENERAL_GENESISMETHOD=file
  ORDERER_GENERAL_GENESISFILE=fabric源码目录下dev-network/config/orderer.block(绝对路径)
  ORDERER_GENERAL_LOCALMSPID=DEFAULT
  ORDERER_GENERAL_LOCALMSPDIR=fabric源码目录下dev-network/msp(绝对路径)
  FABRIC_CFG_PATH=fabric源码目录下dev-network(绝对路径)
  ```
  
* 2，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  node start --peer-chaincodedev=true -o 127.0.0.1:7050
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=peer
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```

* 3，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  channel create -c myc -f fabric源码目录下dev-network/config/myc.tx(绝对路径) -o 127.0.0.1:7050
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```

* 4，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  channel join -b myc.block
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```

* 5，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  chaincode install -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 -n mycc -v 1.0
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```
  
* 6，打开终端  
  cd $GOPATH/src/github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02  
  编译chaincode   
  go build -o chaincode_example02 
  接着运行  
  CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:1.0 ./chaincode_example02
  
* 7，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  chaincode instantiate -n mycc -v 1.0 -c "{\"Args\":[\"init\",\"a\",\"100\",\"b\",\"200\"]}" -C myc
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```
  
* 8，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  chaincode invoke -n mycc -c "{\"Args\":[\"invoke\",\"a\",\"b\",\"10\"]}" -C myc
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```
  
* 9，接着调试网络，用debug模式运行 peer 下的 main.go 文件 并添加配置，绝对路径的地方相对修改，然后运行
  
  Program arguments：
  ```
  chaincode query -n mycc -c "{\"Args\":[\"query\",\"a\"]}" -C myc
  ```
  
  Environment： 
  ```
  CORE_PEER_LOCALMSPID=DEFAULT
  CORE_PEER_ID=cli
  CORE_PEER_MSPCONFIGPATH=fabric源码目录下dev-network/msp(绝对路径)
  CORE_PEER_ADDRESS=127.0.0.1:7051
  FABRIC_CFG_PATH=fabric源码目录下/dev-network(绝对路径)
  ```
  
  看到 ``Query Result: 90`` 表示成功