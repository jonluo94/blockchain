### fabric配置讲解
在 fabric 源码 sampleconfig 目录的两个配置文件，core.yaml，orderer.yaml
* Peer 配置剖析 core.yaml
  ``` 
  ###############################################################################
  #
  #    日志部分
  #
  ###############################################################################
  logging:    ##  日志级别有： critical、 error、 warning、 notice、 info、 debug  级别由大到小， 级别越小输出越详细
      level:       info       ##  全局的默认日志级别
   
      ##  各个模块日志级别， 覆盖全局配置
      cauthdsl:   warning       
      gossip:     warning
      grpc:       error
      ledger:     info
      msp:        warning
      policies:   warning
      peer:
          gossip: warning
   
      # 日志的输出格式
      format: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'
   
  ###############################################################################
  #
  #    peer部分
  #
  ###############################################################################
  peer:
      id: jdoe            ##  Peer节点ID
      networkId: dev      ##  网络ID
      listenAddress: 0.0.0.0:7051     ##  节点监听的本地网络接口地址
      chaincodeAddress: 0.0.0.0:7052  ##  链码容器连接时的监听地址 如未指定, 则使用listenAddress的IP地址和7052端口
      address: 0.0.0.0:7051           ##  节点对外的服务地址 （对外的地址）
      addressAutoDetect: false        ##  是否自动探测服务地址 (默认 关闭， 如果启用TLS时，最好关闭)
      gomaxprocs: -1                  ##  Go的进程限制数 runtime.GOMAXPROCS(n) 默认 -1
      keepalive:                      ##  客户端和peer间的网络心跳连接配置
         
          minInterval: 60s            ##  最小的心跳间隔时间
          client: ##  该节点和客户端 的交互配置
              interval: 60s   ##  和客户端 的心跳间隔 必须 interval >= minInterval
              timeout: 20s    ##  和客户端 间的网络连接超时时间
          deliveryClient: ## 交付客户端用于与订购节点通信的心跳
              interval: 60s
              timeout: 20s
      gossip:     ##  节点间通信的gossip 协议的P2P通信 【主要包含了 启动 及 连接】
          bootstrap: 127.0.0.1:7051   ##  启动节点后 向哪些节点发起gossip连接，以加入网络，且节点相互间都是同一个组织的
          useLeaderElection: true     ##  是否启动动态选举 组织的Leader 节点 与 orgLeader 互斥
          orgLeader: false            ##  是否指定本节点为 组织Leader 节点 与 useLeaderElection 互斥
          endpoint:                   ##  本节点在组织内的gossip id
          maxBlockCountToStore: 100   ##  保存到内存的区块个数上限
          maxPropagationBurstLatency: 10ms    ##  保存消息的最大时间，超过则触发转发给其他节点
          maxPropagationBurstSize: 10         ##  保存的最大消息个数，超过则触发转发给其他节点
          propagateIterations: 1              ##  消息转发的次数
          propagatePeerNum: 3         ##  推送消息给指定个数的节点
          pullInterval: 4s            ##  拉取消息的时间间隔  (second) 必须大于 digestWaitTime + responseWaitTime
          pullPeerNum: 3              ##  从指定个数的节点拉取消息
          requestStateInfoInterval: 4s        ##  从节点拉取状态信息(StateInfo) 消息间隔 (second)
          publishStateInfoInterval: 4s        ##  向其他节点推动状态信息消息的间隔 (second)
          stateInfoRetentionInterval:         ##  状态信息消息的超时时间 (second)
          publishCertPeriod: 10s      ##  启动后在心跳消息中包括证书的等待时间
          skipBlockVerification: false        ##  是否不对区块消息进行校验，默认为false
          dialTimeout: 3s             ##  gRPC 连接拨号的超时 (second)
          connTimeout: 2s             ##  建立连接的超时 (second)
          recvBuffSize: 20            ##  收取消息的缓冲大小
          sendBuffSize: 200           ##  发送消息的缓冲大小
          digestWaitTime: 1s          ##  处理摘要数据的等待时间 (second)  可以大于 requestWaitTime
          requestWaitTime: 1500ms     ##  处理nonce 数据的等待时间 (milliseconds) 可以大于 digestWaitTime
          responseWaitTime: 2s        ##  终止拉取数据处理的等待时间 (second)
          aliveTimeInterval: 5s       ##  定期发送Alive 心跳消息的时间间隔 (second)
          aliveExpirationTimeout: 25s ##  Alive 心跳消息的超时时间 (second)
          reconnectInterval: 25s      ##  断线后重连的时间间隔 (second)
          externalEndpoint:           ##  节点被组织外节点感知时的地址，公布给其他组织的地址和端口, 如果不指定, 其他组织将无法知道本peer的存在
          election:   ## Leader 节点的选举配置
              startupGracePeriod: 15s         ##  leader节点选举等待的时间 (second)
              membershipSampleInterval: 1s    ##  测试peer稳定性的时间间隔 (second)
              leaderAliveThreshold: 10s       ##  pear 尝试进行选举的等待超时 (second)
              leaderElectionDuration: 5s      ##  pear 宣布自己为Leader节点的等待时间 (second)
   
          pvtData:        
              pullRetryThreshold: 60s
              transientstoreMaxBlockRetention: 1000
              pushAckTimeout: 3s
              btlPullMargin: 10
              reconcileBatchSize: 10
              reconcileSleepInterval: 5m
   
      
      events:     ##  事件配置 （在1.3版本后配置去除）
          address: 0.0.0.0:7053       ##  本地服务监听地址 (默认在所有网络接口上进心监听，端口 7053)
          buffersize: 100             ##  最大缓冲消息数，超过则向缓冲中发送事件消息会被阻塞
          timeout: 10ms               ##  事件发送超时时间, 如果事件缓存已满, timeout < 0, 事件被丢弃; timeout > 0, 阻塞直到超时丢弃, timeout = 0, 阻塞直到发送出去
          timewindow: 15m             ##  允许peer和 客户端 时间不一致的最大时间差
   
          keepalive:      ##  客户端到peer间的事件心跳
              minInterval: 60s
          sendTimeout: 60s            ##  在GRPC流上向客户端发送事件的超时时间
   
   
      tls:        ##  tls配置
          
          enabled:  false             ##  是否开启 TLS，默认不开启TLS
          clientAuthRequired: false   ##  客户端连接到peer是否需要使用加密
          
          cert:   ##  证书密钥的位置, 各peer应该填写各自相应的路径
              file: tls/server.crt    ##  本服务的身份验证证书，公开可见，访问者通过该证书进行验证
          key:    
              file: tls/server.key    ##  本服务的签名私钥
          rootcert:   
              file: tls/ca.crt        ##  信任的根CA整数位置
          
          clientRootCAs:              ##  用于验证客户端证书的根证书颁发机构的集合
              files:
                - tls/ca.crt
          clientKey:      ##  当TLS密钥用于制作客户端连接。如果不设置，将使用而不是peer.tls.key.file
              file:
          clientCert:     ##  在进行客户端连接时用于TLS的X.509证书。 如果未设置，将使用peer.tls.cert.file
              file:
   
          serverhostoverride:         ##  是否制定进行TLS握手时的主机名称
      
      authentication:     ##  身份验证包含与验证客户端消息相关的配置参数
          
          timewindow: 15m         ##  客户端请求消息中指定的当前服务器时间与客户端时间之间的可接受差异
   
      fileSystemPath: /var/hyperledger/production     ##  peer数据存储位置(包括账本,状态数据库等)
   
     
      BCCSP:      ##  加密库配置 与Orderer 配置一样
          Default: SW             ##  使用软件加密方式 (默认 SW)
          SW:     
              Hash: SHA2          ##  Hash 算法类型，目前仅支持SHA2
              Security: 256       
            
              FileKeyStore:       ##  本地私钥文件路径，默认指向 <mspConfigPath>/keystore
                  KeyStore:
          # Settings for the PKCS#11 crypto provider (i.e. when DEFAULT: PKCS11)
          PKCS11:     ##  设置 PKCS#11 加密算法 (默认PKCS11)
              Library:            ##  本地PKCS11依赖库  
   
              Label:              ##  token的标识
              Pin:                ##  使用Pin
              Hash:
              Security:
              FileKeyStore:
                  KeyStore:
   
      mspConfigPath: msp          ##  msp 的本地路径
   
      localMspId: SampleOrg       ##  Peer 所关联的MSP 的ID
   
      client:        ##   cli 公共客户端配置选项
          connTimeout: 3s     ##  连接超时时间
   
      
      deliveryclient:     ## 交付服务配置
          reconnectTotalTimeThreshold: 3600s  ##  交付服务交付失败后尝试重连的时间
          connTimeout: 3s     ##  交付服务和 orderer节点的连接超时时间
          reConnectBackoffThreshold: 3600s    ##  设置连续重试之间的最大延迟
   
      localMspType: bccsp     ##  本地MSP类型 （默认为 BCCSP）
   
      profile:            ##  是否启用Go自带的profiling 支持进行调试
          enabled:     false
          listenAddress: 0.0.0.0:6060
   
      adminService:       ##  admin服务用于管理操作，例如控制日志模块严重性等。只有对等管理员才能使用该服务
         
      handlers:
          authFilters:
            -
              name: DefaultAuth
            -
              name: ExpirationCheck       ##  此筛选器检查身份x509证书过期 
          decorators:
            -
              name: DefaultDecorator
          endorsers:
            escc:
              name: DefaultEndorsement
              library:
          validators:
            vscc:
              name: DefaultValidation
              library:
      validatorPoolSize:                  ##  处理交易验证的并发数, 默认是CPU的核数
   
      discovery:      ##  客户端使用发现服务来查询有关peers的信息，例如 - 哪些peer已加入某个channel，最新的channel配置是什么，最重要的是 - 给定chaincode和channel，哪些可能的peer满足认可 policy
          enabled: true
          authCacheEnabled: true
          authCacheMaxSize: 1000
          authCachePurgeRetentionRatio: 0.75
          orgMembersAllowedAccess: false
  ###############################################################################
  #
  #     vm环境配置，目前主要支持 Docker容器
  #
  ###############################################################################
  vm:
   
      endpoint: unix:///var/run/docker.sock   ##  Docker Daemon 地址，默认是本地 套接字
   
      docker:
          tls:    ##  Docker Daemon 启用TLS时的相关证书配置, 包括信任的根CA证书、服务身份证书、签名私钥等等
              enabled: false
              ca:
                  file: docker/ca.crt
              cert:
                  file: docker/tls.crt
              key:
                  file: docker/tls.key
   
          attachStdout: false     ##  是否启用绑定到标准输出，启用后 链码容器 的输出消息会绑定到标准输出，方便进行调试
   
          hostConfig:             ##  Docker 相关的主机配置，包括网络配置、日志、内存等等，这些配置在启动链码容器时进行使用
              NetworkMode: host
              Dns:
                 # - 192.168.0.1
              LogConfig:
                  Type: json-file
                  Config:
                      max-size: "50m"
                      max-file: "5"
              Memory: 2147483648
   
  ###############################################################################
  #
  #    链码相关配置
  #
  ###############################################################################
  chaincode:
   
      id:             ##  记录链码相关信息，包括路径、名称、版本等等，该信息会以标签形式写到链码容器
          path:
          name:
   
      builder: $(DOCKER_NS)/fabric-ccenv:latest       ##  通用的本地编译环境，是一个Docker 镜像
      pull: false     ##      
      golang:         ##  Go语言的链码部署生成镜像的基础Docker镜像
          runtime: $(BASE_DOCKER_NS)/fabric-baseos:$(ARCH)-$(BASE_VERSION)
          dynamicLink: false
      car:            ##  car格式的链码部署生成镜像的基础Docker 镜像
          runtime: $(BASE_DOCKER_NS)/fabric-baseos:$(ARCH)-$(BASE_VERSION)
      java:           ##  java语言的基础镜像
          Dockerfile:  |
              from $(DOCKER_NS)/fabric-javaenv:$(ARCH)-1.1.0
      node:           ##  node.js的基础镜像
          runtime: $(BASE_DOCKER_NS)/fabric-baseimage:$(ARCH)-$(BASE_VERSION)
   
      startuptimeout: 300s    ##  启动链码容器超时，等待超时时间后还没收到链码段的注册消息，则认为启动失败
   
      executetimeout: 30s     ##  invoke 和 initialize 命令执行超时时间
   
      deploytimeout:          ##  部署链码的命令执行超时时间
      mode: net               ##  执行链码的模式，dev: 允许本地直接运行链码，方便调试； net: 意味着在容器中运行链码
      keepalive: 0            ##  Peer 和链码之间的心跳超市时间， <= 0 意味着关闭
      system:                 ##  系统链码的相关配置 (系统链码白名单 ??)
          cscc: enable
          lscc: enable
          escc: enable
          vscc: enable
          qscc: enable
      systemPlugins:          ##  系统链码插件
      logging:                ##  链码容器日志相关配置
        level:  info
        shim:   warning
        format: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'
   
  ###############################################################################
  #
  #    账本相关配置
  #
  ###############################################################################
  ledger:
    blockchain:       
    state:            ##  状态DB的相关配置(包括 golevelDB、couchDB)、DN连接、查询最大返回记录数等
      stateDatabase: goleveldb    ##  stateDB的底层DB配置  (默认golevelDB)
      couchDBConfig:              ##  如果启用couchdb，配置连接信息 (goleveldb 不需要配置这些)
         couchDBAddress: 127.0.0.1:5984
         username:
         password:
         maxRetries: 3    ##  运行时出错重试数
         maxRetriesOnStartup: 10  ##  启动时出错的重试数
         requestTimeout: 35s      ##  请求超时时间
         queryLimit: 10000        ##  每个查询最大返回数
         maxBatchUpdateSize: 1000 ##  批量更新最大记录数
         warmIndexesAfterNBlocks: 1
    history:      
      enableHistoryDatabase: true    ##  是否启用历史数据库，默认开启
   
  ###############################################################################
  #
  #    服务度量监控配置
  #
  ###############################################################################
  metrics:
          enabled: false      ##  是否开启监控服务
          reporter: statsd
          interval: 1s
          statsdReporter:
                address: 0.0.0.0:8125
                flushInterval: 2s
                flushBytes: 1432
          promReporter:       ##  prometheus 普罗米修斯服务监听地址
  ``` 
* Orderer 配置剖析 orderer.yaml
  ```
  ################################################################################
  #
  #   Orderer的配置
  #
  ################################################################################
  General:
      LedgerType: file            ##  账本类型，支持ram、json、file 三种类型【建议用file】，其中ram保存在内存中；json、file保存在本地文件中 (通常为 /var/hyperledger/production/orderer 下)
      ListenAddress: 127.0.0.1    ##  服务监听地址，一般需要制定为服务的特定网络接口地址 或者全网(0.0.0.0)
      ListenPort: 7050            ##  服务监听端口 默认7050
   
      
      TLS:        ##  启用TLS 时的相关配置 (grpc 传输)
          Enabled: false    
          PrivateKey: tls/server.key      ##  Orderer 签名私钥
          Certificate: tls/server.crt     ##  Orderer 身份证书
          RootCAs:
            - tls/ca.crt      ##  根证书
          ClientAuthRequired: false       ##  是否对客户端也进行认证
          ClientRootCAs:      
   
      Keepalive:      ##  设置GRPC 服务心跳检查
          ServerMinInterval: 60s          ##  客户端和 orderer 的 最小心跳间隔
          ServerInterval: 7200s           ##  客户端和 orderer 的心跳间隔时间
          ServerTimeout: 20s              ##  客户端和 Orderer 的超时时间
   
      LogLevel: info      ##  日志等级
   
      ##  日志输出格式
      LogFormat: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'
   
      GenesisMethod: provisional          ##  创世块的提供方式 (系统通道初始区块的提供方式，支持 provisional 或者 file；前者根据GenesisProfile 指定默认的 $FABRIC_CFG_PATH/config.yaml 文件中的profile生成；后者使用GenesisFile 指定现成的初始区块文件)
      
      GenesisProfile: SampleInsecureSolo  ##  创世块使用的Profile；GenesisMethod: provisional 才有效
   
      GenesisFile: genesisblock           ##  使用现成创世块文件时，文件的路径 [创世块的位置]  GenesisMethod: file 才有效
   
      LocalMSPDir: msp                    ##  本地MSP文件的路径 【orderer节点所需的安全认证文件的位置】
      LocalMSPID: SampleOrg               ##  Orderer所关联的MSP的ID  MSP管理器用于注册安全认证文件的ID, 此ID必须与配置系统通道和创世区块时(configtx.yaml的OrdererGenesis部分)指定的组织中的某一个组织的ID一致
      
      Profile:        ##  为Go pprof性能优化工具启用一个HTTP服务以便作性能分析(https://golang.org/pkg/net/http/pprof)
          Enabled: false                  ##  不启用
          Address: 0.0.0.0:6060           ##  Go pprof的HTTP服务监听的地址和端口
   
      BCCSP:      ##  加密库配置  具体参照Peer 配置
          Default: SW
          SW:
              Hash: SHA2
              Security: 256
              FileKeyStore:
                  KeyStore:
      Authentication:
          TimeWindow: 15m
   
  ################################################################################
  #
  #   基于文件账本配置 (file和json两种类型)
  #
  ################################################################################
  FileLedger:
      Location: /var/hyperledger/production/orderer       ##  指定存放文件的位置，一般为 /var/hyperledger/production/orderer, 该目录下的 chains目录存放各个chain的区块，index目录存放 索引文件 (如果这项不指定, 每次节点重启都将使用一个新的临时位置) 
      Prefix: hyperledger-fabric-ordererledger            ##  如果不指定Location，则在临时目录下创建账本时目录的名称
   
  ################################################################################
  #
  #   基于内存账本配置 
  #
  ################################################################################
  RAMLedger:
      HistorySize: 1000           ##  内存账本所支持存储的区块的数量, 如果内存中存储的区块达到上限, 继续追加区块会导致最旧的区块被丢弃
   
  ################################################################################
  #
  #   kafka 集群配置
  #
  ################################################################################
  Kafka:
  # kafka是一种基于发布/订阅模式的分布式消息系统
  # fabric网络中, orderer节点集群组成kafka集群, 客户端是kafka集群的Producer(消息生产者), peer是kafka集群的Consumer(消息消费者)
  # kafka集群使用ZooKeeper(分布式应用协调服务)管理集群节点, 选举leader.
      Retry:      ##  连接时的充实操作 kafka 会利用 sarama 客户端为chennel创建一个producer 负责向kafka 写数据，一个comsumer负责kafka读数据
          ShortInterval: 5s           ##  操作失败后的快速重试间隔时间
          ShortTotal: 10m             ##  快速重试阶段最对重试多久
          LongInterval: 5m            ##  快速充实阶段仍然失败后进入 慢重试阶段，慢重试的时间间隔
          LongTotal: 12h              ##  慢重试最多重试多久
         
          # https://godoc.org/github.com/Shopify/sarama#Config
          NetworkTimeouts:            ##  Sarama 网络超时时间
              DialTimeout: 10s         
              ReadTimeout: 10s
              WriteTimeout: 10s
          Metadata:                   ##  kafka集群leader 选举中的metadata 请求参数
              RetryBackoff: 250ms     ##  leader选举过程中元数据请求失败的重试间隔
              RetryMax: 3             ##  最大重试次数
          Producer:                   ##  发送消息到kafka集群时的超时
              RetryBackoff: 100ms     ##  向kafka集群发送消息失败后的重试间隔
              RetryMax: 3             ##  最大重试次数
          Consumer:                   ##  从kafka集群接收消息时的超时
              RetryBackoff: 2s        ##  从kafka集群拉取消息失败后的重试间隔
      Verbose: false                  ##  是否开启kafka的客户端的调试日志 (orderer与kafka集群交互是否生成日)
   
      TLS:        ##  与kafka集群的连接启用TLS时的相关配置
        Enabled: false                ##  是否开启TLS，默认不开启
        PrivateKey:                   ##  Orderer 身份签名私钥
          # File:                       ##    私钥文件路径 
        Certificate:                  ##  kafka的身份证书
          # File:                       ##    证书文件路径 
        RootCAs:                      ##  验证kafka证书时的根证书
          # File:                       ##    根证书文件路径 
      Version:                        ##  kafka的版本
   
  ################################################################################
  #
  #   Orderer节点的调试参数
  #
  ################################################################################
  Debug:
      BroadcastTraceDir:      ##  该orderer节点的广播服务请求保存的位置
      DeliverTraceDir:        ##  该orderer节点的传递服务请求保存的位置
   
  ##  以下配置是1.4最新的配置
  ################################################################################
  #
  #   操作配置
  #
  ################################################################################
  Operations:
      # 操作服务地址端口
      ListenAddress: 127.0.0.1:8443
  
      # TLS 配置
      TLS:
          Enabled: false
          Certificate:
          PrivateKey:
          ClientAuthRequired: false
          RootCAs: []
  
  ################################################################################
  #
  #   度量配置
  #
  ################################################################################
  Metrics:
      # 度量提供程序是 prometheus或disabled
      Provider: disabled
  
      # statsd 配置
      Statsd:
        # 网络协议
        Network: udp
  
        # 服务地址
        Address: 127.0.0.1:8125
  
        # 将本地缓存的计数器和仪表推送到 statsd 的时间间隔;时间被立即推送
        WriteInterval: 30s
  
        # 前缀预先添加到所有发出的 statsd 指标
        Prefix:
  
  ################################################################################
  #
  #   共识配置
  #
  ################################################################################
  Consensus:
      # 这里允许的键值对取决于共识插件。对于 etcd/raft，
      # 我们使用以下选项:
      # WALDir 指定存储 etcd/raft 的Write Ahead Logs的位置。每个channel都有自己的以channelID命名的子目录
      WALDir: /var/hyperledger/production/orderer/etcdraft/wal
  
      # SnapDir 指定存储 etcd/raft 快照的位置。每个channel都有自己的以channelID命名的子目录
      SnapDir: /var/hyperledger/production/orderer/etcdraft/snapshot
  ```
  
搭建网络核心crypto-config.yaml 配置
* crypto-config.yaml 文件
  ```
  # ---------------------------------------------------------------------------
  # "OrdererOrgs"
  # ---------------------------------------------------------------------------
  OrdererOrgs:  ##  定义Orderer组织
    - Name: Orderer  ##  名称
      Domain: example.com  ##  组织的命名域
      # ---------------------------------------------------------------------------
      # "Specs" - 有关完整说明，请参阅下面的PeerOrgs
      # ---------------------------------------------------------------------------
      Specs:
        - Hostname: orderer
  # ---------------------------------------------------------------------------
  # "PeerOrgs" 
  # ---------------------------------------------------------------------------
  PeerOrgs:
    - Name: Org1  ##  名称
      Domain: org1.example.com  ##  组织的命名域
      EnableNodeOUs: true         ##  如果设置了EnableNodeOUs，就在msp下生成config.yaml文件
      Template:                   ##  允许定义从模板顺序创建的1个或多个主机。 默认情况下，这看起来像是从0到Count-1的“peer”。 您可以覆盖节点数（Count），起始索引（Start）或用于构造名称的模板（Hostname）。
        Count: 1                  ##  表示生成几个Peer
        # Start: 5
        # Hostname: {{.Prefix}}{{.Index}} # default
      Users:
        Count: 1  ##  表示生成普通User数量
        
    - Name: Org2  
      Domain: org2.example.com  
      EnableNodeOUs: true
      Template:                  
        Count: 1                 
      Users:
        Count: 1 
  ```
  
通道及锚节点的配置 configtx.yaml 配置剖析
* configtx.yaml 文件
  ```
  ################################################################################
  #
  #    Organizations部分
  #   【注意】：本文件中 &KEY 均为  *KEY 所引用；  xx：&KEY 均为  <<: *KEY 所引用
  ################################################################################
  Organizations:
    ##  定义Orderer组织  
    - &OrdererOrg
         Name: OrdererOrg        ##  Orderer的组织的名称
         ID: OrdererMSP          ##  Orderer 组织的ID （ID是引用组织的关键）
         MSPDir: crypto-config/ordererOrganizations/example.com/msp       ##  Orderer的 MSP 证书目录路径
         AdminPrincipal: Role.ADMIN ##  【可选项】 组织管理员所需要的身份，可选项: Role.ADMIN 和 Role.MEMBER 
 
    ##  定义Peer组织1
    - &Org1
         Name: Org1MSP           ##  组织名称 
         ID: Org1MSP             ##  组织ID
         MSPDir: crypto-config/peerOrganizations/org1.example.com/msp    ##  Peer的MSP 证书目录路径
         AnchorPeers:            ##  定义组织锚节点 用于跨组织 Gossip 通信
            - Host: peer0.org1.example.com      ##  锚节点的主机名
              Port: 7051                        ##  锚节点的端口号
    ##  定义Peer组织 2
    - &Org2
         Name: Org2MSP
         ID: Org2MSP
         MSPDir: crypto-config/peerOrganizations/org2.example.com/msp
         AnchorPeers:
            - Host: peer0.org2.example.com
              Port: 7051
  ################################################################################
  #   本节定义了 fabric 网络的功能. 
  ################################################################################
  Capabilities:
      ## 通道功能适用于orderers and the peers，并且必须得到两者的支持。 将功能的值设置为true.
      Global: &ChannelCapabilities
          ## V1.1 的 Global是一个行为标记，已被确定为运行v1.0.x的所有orderers和peers的行为，但其修改会导致不兼容。 用户应将此标志设置为true.
          V1_1: true
   
      ## Orderer功能仅适用于orderers，可以安全地操纵，而无需担心升级peers。 将功能的值设置为true
      Orderer: &OrdererCapabilities
          ## Orderer 的V1.1是行为的一个标记，已经确定为运行v1.0.x的所有orderers 都需要，但其修改会导致不兼容。 用户应将此标志设置为true
          V1_1: true
   
      ## 应用程序功能仅适用于Peer 网络，可以安全地操作，而无需担心升级或更新orderers。 将功能的值设置为true
      Application: &ApplicationCapabilities
          ## V1.2 for Application是一个行为标记，已被确定为运行v1.0.x的所有peers所需的行为，但其修改会导致不兼容。 用户应将此标志设置为true
          V1_2: true
  ################################################################################
  #
  #   应用通道相关配置，主要包括 参与应用网络的可用组织信息
  #
  ################################################################################
  Application: &ApplicationDefaults   ##  自定义被引用的地址
      Organizations:              ##  加入通道的组织信息
  ################################################################################
  #
  #   Orderer 系统通道相关配置，包括 Orderer 服务配置和参与Orderer 服务的可用组织
  #   Orderer 默认是 solo 的 且不包含任何组织 【主要被 Profiles 部分引用】
  ################################################################################
  Orderer: &OrdererDefaults   ##  自定义被引用的地址
      OrdererType: solo       ##  Orderer 类型，包含 solo 和 kafka 集群
      Addresses:              ##  服务地址
          - orderer.example.com:7050
      BatchTimeout: 2s        ##  区块打包的最大超时时间 (到了该时间就打包区块)
      BatchSize:              ##  区块打包的最大包含交易数
          MaxMessageCount: 10         ##  一个区块里最大的交易数
          AbsoluteMaxBytes: 98 MB     ##  一个区块的最大字节数， 任何时候都不能超过
          PreferredMaxBytes: 512 KB   ##  一个区块的建议字节数，如果一个交易消息的大小超过了这个值, 就会被放入另外一个更大的区块中
   
      MaxChannels: 0          ##  【可选项】 表示Orderer 允许的最大通道数， 默认 0 表示没有最大通道数
      Kafka:
          Brokers:                    ##  kafka的 brokens 服务地址 允许有多个
              - 127.0.0.1:9092
      Organizations:          ##  参与维护 Orderer 的组织，默认为空
  ################################################################################
  #
  #   Profile 
  #
  #   - 一系列通道配置模板，包括Orderer 系统通道模板 和 应用通道类型模板
  #
  ################################################################################
  Profiles:
      ##  Orderer的 系统通道模板 必须包括 Orderer、 Consortiums 两部分
      TwoOrgsOrdererGenesis:              ##  Orderer 系统的通道及创世块配置。通道为默认配置，添加一个OrdererOrg 组织， 联盟为默认的 SampleConsortium 联盟，添加了两个组织 【该名称可以自定义 ？？】
          Capabilities:
              <<: *ChannelCapabilities
          Orderer:    ##  指定Orderer系统通道自身的配置信息
              <<: *OrdererDefaults        ##  引用 Orderer 部分的配置  &OrdererDefaults
              Organizations:
                  - *OrdererOrg           ##  属于Orderer 的通道组织  该处引用了 【 &OrdererOrg 】位置内容
              Capabilities:
                  <<: *OrdererCapabilities
   
          Consortiums:    ##  Orderer 所服务的联盟列表
              SampleConsortium:           ##  创建更多应用通道时的联盟 引用 TwoOrgsChannel 所示
                  Organizations:
                      - *Org1
                      - *Org2
      ##  应用通道模板 必须包括 Application、  Consortium 两部分              
      TwoOrgsChannel:                     ##  应用通道配置。默认配置的应用通道，添加了两个组织。联盟为SampleConsortium
          Consortium: SampleConsortium    ##  通道所关联的联盟名称
          Application:    ##  指定属于某应用通道的信息，主要包括 属于通道的组织信息
              <<: *ApplicationDefaults
              Organizations:              ##  初始 加入应用通道的组织
                  - *Org1
                  - *Org2                 
              Capabilities:
                  <<: *ApplicationCapabilities
  ```