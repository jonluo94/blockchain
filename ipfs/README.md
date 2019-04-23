#### IPFS
星际文件系统（InterPlanetary File System，缩写IPFS）是一个旨在创建持久且分布式存储和共享文件的网络传输协议，它是一种内容可寻址的对等超媒体分发协议。在IPFS网络中的节点将构成一个分布式文件系统。

- 源码安装(多种方式，详情请看 [go-ipfs](https://github.com/ipfs/go-ipfs)
  * 安装go，配置GOPATH
  * go get -u -d github.com/ipfs/go-ipfs
  * cd $GOPATH/src/github.com/ipfs/go-ipfs
  * make install
  * 查看安装成功 ipfs version
- 创建 IPFS 网络节点
  * ipfs init
- 常用命令
  * ipfs id (节点唯一id)
  * ipfs add file (添加file文件，返回文件hash)
  * ipfs add -r filedir (添加 filedir文件夹所有文件，通过http://127.0.0.1:8080/ipfs/hash/file 访问,其中hash为文件夹hash，file为文件夹下的文件名)
  * ipfs cat hash (通过hash查看文件内容，http://127.0.0.1:8080/ipfs/hash 本地查看)
  * ipfs daemon (启动服务器，同步数据，通过 https://ipfs.io/ipfs/hash 查询，需要翻墙)
- web端（http://localhost:5001/webui）
- 跨域资源共享CORS配置
  * ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methos '["PUT","GET","POST","OPTIONS"]'
  * ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
  * ipfs config --json API.HTTPHeaders.Access-Control-Allow-Headers '["Authorization"]'
  * ipfs config --json API.HTTPHeaders.Access-Control-Allow-Credentials '["true"]'
  * ipfs config --json API.HTTPHeaders.Access-Control-Expose-Headers '["Location"]'
- IPFS+IPNS静态个人博客小案例
  * 创建博客文件夹，编写静态文件
  * 文件夹 ipfs-blog 为静态blog代码 
  * ipfs add -r ipfs-blog/ 添加到本地节点
  * 因为访问固定不变的地址访问，需要ipns
  * 首先发布ipns ： ipfs name publish hash (hash为项目ipfs-blog文件夹的hash)
  * 就可以通 http://127.0.0.1:8080/ipns/id 访问blog (id 为 ipfs id 查看的id)
  * 每次跟新代码都要重新执行 ipfs name publish hash
- JS-IPFS-API 实现读写数据，文件上传下载 (公链)
  * 项目文件夹 jsipfs-app
  * npm install 安装依赖
  * 通过 npm run dev 运行
  * 详细 api 文档 [js-ipfs-api](https://github.com/ipfs/js-ipfs-api)
- JS-IPFS-HTTP-CLIENT 实现读写数据，文件上传下载 (私链)
  * 项目文件夹 priv-ipfs-fileservice