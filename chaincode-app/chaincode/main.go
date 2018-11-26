package main

import (
	"fmt"
	/*导入 chaincode shim 包和 peer protobuf 包*/
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

)

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
	if fn == "set" {
		return s.set(stub, args)
	}else if fn == "get"{
		return s.get(stub, args)
	}

	return shim.Error("方法不存在:"+fn)
}

func (s *HelloWorld) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{

	if len(args) != 2 {
		return shim.Error("要输入键和值")
	}
	//写入
	err := stub.PutState(args[0],[]byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	// 通知监听器，事件 eventInvokeSet 已被执行
	err = stub.SetEvent("eventInvokeSet", []byte{})
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

func main(){

	if err := shim.Start(new(HelloWorld)); err != nil {
		fmt.Println("HelloWorld start error")
	}
}