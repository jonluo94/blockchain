package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type VoteChaincode struct {

}

type Vote struct {
	Name string `json:"name"`
	VoteNum int `json:"votenum"`
}

func (t * VoteChaincode) Init (stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t * VoteChaincode) Invoke (stub shim.ChaincodeStubInterface) peer.Response {

	fn , args := stub.GetFunctionAndParameters()

	if fn == "voteUser" {
		return t.voteUser(stub, args)
	} else if fn == "getUserVote" {
		return t.getUserVote(stub)
	}

	return shim.Error("调用方法不存在！")
}

func (t *VoteChaincode) voteUser (stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	name := args[0]
	// 判断当前用户是否存在
	userAsBytes , err := stub.GetState(name)

	if err != nil {
		return shim.Error(err.Error())
	}

	vote := Vote{}

	if userAsBytes != nil {
		err = json.Unmarshal(userAsBytes , &vote)

		if err != nil {
			return shim.Error(err.Error())
		}

		vote.VoteNum += 1
	} else {
		vote = Vote{Name: name, VoteNum: 1}
	}

	voteJsonAsBytes , err := json.Marshal(vote)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(name, voteJsonAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	// 通知监听器，事件 eventInvokeVoteUser 已被执行，投票成功事件
	err = stub.SetEvent("eventInvokeVoteUser", []byte(name))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *VoteChaincode) getUserVote( stub shim.ChaincodeStubInterface) peer.Response {
	//获取投票列表
	resultIterator, err := stub.GetStateByRange("","")
	defer resultIterator.Close()
	if err != nil {
		return shim.Error(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")
	isFrist := true
	for resultIterator.HasNext() {
		queryResponse , err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
        //第一个不加，
		if !isFrist {
			buffer.WriteString(",")
		}
		isFrist = false
		buffer.WriteString(string(queryResponse.Value))
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

func main() {

	if err := shim.Start(new(VoteChaincode)); err != nil {
		fmt.Println("chaincode start error")
	}
}