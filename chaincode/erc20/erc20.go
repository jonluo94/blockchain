/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"bytes"
	"encoding/pem"
	"crypto/x509"
)

const (
	TokenId      = "MyToken"
	TokenOwner   = TokenId + "-Owner"
	TokenBalance = TokenId + "-%s-Balance"
	TokenFreeze  = TokenId + "-%s-Freeze"
	TokenApprove = TokenId + "-%s-Approve-%s"
)

type TokenChaincode struct {
}

type ERC20Token struct {
	Name        string  `json:"name"`        //名称
	Symbol      string  `json:"symbol"`      //符号
	Decimals    uint8   `json:"decimals"`    //小数位
	TotalSupply float64 `json:"totalSupply"` //总数
}

func (t *TokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fcn, args := stub.GetFunctionAndParameters()
	fmt.Printf("方法: %s  参数 ： %s \n", fcn, args)
	if len(args) != 1 {
		return shim.Error("参数个数不是1")
	}
	tokenBts := []byte(args[0])

	var token ERC20Token
	err := json.Unmarshal(tokenBts, &token)
	if err != nil {
		return shim.Error(err.Error())
	}
	//检查
	err = checkToken(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	//添加代币
	err = stub.PutState(TokenId, tokenBts)
	if err != nil {
		return shim.Error(err.Error())
	}
	//创建人
	creator := initiator(stub)
	err = stub.PutState(TokenOwner, []byte(creator))
	if err != nil {
		return shim.Error(err.Error())
	}
	//拥有者初始化拥有所有代币
	err = stub.PutState(balanceKey(creator), parseDecimals(token.Decimals, token.TotalSupply))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	//获取要调用的方法名和方法参数
	fn, args := stub.GetFunctionAndParameters()
	fmt.Printf("方法: %s  参数 ： %s \n", fn, args)

	if fn == "tokenInfo" {
		return t.tokenInfo(stub)
	} else if fn == "balanceOf" {
		return t.balanceOf(stub, args)
	} else if fn == "minter" {
		return t.minter(stub, args)
	} else if fn == "transfer" {
		return t.transfer(stub, args)
	} else if fn == "freezeAccount" {
		return t.freezeAccount(stub, args)
	} else if fn == "approve" {
		return t.approve(stub, args)
	} else if fn == "transferFrom" {
		return t.transferFrom(stub, args)
	} else if fn == "allowance" {
		return t.allowance(stub, args)
	} else if fn == "transferOwner" {
		return t.transferOwner(stub, args)
	} else if fn == "increaseAllowance" {
		return t.increaseAllowance(stub, args)
	} else if fn == "decreaseAllowance" {
		return t.decreaseAllowance(stub, args)
	} else if fn == "burn" {
		return t.burn(stub, args)
	}

	return shim.Error("方法不存在")
}

//获取token信息
func (t *TokenChaincode) tokenInfo(stub shim.ChaincodeStubInterface) pb.Response {
	token, err := stub.GetState(TokenId)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(token)
}

//输入地址，可以获取该地址代币的余额
func (t *TokenChaincode) balanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不为1")
	}
	name := args[0]
	balance, err := stub.GetState(balanceKey(name))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(balance)
}

//挖矿
func (t *TokenChaincode) minter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不为2")
	}
	to := args[0]
	v, err := strconv.ParseFloat(args[1], 64)
	//v不能小于零
	if v < 0 {
		return shim.Error("v less than 0")
	}
	//检查是否是创建人
	err = checkOwner(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	//获取to的balance
	a, err := getBalance(stub, to)
	if err != nil {
		return shim.Error(err.Error())
	}
	a += v
	if a < 0 {
		return shim.Error("a balance less than 0")
	}

	//代币总数增加
	tks, err := stub.GetState(TokenId)
	if err != nil {
		return shim.Error(err.Error())
	}
	var token ERC20Token
	err = json.Unmarshal(tks, &token)
	if err != nil {
		return shim.Error(err.Error())
	}
	token.TotalSupply += v

	tks, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	//跟新代币
	err = stub.PutState(TokenId, tks)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 重新写回账本
	err = stub.PutState(balanceKey(to), parseDecimals(token.Decimals, a))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//调用transfer函数将自己的token转账给to地址，value为转账个数
func (t *TokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不为2")
	}
	//to
	to := args[0]
	//交易数量
	val := args[1]
	//from
	from := initiator(stub)
	//保留获取小数位
	decimals, err := getDecimals(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//交易
	err = deal(stub, from, to, val, decimals)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//实现账户的冻结和解冻 (true 冻结，false 解冻)
func (t *TokenChaincode) freezeAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不为2")
	}
	to := args[0]
	isFreeze := args[1]
	if isFreeze != "true" && isFreeze != "false" {
		return shim.Error("isFreeze is true or false")
	}
	//检查是否是创建人
	err := checkOwner(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//账户冻结和解冻
	err = stub.PutState(freezeKey(to), []byte(isFreeze))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//转移拥有者
func (t *TokenChaincode) transferOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不为1")
	}
	toOwner := args[0]
	//检查是否是创建人
	err := checkOwner(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//改变owner
	err = stub.PutState(TokenOwner, []byte(toOwner))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//批准spender账户从自己的账户转移value个token。可以分多次转移
func (t *TokenChaincode) approve(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不为2")
	}
	//授权人
	spender := args[0]
	val := args[1]
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return shim.Error("Invalid transaction amount")
	}
	//发起人
	sponsor := initiator(stub)
	//批准
	err = stub.PutState(approveKey(sponsor, spender), []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//与approve搭配使用，approve批准之后，调用transferFrom函数来转移token。
func (t *TokenChaincode) transferFrom(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("参数个数不为3")
	}
	//from
	sponsor := args[0]
	//to
	to := args[1]
	//val
	val := args[2]
	valf, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return shim.Error("Invalid transaction amount")
	}
	//spender
	spender := initiator(stub)
	//授权数量
	v, err := getApprove(stub, sponsor, spender)
	if err != nil {
		return shim.Error(err.Error())
	}
	//超出授权
	if valf > v {
		return shim.Error("approve not enough")
	}
	//保留获取小数位
	decimals, err := getDecimals(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//交易
	err = deal(stub, sponsor, to, val, decimals)
	if err != nil {
		return shim.Error(err.Error())
	}
	//计算approve剩余
	v -= valf
	//跟新授权数量
	err = stub.PutState(approveKey(sponsor, spender), parseDecimals(decimals, v))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//返回spender还能提取sponsor的token的个数
func (t *TokenChaincode) allowance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不为1")
	}
	//批准人
	sponsor := args[0]
	//发起人
	spender := initiator(stub)
	val, err := stub.GetState(approveKey(sponsor, spender))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(val)
}

//增加spender账户从自己的账户转移value个token。可以分多次转移
func (t *TokenChaincode) increaseAllowance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return changeAllowance(stub, args, "+")
}

//减少spender账户从自己的账户转移value个token。可以分多次转移
func (t *TokenChaincode) decreaseAllowance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return changeAllowance(stub, args, "-")
}

//改变spender账户从自己的账户转移value个token。可以分多次转移
func changeAllowance(stub shim.ChaincodeStubInterface, args []string, operation string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数不为2")
	}
	//授权人
	spender := args[0]
	v, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Invalid transaction amount")
	}
	//v不能小于零
	if v < 0 {
		return shim.Error("v less than 0")
	}
	//发起人
	sponsor := initiator(stub)
	//获取当前allowance
	val, err := stub.GetState(approveKey(sponsor, spender))
	if err != nil {
		return shim.Error(err.Error())
	}
	a, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	if operation == "+" {
		//增加
		a += v
	}
	if operation == "-" {
		//减少
		a -= v
	}
	//不能溢出
	if a < 0 {
		return shim.Error("a less than 0")
	}
	decimals, err := getDecimals(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//批准
	err = stub.PutState(approveKey(sponsor, spender), parseDecimals(decimals, a))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//销毁代币
func (t *TokenChaincode) burn(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不为1")
	}
	v, err := strconv.ParseFloat(args[0], 64)
	//v不能小于零
	if v < 0 {
		return shim.Error("v less than 0")
	}
	from := initiator(stub)
	//获取from的balance
	a, err := getBalance(stub, from)
	if err != nil {
		return shim.Error(err.Error())
	}
	a -= v
	if a < 0 {
		return shim.Error("a balance less than 0")
	}
	//代币总数减少
	tks, err := stub.GetState(TokenId)
	if err != nil {
		return shim.Error(err.Error())
	}
	var token ERC20Token
	err = json.Unmarshal(tks, &token)
	if err != nil {
		return shim.Error(err.Error())
	}
	token.TotalSupply -= v

	tks, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	//跟新代币
	err = stub.PutState(TokenId, tks)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 重新写回账本
	err = stub.PutState(balanceKey(from), parseDecimals(token.Decimals, a))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//交易处理 from to val decimals
func deal(stub shim.ChaincodeStubInterface, from, to, val string, decimals uint8) (error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	//v不能小于零
	if v < 0 {
		return fmt.Errorf("v less than 0")
	}
	//获取from的balance
	a, err := getBalance(stub, from)
	if err != nil {
		return err
	}
	//判断两个帐号不能相同
	if from == to {
		return fmt.Errorf("from and to is the same address")
	}
	//获取to的balance
	b, err := getBalance(stub, to)
	if err != nil {
		return err
	}
	//execution a b v 交易
	a -= v
	if a < 0 {
		return fmt.Errorf("from balance not enough")
	}
	b += v
	if b < 0 {
		return fmt.Errorf("to balance less than 0")
	}
	fmt.Println(" a: ", a, " b: ", b)
	// 重新写回账本
	err = stub.PutState(balanceKey(from), parseDecimals(decimals, a))
	if err != nil {
		return err
	}
	err = stub.PutState(balanceKey(to), parseDecimals(decimals, b))
	if err != nil {
		return err
	}
	return nil
}

//获取批准数量
func getApprove(stub shim.ChaincodeStubInterface, sponsor, spender string) (float64, error) {
	//批准数量
	val, err := stub.GetState(approveKey(sponsor, spender))
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		return 0, err
	}
	return b, nil
}

//获取balance
func getBalance(stub shim.ChaincodeStubInterface, addr string) (float64, error) {
	//检查帐号是否冻结
	isFreeze, err := stub.GetState(freezeKey(addr))
	if err != nil {
		return 0, err
	}
	if isFreeze != nil && bytes.Equal(isFreeze, []byte("true")) {
		return 0, fmt.Errorf("addr is freeze")
	}
	//查询balance
	addrVal, err := stub.GetState(balanceKey(addr))
	if err != nil {
		return 0, err
	}
	//为空返回0
	if addrVal == nil {
		return 0, nil
	}
	b, err := strconv.ParseFloat(string(addrVal), 64)
	if err != nil {
		return 0, err
	}
	return b, nil
}

//校验创建人
func checkOwner(stub shim.ChaincodeStubInterface) error {
	creator := initiator(stub)
	owner, err := stub.GetState(TokenOwner)
	if err != nil {
		return err
	}
	if !bytes.Equal([]byte(creator), owner) {
		return fmt.Errorf("is not owner")
	}
	return nil
}

//校验token
func checkToken(token ERC20Token) error {
	if token.Name == "" {
		return fmt.Errorf("name不能为空")
	}
	if token.Symbol == "" {
		return fmt.Errorf("symbol不能为空")
	}
	if token.TotalSupply <= 0 {
		return fmt.Errorf("totalSupply要大于0")
	}
	return nil
}

//转换为token decimals
func parseDecimals(decimals uint8, value float64) []byte {
	val := strconv.FormatFloat(value, 'f', int(decimals), 64)
	return []byte(val)
}

//获取token decimals
func getDecimals(stub shim.ChaincodeStubInterface) (uint8, error) {
	tokenBts, err := stub.GetState(TokenId)
	if err != nil {
		return 0, err
	}
	var token ERC20Token
	err = json.Unmarshal(tokenBts, &token)
	if err != nil {
		return 0, err
	}
	return token.Decimals, nil
}

//交易发起人
func initiator(stub shim.ChaincodeStubInterface) string {
	//获取当前用户
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Println("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Println("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Println("ParseCertificate failed")
	}
	name := cert.Subject.CommonName
	fmt.Println("initiator:" + name)
	return name
}

//封装balance key
func balanceKey(name string) string {
	return fmt.Sprintf(TokenBalance, name)
}

//封装freeze key
func freezeKey(name string) string {
	return fmt.Sprintf(TokenFreeze, name)
}

//封装approve key
func approveKey(from, to string) string {
	return fmt.Sprintf(TokenApprove, from, to)
}

func main() {
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		fmt.Printf("Error starting Token chaincode: %s", err)
	}
}
