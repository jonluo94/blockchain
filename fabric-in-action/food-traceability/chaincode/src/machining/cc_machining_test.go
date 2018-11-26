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
//	res := stub.MockInit("1", args)
//	if res.Status != shim.OK {
//		fmt.Println("Init failed : ", string(res.Message))
//		t.FailNow()
//	}
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

	scc := new(MachiningChaincode)
	stub = shim.NewMockStub("machining", scc)

	testInvoke(t, [][]byte{[]byte("addMachining"), []byte("FA001"), []byte("加工厂1")})
	testInvoke(t, [][]byte{[]byte("addMachining"), []byte("FA002"), []byte("加工厂2")})

	testInvoke(t, [][]byte{[]byte("get"), []byte("M002")})

	testInvoke(t, [][]byte{[]byte("addBucket"), []byte(`{"id":"F001000001000001","machiningId":"FA001","time":"2018-10-08 15:26:37","inMachiningTime":"2018-10-08 15:26:37","stutas":0}`)})
	testInvoke(t, [][]byte{[]byte("addBucket"), []byte(`{"id":"F001000001000001","machiningId":"FA0012","time":"2018-10-08 15:26:37","inMachiningTime":"2018-10-08 15:26:37","stutas":0}`)})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F001000001000001")})

	testInvoke(t, [][]byte{[]byte("addMilkOperation"), []byte(`{"bucketId":"F001000001000001","operation":1,"consumptionOrOutput":"灌装"}`)})
	testInvoke(t, [][]byte{[]byte("addMilkOperation"), []byte(`{"bucketId":"F001000001000001","operation":0,"consumptionOrOutput":"消毒"}`)})

	testInvoke(t, [][]byte{[]byte("addMilkPack"), []byte(`{"bucketId":"F001000001000001","consumptionOrOutput":"包装"}`)})
	testInvoke(t, [][]byte{[]byte("addMilkPack"), []byte(`{"bucketId":"F001000001000001","consumptionOrOutput":"包装"}`)})

	testInvoke(t, [][]byte{[]byte("sentSale"), []byte("F0010000010000010001"), []byte("S001")})
	testInvoke(t, [][]byte{[]byte("sentSale"), []byte("F0010000010000010002"), []byte("S001")})

	testInvoke(t, [][]byte{[]byte("checkMilkForSaleterminal"), []byte("S001")})

	testInvoke(t, [][]byte{[]byte("confirmMilk"), []byte("F0010000010000010002"), []byte("S001"), []byte("2")})

	testInvoke(t, [][]byte{[]byte("checkMilkForSaleterminal"), []byte("S001")})

	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2")})
	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2-1")})

}
