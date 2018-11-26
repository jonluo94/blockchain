#### 在 hlf-network 下依次执行创建所需文件
```
./bin/cryptogen generate --config=./crypto-config.yaml
```
```
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ChainProfile -outputBlock ./channel-artifacts/orderer.genesis.block
```
```
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ChainProfile -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID jonluo-chain
```
```
FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile ChainProfile -outputAnchorPeersUpdate ./channel-artifacts/org1.anchors.tx -channelID jonluo-chain -asOrg Org1
```
#### 编写 docker-compose.yaml 文件
* 运行： docker-compose up
* 结束： docker-compose down
* 清楚所有容器和chaincode镜像：./cleanDocker.sh