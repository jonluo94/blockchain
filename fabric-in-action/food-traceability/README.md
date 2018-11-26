### food-traceability 食品溯源实战
> 逐步开发实现
#### 搭建网络
* 创建项目 food-traceability
* 创建 network 目录
* 搭建 fabric 网络工具：bin文件夹下 (具体来源fabric-samples项目)
* 在channel下编写configtx.yaml，cryptogen.yaml，fabric-ca-server-config.yaml 文件
  * 注意 fabric-ca-server-config.yaml 下的 affiliations节点的组织与网络对应
* 执行 create-channel-artifacts.sh 创建证书和密钥文件
* 在docker-compose下编写base.yaml，docker-compose.yaml 文件
  * 注意 base.yaml的 CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE 要与docker-compose的启动网络对应
  * 注意 docker-compose.yaml的 *_KEYFILE 修改为相应的密钥文件
* 执行 start-network.sh
* 完成网络
#### 编写chaincode
* 创建 chaincode
* 在 src 下编写chaincode源代码，每个目录一个chaincode
* 相应目录下编写测试用例
* 相应目录下测试： go test -v ..._test.go ...go
* 测试成功，完成chaincode
#### 编写app
* 创建 app 目录
* 在 lib 下添加相应的 fabric-samples/balance-transfer/app的文件
* 编写 package.json
* 在 config 下编写网络配置 network-config.yaml 和各组织的用户密钥证书配置org1.yaml，org2.yaml，org3.yaml
* 在 config.js 引入config的配置
* 在 config.json 编写服务端配置
* 编写app.js
* 安装相应包：npm install
* 运行：node app.js
* 编写运行建本 bulid.sh 创建用户，创建channel，加入channel，安装chaincode，实例化chaincode
* 完成服务端