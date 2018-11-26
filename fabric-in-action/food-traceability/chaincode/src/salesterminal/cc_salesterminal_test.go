package main

import (
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var stub *shim.MockStub

//测试 Init 方法
//func testInit(t *testing.T, args [][]byte) {
//
//    res := stub.MockInit("1", args)
//    if res.Status != shim.OK {
//        fmt.Println("Init failed : ", string(res.Message))
//        t.FailNow()
//    }
//}

func testInvoke(t *testing.T, args [][]byte) {
	res := stub.MockInvoke("1", args)
	fmt.Printf("res:" + string(res.Payload) + "\n")

	if res.Status != shim.OK {
		fmt.Println("Invoke", "failed", string(res.Message))
		//t.FailNow()
	}

}

func TestDemo(t *testing.T) {

	scc := new(SalesterminalChaincode)
	stub = shim.NewMockStub("salesterminal", scc)

	testInvoke(t, [][]byte{[]byte("addSalesterminal"), []byte("S001"), []byte("销售终端1")})
	testInvoke(t, [][]byte{[]byte("addSalesterminal"), []byte("S002"), []byte("销售终端2")})
	testInvoke(t, [][]byte{[]byte("get"), []byte("S002")})

	testInvoke(t, [][]byte{[]byte("addMilk"), []byte(`{"id":"F0010000010000010001","time":"2018-10-08 15:48:27","InSaleTime":"2018-10-08 15:48:27","saleId":"S001","stutas":0}`)})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F0010000010000010001")})

	testInvoke(t, [][]byte{[]byte("addOperation"), []byte(`{"milkId":"F0010000010000010001","operation":1,"consumptionOrOutput":"售出"}`)})
	testInvoke(t, [][]byte{[]byte("addOperation"), []byte(`{"milkId":"F0010000010000010001","operation":2,"consumptionOrOutput":"下架"}`)})
	testInvoke(t, [][]byte{[]byte("addOperation"), []byte(`{"milkId":"F0010000010000010001","operation":0,"consumptionOrOutput":"上架"}`)})

	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2")})
	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2-1")})

}
