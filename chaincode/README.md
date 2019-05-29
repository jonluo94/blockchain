#### Chaincode学习例子
* hello world 
* simple stub (shim.ChaincodeStubInterface 的方法的使用)
* vote (投票案例)
* couchdbquery (couchdb复杂查询案例)
* erc20 (chaincode实现erc20标准的代币案例)
#### 以 hello world 实现操作，其他 Chaincode 同理实现即可
* 建 helloworld 文件夹，并在文件夹下新建 helloworld.go 和 helloworld_test.go
* 编写链上代码 helloworld.go
  * 关键依赖
    ```
    import (
        "fmt"
        /*导入 chaincode shim 包和 peer protobuf 包*/
        "github.com/hyperledger/fabric/core/chaincode/shim"
        "github.com/hyperledger/fabric/protos/peer"
    )
    ```
  * 实现了 Chaincode 接口 Init 方法和 Invoke 方法
    ```
    // helloworld 结构体
    type HelloWorld struct {
    
    }
    //结构体中的初始化方法
    func (s *HelloWorld) Init(stub shim.ChaincodeStubInterface) peer.Response {
  
    	//接受 string 数组
        args:= stub.GetStringArgs()
        
        if len(args) != 2 {
        	return shim.Error("要输入键和值")
        }
        //转为 k-v 写入区块链
        err := stub.PutState(args[0],[]byte(args[1]))
        
        if err != nil {
        	shim.Error(err.Error())
        }
        
        return shim.Success(nil)
    } 
    //调用Chaincode
    func (s *HelloWorld) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
    
    	//获取要调用的方法名和方法参数
    	fn, args := stub.GetFunctionAndParameters()
        //根据方法名调用方法（set，get）
    	if fn =="set" {
    		return s.set(stub, args)
    	}else if fn == "get"{
    		return s.get(stub, args)
    	}
    
    	return shim.Error("方法不存在")
    }
    ```
  * set 和 get 业务方法实现
    ```
    func (s *HelloWorld) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
    
    	if len(args) != 2 {
    		return shim.Error("要输入键和值")
    	}
    	//写入
    	err := stub.PutState(args[0],[]byte(args[1]))
    	if err != nil {
    		return shim.Error(err.Error())
    	}
    	return shim.Success(nil)
    }
    
    func (s *HelloWorld) get(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    
    	if len(args) != 1 {
    		return shim.Error("要输入一个键")
    	}
        //读出
    	value, err := stub.GetState(args[0])
    
    	if err != nil {
    		return shim.Error(err.Error())
    	}
    
    	return shim.Success(value)
    }
    ```
  * 最后 main 方法，启动实例
    ```
    func main(){
    
    	if err := shim.Start(new(HelloWorld)); err != nil {
    		fmt.Println("HelloWorld start error")
    	}
    } 
    ```
* 编写链上代码测试用例 helloworld_test.go
  * 利用 chaincode shim包 MockStub 的 MockInit 方法和 MockInvoke 方法 进行模拟测试
    ```
    import (
        "fmt"
        "testing"
        "github.com/hyperledger/fabric/core/chaincode/shim"
    )
    
    var stub *shim.MockStub
    
    //测试 Init 方法
    func testInit(t *testing.T, args [][]byte) {
    
        res := stub.MockInit("1", args)
        if res.Status != shim.OK {
            fmt.Println("Init failed : ", string(res.Message))
            t.FailNow()
        }
    }
    
    //测试 set 方法
    func testSet(t *testing.T, key string, value string) {
    
        res := stub.MockInvoke("1", [][]byte{[]byte("set"), []byte(key),[]byte(value)})
        if res.Status != shim.OK {
            fmt.Println("set", key, " failed : ", string(res.Message))
            t.FailNow()
        }
    
    }
    
    //测试 get 方法
    func testGet(t *testing.T, key string) {
    
        res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte(key)})
        if res.Status != shim.OK {
            fmt.Println("get", key, "failed", string(res.Message))
            t.FailNow()
        }
        if res.Payload == nil {
            fmt.Println("get", key, "failed to get value")
            t.FailNow()
        }
    
        fmt.Println("get value", key, " : ", string(res.Payload))
    
    }
    
    //测试
    func TestHelloWorld(t *testing.T) {
    
        //模拟实例
        stub = shim.NewMockStub("helloworld", new(HelloWorld))
    
        testInit(t, [][]byte{[]byte("hi"), []byte("jonluo")})
        testGet(t, "hi")
        testSet(t, "say","helloworld")
        testGet(t,"say")
    }
    ```
  * 在helloworld文件夹下执行 : go test -v helloworld_test.go helloworld.go
    结果：
    ```
    === RUN   TestHelloWorld
    get value hi  :  jonluo
    get value say  :  helloworld
    --- PASS: TestHelloWorld (0.00s)
    PASS
    ok      command-line-arguments  0.021s
    ```
* 在开发链上部署测试链上代码
  * 利用 hyperledger/fabric-samples 进行部署（安装详情请看 hyperledger-fabric 目录）
  * 将链码文件夹 helloworld 复制到 fabric-samples/chaincode 下
  * 到 fabric-samples 下：cd chaincode-docker-devmode
  * 在 chaincode-docker-devmode 下打开三个终端
  * 终端1-启动网络
    ```
    docker-compose -f docker-compose-simple.yaml up
    ```
  * 终端2-构建并启动链码
    进入容器
    ```
    docker exec -it chaincode bash
    ```
    到链码目录编译代码
    ```
    cd helloworld
    go build
    ```
    运行链码
    ```
    CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./helloworld 
    ```
    日志
    ```
    2018-08-25 03:58:31.237 UTC [shim] SetupChaincodeLogging -> INFO 001 Chaincode log level not provided; defaulting to: INFO
    2018-08-25 03:58:31.237 UTC [shim] SetupChaincodeLogging -> INFO 002 Chaincode (build level: ) starting up ...
    ```
  * 终端3-使用链码
    进入容器
    ``` 
    docker exec -it cli bash
    ```
    节点安装链码
    ```
    peer chaincode install -p chaincodedev/chaincode/helloworld -n mycc -v 0
    ```
    节点实例化链码
    ```
    peer chaincode instantiate -n mycc -v 0 -c '{"Args":["hi","jonluo"]}' -C myc 
    ```
    验证
    ```
    peer chaincode query -n mycc -c '{"Args":["get","hi"]}' -C myc 
    peer chaincode invoke -n mycc -c '{"Args":["set", "hello", "world"]}' -C myc
    peer chaincode query -n mycc -c '{"Args":["get","hello"]}' -C myc 
    ```
    完成

#### 相关文档 
- [chaincode shim API文档](https://godoc.org/github.com/hyperledger/fabric/core/chaincode/shim)
- [protos peer API文档](https://godoc.org/github.com/hyperledger/fabric/protos/peer)

