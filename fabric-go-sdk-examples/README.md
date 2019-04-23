### fabric-go-sdk-examples
* fabric 1.4
* fabric-samples
* go 1.11.6

### 步骤
* go get github.com/hyperledger/fabric-samples
* cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
  * 修改背书策略(修改 scripts/utils.sh)
    ```
    "AND ('Org1MSP.peer','Org2MSP.peer')" 
    改为
    "OR ('Org1MSP.peer','Org2MSP.peer')"
    ```
* 替换chaincode 
  * 将 chaincode/chaincode_example02.go 替换 fabric-samples下的chaincode/chaincode_example02/go/chaincode_example02.go

* 启动网络
  ``` 
  ./byfn.sh generate
  ./byfn.sh up
  ```
* 测试
  * invoke-chaincode 
    调用链码
  * query-chaincode
    查询链码
  * query-ledger
    查询账本


