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

	scc := new(DairyfarmChaincode)
	stub = shim.NewMockStub("dairyfarm", scc)

	testInvoke(t, [][]byte{[]byte("addDairyFarm"), []byte("F001"), []byte("奶牛场1")})
	testInvoke(t, [][]byte{[]byte("addDairyFarm"), []byte("F002"), []byte("奶牛场2")})
	testInvoke(t, [][]byte{[]byte("addDairyFarm"), []byte("F003"), []byte("奶牛场3")})

	testInvoke(t, [][]byte{[]byte("addCow"), []byte(`{"farmId":"F001","healthy":true,"quarantine":true,"feedSource":"国产","stutas":0}`)})
	testInvoke(t, [][]byte{[]byte("addCow"), []byte(`{"farmId":"F001","healthy":true,"quarantine":true,"feedSource":"国产","stutas":0}`)})

	testInvoke(t, [][]byte{[]byte("addCow"), []byte(`{"farmId":"F004","healthy":true,"quarantine":true,"feedSource":"国产","stutas":0}`)})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F001000001")})

	testInvoke(t, [][]byte{[]byte("addCowOperate"), []byte(`{"cowId":"F001000001","operation":1,"consumptionOrOutput":"食物1"}`)})
	testInvoke(t, [][]byte{[]byte("addCowMilking"), []byte(`{"cowId":"F001000001","consumptionOrOutput":"b1"}`)})
	testInvoke(t, [][]byte{[]byte("addCowMilking"), []byte(`{"cowId":"F001000001","consumptionOrOutput":"b12"}`)})

	testInvoke(t, [][]byte{[]byte("sentProcess"), []byte("F001000001000001"), []byte("M001")})
	testInvoke(t, [][]byte{[]byte("sentProcess"), []byte("F001000001000002"), []byte("M001")})

	testInvoke(t, [][]byte{[]byte("checkBucketForMachining"), []byte("M001")})

	testInvoke(t, [][]byte{[]byte("confirmBucket"), []byte("F001000001000002"), []byte("M001"), []byte("2")})

	testInvoke(t, [][]byte{[]byte("checkBucketForMachining"), []byte("M001")})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F001000001000002")})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F001000001-opr")})

	testInvoke(t, [][]byte{[]byte("delCow"), []byte("F001000001")})

	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2")})
	//testInvoke(t,[][]byte{[]byte("getOperationHistory"), []byte("F001_2-1")})

	testInvoke(t, [][]byte{[]byte("get"), []byte("F0010000012222")})

}
