### IPFS私链搭建
星际文件系统（InterPlanetary File System，缩写IPFS）是一个旨在创建持久且分布式存储和共享文件的网络传输协议，它是一种内容可寻址的对等超媒体分发协议。在IPFS网络中的节点将构成一个分布式文件系统。

* 准备两台 centos7 虚拟机，设置并网络连通
  * 192.168.1.210 
  * 192.168.1.213
*  到 https://github.com/ipfs/go-ipfs 下在安装包 https://github.com/ipfs/go-ipfs/releases/tag/v0.4.18
   * go-ipfs_v0.4.18_linux-amd64.tar.gz
* 下载并安装密钥创建工具
   * go get -u github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen
* 将在 $GOPATH/bin/ipfs-swarm-key-gen 和  go-ipfs_v0.4.18_linux-amd64.tar.gz 复制到两台虚拟机里
* 两台虚拟机安装ipfs
  * 解压  
    ```
    tar -zxvf go-ipfs_v0.4.18_linux-amd64.tar.gz
    ```
  * 移动文件  
    ``` 
    cd go-ipfs && sudo mv ipfs /usr/local/bin/ipfs
    ```
* 初始化IPFS节点（无需在内网寻找相邻节点）
  ```
  ipfs init
  ```
* 在 192.168.1.210 虚拟机里创建共享密钥，同一个IPFS私链内的所有节点必须共享同一个密钥才能加入。  
  创建密钥： 
  ```
  ./ipfs-swarm-key-gen > /root/.ipfs/swarm.key
  ```
  创建完密钥放在了自己的ipfs默认配置文件夹下面（~/.ipfs/）  
  并清除所有缺省启动节点 
  ```
  ipfs bootstrap rm all
  ``` 
* 在 192.168.1.213 里清除所有缺省启动节点 
  ```
  ipfs bootstrap rm all
  ``` 
  将 192.168.1.210 创建的 /root/.ipfs/swarm.key 复制到  192.168.1.213 的 /root/.ipfs/swarm.key  
  在 192.168.1.213  添加 192.168.1.210 为默认节点 （节点地址通过 `ipfs id` 查看）
  ```   
  ipfs bootstrap add /ip4/192.168.1.210/tcp/4001/ipfs/QmRBWWrTwd7d1QCKEjcLdGgZAvSKtSNgZdQznzF58RBwZ2
  ```
* 修改配置
  在 192.168.1.210  添加
  ```
  ipfs config  Addresses.API "/ip4/192.168.1.210/tcp/5001"
  ipfs config  Addresses.Gateway "/ip4/192.168.1.210/tcp/8080"
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin "[\"*\"]"
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials "[\"true\"]"
  ```
  在 192.168.1.213  添加
  ```
  ipfs config  Addresses.API "/ip4/192.168.1.213/tcp/5001"
  ipfs config  Addresses.Gateway "/ip4/192.168.1.213/tcp/8080"
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin "[\"*\"]"
  ipfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials "[\"true\"]"
  ```

* 两台虚拟机启动IPFS节点 
  关闭防火墙
  ```
  iptables -F 
  ``` 
  ```
  ipfs daemon & 
  ```
  查看peer
  ```
  ipfs swarm peers 
  ```
  输出日志，成功
  ```
  /ip4/192.168.1.210/tcp/4001/ipfs/QmRBWWrTwd7d1QCKEjcLdGgZAvSKtSNgZdQznzF58RBwZ2
  ```
