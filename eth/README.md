#### 以太坊
- [以太坊白皮书](http://ethfans.org/wikis/以太坊白皮书) 讲解以太坊设计思路
- [以太坊黄皮书](https://ethereum.github.io/yellowpaper/paper.pdf) 讲解以太坊底层技术
- [以太坊设计原理](http://ethfans.org/posts/510)
- [以太坊代码剖析](http://ethfans.org/topics/227)
- [以太坊源码分析：地址如何生成的](http://baijiahao.baidu.com/s?id=1597007700502456856&wfr=spider&for=pc)
- [以太坊crypto](http://www.cnblogs.com/baizx/p/6936258.html) 签名,验证,以及公钥与以太坊地址转换
- [以太坊RLP](http://www.cnblogs.com/baizx/p/6928622.html) 对象进行序列化的主要编码方式
- [以太坊的代码来学习区块链技术](https://github.com/ZtesoftCS/go-ethereum-code-analysis)
- [以太坊开发入门](http://me.tryblockchain.org/getting-up-to-speed-on-ethereum.html)
- [以太坊Gas目前定价](https://ethgasstation.info/index.php) 
- [Solidity校验椭圆曲线加密数字签名（附实例）](http://www.toutiao.com/i6401418700217385473/?tt_from=weixin&utm_campaign=client_share&from=groupmessage&app=news_article&utm_source=weixin&iid=8932715408&utm_medium=toutiao_ios&wxshare_count=2&pbid=35867484354)
#### 以太坊钱包
- [My Ether Wallet](https://myetherwallet.com) 网页版以太坊钱包 [源码](https://github.com/kvhnuke/etherwallet)
- [MetaMask](https://metamask.io/) Chrome Extension浏览器插件版
- [Mist Wallet](https://github.com/ethereum/mist/releases/latest) 官方版钱包
- [imToken](https://token.im/) 移动App版钱包

#### 以太坊开发环境
* geth安装 (以太坊客户端)
  * sudo apt-get install software-properties-common
  * sudo add-apt-repository -y ppa:ethereum/ethereum
  * sudo apt-get update
  * sudo apt-get install ethereum
* solc安装 (solidity是以太坊智能合约的开发语言，想要测试智能合约，开发DAPP的需要安装solc)
  * sudo add-apt-repository ppa:ethereum/ethereum
  * sudo apt-get update
  * sudo apt-get install solc
* truffle安装 (truffle是配套的以太坊开发框架。通过truffle可以快速的编译和部署合约并进行测试，同时还有web前端交互界面)
  * npm config set registry https://registry.npm.taobao.org 
  * npm install -g truffle
  * 运行命令不成功，可能是保存路径的问题 创建软链接(sudo ln -s /usr/nodejs/lib/node_modules/truffle/build/cli.bundled.js /usr/local/bin/truffle)
* testrpc安装 (testrpc可以理解为快速生成以太坊测试账号)
  * npm install -g ethereumjs-testrpc
  * npm install -g ganache-cli (testrpc 已经重命名为 ganache-cli)
  * 运行命令不成功，可能是保存路径的问题 创建软链接(sudo ln -s /usr/nodejs/lib/node_modules/ganache-cli/build/cli.node.js /usr/local/bin/ganache-cli)
* mist安装 (web3浏览器和钱包)
  * 安装依赖
    * Node.js v7.x 以上
    * Meteor (需要翻墙)
      * curl https://install.meteor.com/ | sh
    * Yarn (包管理)
      *  curl -o- -L https://yarnpkg.com/install.sh | bash
  * 源码安装
    * git clone https://github.com/ethereum/mist.git
    * cd mist
    * yarn
    * 更新用mist (这步只是升级时用到)
      * cd mist
      * git pull
      * yarn
  * 运行mist
    * yarn dev:meteor
    * cd mist
    * yarn dev:electron
* remix安装(ide)
  * git clone https://github.com/ethereum/remix-ide.git
  * cd remix-ide
  * npm install
  * npm run setupremix  
  * npm start
#### 以太坊私有网络
* 创建文件夹 mkdir private-geth
* 切换到文件夹下 cd private-geth
* 建立创世纪区块文件,是一个json格式的文件 vim genesis.json
  ```
   {
    "nonce": "0x0000000000000042",     
    "timestamp": "0x00",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "extraData": "0x00",     
    "gasLimit": "0xffffffff",     
    "difficulty": "0x400",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "coinbase": "0x0000000000000000000000000000000000000000",     
    "alloc": {
     },
     "config": {
        "chainId": 101,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
    }
   }
  ```
  * 参数说明
    * mixhash : 与nonce配合用于挖矿，由上一个区块的一部分生成的hash。注意他和nonce的设置需要满足以太坊的Yellow paper, 4.3.4. Block Header Validity, (44)章节所描述的条件
    * nonce : nonce就是一个64位随机数，用于挖矿，注意他和mixhash的设置需要满足以太坊的Yellow paper, 4.3.4. Block Header Validity, (44)章节所描述的条件
    * difficulty : 设置当前区块的难度，如果难度过大，cpu挖矿就很难，这里设置较小难度
    * alloc : 用来预置账号以及账号的以太币数量，因为私有链挖矿比较容易，所以我们不需要预置有币的账号，需要的时候自己创建即可以
    * coinbase : 矿工的账号，随便填
    * timestamp	: 设置创世块的时间戳
    * parentHash : 上一个区块的hash值，因为是创世块，所以这个值是0
    * extraData ： 附加信息，随便填，可以填你的个性信息
    * gasLimit ： 该值设置对GAS的消耗总量限制，用来限制区块能包含的交易信息总和，因为我们是私有链，所以填最大
    * config.chainId 以太坊区块链网络Id，ethereum主链是1，私有链只用不要与主链冲突即可
* 初始化创世纪节点,并设置data目录
  * geth --datadir ./data/00 init genesis.json
* 启动节点, 加上console 表示启动后,启用命令行:
  * geth --identity EthNode0 --datadir ./data/00 --rpc --rpcapi db,eth,net,web3 --rpcaddr 127.0.0.1 --ipcpath ./data/00/geth/geth.ipc --rpcport 8180 --port 30300 --networkid 20000 console
  * 使用命令 geth -h 可以查看geth 相关的帮助文档。这里几个常用的属性
  ```--identity : 节点身份标识，起个名字
     --datadir : 指定节点存在位置，“data0”
     --rpc : 启用http-rpc服务器
     --rpcapi : 基于http-rpc提供的api接口。eth,net,web3,db...
     --rpcaddr : http-rpc服务器接口地址：默认“127.0.0.1”
     --rpcport : http-rpc 端口(多节点时，不要重复)
     --port : 节点端口号（多节点时，不要重复）
     --networkid : 网络标识符 随便指定一个id（确保多节点是统一网络，保持一致）
  ```
  * 以太坊客户端JavaScript控制台常用命令
    * 创建账户 : personal.newAccount("123456")
    * 获取账户数组 : eth.accounts  
    * 解锁账户，转账时可使用 : personal.unlockAccount(eth.accounts[0], "123456")
    * 节点主账户 : eth.coinbase
    * 查看账户余额 : eth.getBalance(eth.accounts[0])
    * 启动，结束挖矿，写区块 : miner.start(),miner.stop() 
    * 通过区块号查看区块 : eth.getBlock(33)
    * 查看区块数 : eth.blockNumber
    * 通过交易hash查看交易 : eth.getTransaction("0x0c59f431068937cbe9e230483bc79f59bd7146edc8ff5ec37fea6710adcab825") 
    * 通过查看txpool : txpool.status 
  * 配置多节点
    * geth --datadir ./data/01 init genesis.json
    * geth --identity EthNode1 --datadir ./data/01 --rpc --rpcapi db,eth,net,web3 --rpcaddr 127.0.0.1 --ipcpath ./data/01/geth/geth.ipc --rpcport 8181 --port 30301 --networkid 20000 console
    * 查看：新节点enode信息，使用你本机ip替换[::]
      在第一个节点 输入： admin.nodeInfo.enode
      结果："enode://4311dc8d4aa85f7cc20794e49772c5b187020fd55860e1c4ce400d8eb4787224f84150ae57f79e04edbe8bb5cdfe7a17c38c06b4649fad2ffb5328265fe182d9@[::]:30300"
      在第二个节点链接第一节点：admin.addPeer("enode://4311dc8d4aa85f7cc20794e49772c5b187020fd55860e1c4ce400d8eb4787224f84150ae57f79e04edbe8bb5cdfe7a17c38c06b4649fad2ffb5328265fe182d9@127.0.0.1:30300")
    * 链接后会自动同步区块
#### 以太坊DAPP开发工具
* web3.js 以太坊RPC接口调用(不包含帐号创建和密钥等接口)
* lightwaller 以太坊轻钱包(帐号创建，交易，合约调用)
* truffle + testrpc/ganache-cli (dapp框架)