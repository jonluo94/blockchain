## chaincode-app (fabric-sdk-go 未正式发布，目前不推荐)
* Hyperledger Fabric SDK Go 简单app
  * hlf-network -- 网络
  * hello-service -- sdk 代码实现
  * chaincode -- 链码
* 运行步骤
* cd $GOPATH/src
* 创建 gitee.com 文件夹
* 将 chaincode-app 整个文件夹复制到 gitee.com 文件夹下
* 打开两个终端
* 第一个终端在 hlf-network 下运行：docker-compose up
* 第二个终端在 hello-service 下运行：go run main.go
* 在浏览器打开 localhost:8888 访问

> 由于 Go SDK 未正式发布的，目前有很多坑，填坑中......

> 基于 fabric版本1.1.0 到 1.2.0 已退休

