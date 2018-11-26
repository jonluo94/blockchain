## fabric-go-sdk 创建第一个app
1. 创建config.yaml，根据hlf-network配置信息
2. 编写blockchain文件夹的三个go文件
   * setup.go -- sdk配置包含创建通道，安装和实例化链码
   * query.go -- 链码查询方法的调用
   * invoke.go -- 写入通道的方法的调用
3. 编写main方法
4. 通过使用dep工具在vendor目录中展平这些依赖项来处理依赖冲突
   * 新建 Gopkg.toml 文件，写入：
     ``` 
     [[constraint]]
       name = "github.com/hyperledger/fabric-sdk-go"
       revision = "614551a752802488988921a730b172dada7def1d"
     ```
   * 执行： dep ensure （需要翻墙）
   * 会生成 vendor 和 Gopkg.lock
5. 编写web端集成到web下，并在main方法启动web服务
