package main

import (
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

var stub *shim.MockStub

func testInvoke(t *testing.T, args [][]byte) {
    res := stub.MockInvoke("1",args)
    if res.Status != shim.OK {
        fmt.Println("Invoke", "failed", string(res.Message))
        t.FailNow()
    }
}


func TestVote(t *testing.T) {

    scc := new(VoteChaincode)
    stub = shim.NewMockStub("vote", scc)

    testInvoke(t,[][]byte{[]byte("voteUser"), []byte("a")})

    testInvoke(t,[][]byte{[]byte("voteUser"), []byte("b")})

    testInvoke(t,[][]byte{[]byte("getUserVote"), []byte("")})

}