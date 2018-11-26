package main

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
