package main

import (
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

var stub *shim.MockStub

func testInit(t *testing.T, args [][]byte) {
    res := stub.MockInit("1", args)
    if res.Status != shim.OK {
        fmt.Println("Init", args, "failed", string(res.Message))
        t.FailNow()
    }
}

func testInvoke(t *testing.T, args [][]byte) {
    res := stub.MockInvoke("1", args)
    if res.Status != shim.OK {
        fmt.Println("Invoke", args, "failed", string(res.Message))
        t.FailNow()
    }
}

func TestSimpleStub(t *testing.T) {

    simple := new(SimpleStub)
    stub = shim.NewMockStub("simple", simple)
    testInit(t,[][]byte{[]byte("a"),[]byte("100")})
    testInvoke(t, [][]byte{[]byte("set"),[]byte("a"),[]byte("90")})
}
