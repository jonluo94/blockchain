package main

import (
	"fmt"
	/*导入 chaincode shim 包和 peer protobuf 包*/
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"encoding/json"
	"time"
	"strconv"
	"bytes"
	"strings"
)

//参考: https://github.com/cloudant/mango

const prefix = "jonluo"

type CloudCertificateChaincode struct {
}

// 云证
type CloudCertificate struct {
	CloudCardNumber   string `json:"cloudCardNumber"`   //云证编号
	CloudCardPerson   string `json:"cloudCardPerson"`   //存证方
	CloudCardPlatform string `json:"cloudCardPlatform"` //传证平台
	Time              int64  `json:"time"`              //存证时间
	BlockNumber       string `json:"blockNumber"`       //存证区块号
	CloudCardHash     string `json:"cloudCardHash"`     //存证hash
	FileType          string `json:"fileType"`          //文件类型
	FileLabel         string `json:"fileLabel"`         //文件标签
	FileName          string `json:"fileName"`          //文件名
	FileAddress       string `json:"fileAddress"`       //下载地址
}

//初始化方法
func (s *CloudCertificateChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

//调用Chaincode
func (s *CloudCertificateChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	//获取要调用的方法名和方法参数
	fn, args := stub.GetFunctionAndParameters()

	fmt.Printf("方法: %s  参数 ： %s \n", fn, args)

	if fn == "addCard" {
		return s.addCard(stub, args)
	} else if fn == "getList" {
		return s.getList(stub, args)
	} else if fn == "get" {
		return s.get(stub, args)
	}

	return shim.Error("方法不存在")
}

func (s *CloudCertificateChaincode) addCard(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	cardStr := args[0]

	var card CloudCertificate
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(cardStr), &card); err != nil {
		return shim.Error("json反序列化失败")
	}

	t := time.Now()
	id := prefix + strconv.FormatInt(t.UnixNano(), 10)
	card.CloudCardNumber = id
	card.Time = t.Unix()

	bys, err := json.Marshal(card)
	fmt.Println("json:" + string(bys))

	if err != nil {
		return shim.Error("json序列化失败")
	}

	err = stub.PutState(id, bys)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *CloudCertificateChaincode) getList(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("要输入一个键")
	}

	page, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("page 出错")
	}
	size, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("size 出错")
	}
	index := (page - 1) * size

	pmap := map[string]string{}
	if err := json.Unmarshal([]byte(args[2]), &pmap); err != nil {
		return shim.Error("json反序列化失败")
	}
	//封装条件
	selector := selectionCriteria(pmap)
	fmt.Println(selector)

	queryIterator, err := stub.GetQueryResult(selector)
	defer queryIterator.Close()

	var list = make([]CloudCertificate, 0)

	if err != nil {
		return shim.Error("GetQueryResult 出错")
	} else {
		var next = 0

		for queryIterator.HasNext() {

			if next == page*size {
				break
			}

			if next >= index {

				item, err := queryIterator.Next()
				if err != nil {
					return shim.Error("queryIterator.Next 出错")
				}

				var c CloudCertificate
				err = json.Unmarshal(item.Value, &c)
				if err != nil {
					return shim.Error("json反序列化失败")
				}
				list = append(list, c)
			}

			next++

		}
	}

	msg, err := json.Marshal(list)
	fmt.Println("json:" + string(msg))

	if err != nil {
		return shim.Error("json序列化失败")
	}

	return shim.Success(msg)

}

func (s *CloudCertificateChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {

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

func selectionCriteria(pmap map[string]string) string {

	var buffer bytes.Buffer
	buffer.WriteString(`{"selector":{`)
	buffer.WriteString(`"cloudCardNumber":{"$regex": "^` + prefix + `.*"},`)

	for k, v := range pmap {

		switch k {

		case "startTime":
			if v != "" {
				buffer.WriteString(`"time":{"$gte": ` + v + `},`)
			}
		case "endTime":
			if v != "" {
				buffer.WriteString(`"time":{"$lte": ` + v + `},`)
			}
		case "startTime-endTime":
			if v != "" {
				args := strings.Split(v,"-")
				buffer.WriteString(`"time":{"$gte": ` + args[0] + `,"$lte": ` + args[1] + `},`)
			}

		case "fileType":
			if v != "" && v != "," {
				types := `"fileType":{"$or":[`

				args := strings.Split(v,",")
				for i,tyv := range args  {
					if i != 0 {
						types += `,`
					}
					types +=`{"$eq":"`+tyv+`"}`
				}
				types += `]},`

				buffer.WriteString(types)
			}

		default:
			if k != "" && v != "" {
				buffer.WriteString(`"` + k + `":{"$eq": "` + v + `"},`)
			}

		}
	}
	buffer.Truncate(buffer.Len()-1)

	buffer.WriteString("}}")

	return buffer.String()
}

func main() {

	if err := shim.Start(new(CloudCertificateChaincode)); err != nil {
		fmt.Println("CloudCertificateChaincode start error")
	}
}
