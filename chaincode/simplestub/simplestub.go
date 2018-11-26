package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleStub struct {
}

func (t *SimpleStub) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	stub.PutState(args[0], []byte(args[1]))
	return shim.Success(nil)
}

func (t *SimpleStub) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	//GetFunctionAndParameters
	//获取方法名和参数
	//{"Args":["set","a","100"]}
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("GetFunctionAndParameters方法获取方法名和参数：", fn, args)
	stub.PutState(args[0], []byte(args[1]))

	//GetStringArgs
	//获取所有参数字符串数组，包含方法名
	args = stub.GetStringArgs()
	fmt.Println("GetStringArgs方法获取所有参数字符串数组，包含方法名：", args)

	//PutState
	//设置key-value 4个（str0-hello0，str1-hello1，str2-hello2,str3-hello3）
	stub.PutState("str0", []byte("hello0"))
	stub.PutState("str1", []byte("hello1"))
	stub.PutState("str2", []byte("hello2"))
	stub.PutState("str3", []byte("hello3"))
	fmt.Println("PutState方法设置4个值：（str0：hello0，str1：hello1，str2：hello2，str3：hello3）")

	//GetStateByRange
	//获取从key为str0到key为str2的值，返回k-v的迭代器
	resultIterator, _ := stub.GetStateByRange("str0", "str2")
	defer resultIterator.Close()
	fmt.Println("GetStateByRange获取从key为str0到key为str2的值，返回k-v的迭代器，遍历：")
	for resultIterator.HasNext() {
		item, _ := resultIterator.Next()
		fmt.Println(string(item.Key), string(item.Value))
	}

	//GetState
	//获取对应key的值
	a, _ := stub.GetState("a")
	fmt.Println("GetState方法获取a的值:", string(a))

	//GetHistoryForKey
	//获取key的历史值，返回历史迭代器，可以获取历史所在交易id和值
	//GetHistoryForKey 请求节点配置 core.ledger.history.enableHistoryDatabase 为 true
	historyIterator, err := stub.GetHistoryForKey("a")
	if err != nil {
		fmt.Println(err)
	} else {
		defer historyIterator.Close()
		for historyIterator.HasNext() {
			item, _ := historyIterator.Next()
			fmt.Println(item.TxId, string(item.Value))
		}
	}

	//DelState
	//删除key为a的键值
	stub.DelState("a")
	fmt.Printf("DeDelState方法删除key为a的键值，")
	a, _ = stub.GetState("a")
	fmt.Println("删除后获取a的值:", string(a))

	//CreateCompositeKey
	//创建组合键
	indexName := "sex~name"
	indexKey, _ := stub.CreateCompositeKey(indexName, []string{"boy", "jon"})
	fmt.Println("CreateCompositeKey方法创建组合键:", indexKey)

	stub.PutState(indexKey, []byte("0"))

	indexKey, _ = stub.CreateCompositeKey(indexName, []string{"boy", "luo"})
	fmt.Println("CreateCompositeKey方法创建组合键:", indexKey)
	stub.PutState(indexKey, []byte("0"))

	indexKey, _ = stub.CreateCompositeKey(indexName, []string{"girl", "wen"})
	fmt.Println("CreateCompositeKey方法创建组合键:", indexKey)
	stub.PutState(indexKey, []byte("0"))

	//GetStateByPartialCompositeKey
	//获取组合键的集合迭代器
	resultIterator, _ = stub.GetStateByPartialCompositeKey(indexName, []string{"boy"})
	defer resultIterator.Close()
	fmt.Println("GetStateByPartialCompositeKey方法获取有boy的集合迭代器，遍历：")
	for resultIterator.HasNext() {
		item, _ := resultIterator.Next()
		fmt.Println("key: " + item.Key)
		fmt.Println("value: " + string(item.Value))
		objectType, compositeKeyParts, _ := stub.SplitCompositeKey(item.Key)
		fmt.Println("objectType: " + objectType)
		fmt.Println("sex : " + compositeKeyParts[0])
		fmt.Println("name : " + compositeKeyParts[1])

	}

	//GetQueryResult
	//https://github.com/cloudant/mango
	//只支持支持丰富查询的状态数据库CouchDB
	//queryIterator, err1 := stub.GetQueryResult(`{"selector": {"sex": "boy"}}`)
	//if err1 != nil {
	//	fmt.Println(err1)
	//} else {
	//	defer queryIterator.Close()
	//	for queryIterator.HasNext() {
	//		item, _ := queryIterator.Next()
	//		fmt.Println(item.Key, string(item.Value))
	//	}
	//}

	//InvokeChaincode
	//InvokeChaincode方法在本地调用指定的chaincode的方法
	//同一通道的链码调用另一链码会影响另一链码状态
	//不同通道的链码调用另一链码不会影响被调用的链码状态，相当于一个查询
	/*
	  trans:=[][]byte{[]byte("invoke"),[]byte("a"),[]byte("b"),[]byte("11")}
	  result := stub.InvokeChaincode("mycc",trans,"myc")
	  fmt.Println(result)
	*/

	//获取当前用户
	//creatorByte, _ := stub.GetCreator()
	//certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	//if certStart == -1 {
	//	fmt.Errorf("No certificate found")
	//}
	//certText := creatorByte[certStart:]
	//bl, _ := pem.Decode(certText)
	//if bl == nil {
	//	fmt.Errorf("Could not decode the PEM structure")
	//}
	//
	//cert, err := x509.ParseCertificate(bl.Bytes)
	//if err != nil {
	//	fmt.Errorf("ParseCertificate failed")
	//}
	//uname := cert.Subject.CommonName
	//fmt.Println("Name:" + uname)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleStub))
	if err != nil {
		fmt.Println("start err")
	}
}
