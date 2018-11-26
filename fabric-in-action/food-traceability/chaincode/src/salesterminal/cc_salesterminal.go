package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"time"
	"encoding/json"
	"strings"
)

const (
	channelName            = "mychannel" //通道名
	machiningChaincodeName = "machining"         // 加工厂chaincode名
	dairyfarmChaincodeName = "dairyfarm"         // 奶牛场chaincode名
	operationSuffix        = "-opr"                 //奶牛操作后缀
	dateFomat              = "2006-01-02 15:04:05"  //日期格式
)

type MilkHistory struct {
	MilkInfo map[string]interface{} `json:"milkInfo"`
	BucketInfo map[string]interface{} `json:"bucketInfo"`
	CowInfo map[string]interface{} `json:"cowInfo"`
	SaleHistory []map[string]interface{} `json:"saleHistory"`
	MachHistory []map[string]interface{} `json:"machHistory"`
	DairHistory []map[string]interface{} `json:"dairHistory"`
}

type SalesterminalChaincode struct {
}

//牛奶
type Milk struct {
	Id         string `json:"id"`         //批次编号
	Time       string `json:"time"`       //生产日期
	InSaleTime string `json:"inSaleTime"` //进入售买日期
	SaledTime  string `json:"saledTime"`  //售出日期
	SaleId     string `json:"saleId"`     // 销售终端id
	Stutas     int8   `json:"stutas"`     //状态
}

//销售终端
type Salesterminal struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`       //名称
	CreateTime string   `json:"createTime"` //创建时间
	MaxMilkNo  uint     `json:"maxMilkNo"`  //最大牛奶编号
	MilkIds    []string `json:"milkIds"`    //牛奶批次
}

//加工操作
type SaleOperation struct {
	MilkId              string `json:"milkId"`              //装满奶牛的桶id
	Time                string `json:"time"`                //时间
	Operation           int8   `json:"operation"`           // 操作类型 0为上架，1为售出，2为下架
	ConsumptionOrOutput string `json:"consumptionOrOutput"` //消耗或产出
}



func (t *SalesterminalChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SalesterminalChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	fn, args := stub.GetFunctionAndParameters()

	fmt.Printf("方法: %s  参数 ： %s \n", fn, args)

	if fn == "addSalesterminal" {
		return t.addSalesterminal(stub, args)
	} else if fn == "addMilk" {
		return t.addMilk(stub, args)
	} else if fn == "addOperation" {
		return t.addOperation(stub, args)
	} else if fn == "getOperationHistory" {
		return t.getOperationHistory(stub, args)
	} else if fn == "getMilkHistory" {
		return t.getMilkHistory(stub, args)
	} else if fn == "get" {
		return t.get(stub, args)
	} else if fn == "set" {
		return t.set(stub, args)
	}

	return shim.Error("Salesterminal No operation:" + fn)

}

//添加销售终端
func (t *SalesterminalChaincode) addSalesterminal(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("参数出错")
	}

	id := args[0]
	name := args[1]

	if isExisted(stub, id) {
		return shim.Error("已存在")
	}

	salesterminal := Salesterminal{}
	salesterminal.Id = id
	salesterminal.Name = name
	salesterminal.CreateTime = time.Now().Format(dateFomat)
	salesterminal.MaxMilkNo = 0
	salesterminal.MilkIds = []string{}

	jsonString, err := json.Marshal(salesterminal)

	fmt.Println("json:" + string(jsonString))

	if err != nil {
		return shim.Error("json序列化失败")
	}

	err = stub.PutState(id, jsonString)

	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(jsonString)

}

//根据id获取加销售终端
func (t *SalesterminalChaincode) getSalesterminalById(stub shim.ChaincodeStubInterface, id string) (Salesterminal,error) {

	dvalue, err := stub.GetState(id)
	var d Salesterminal
	err = json.Unmarshal([]byte(dvalue), &d)
	return d,err

}

//添加牛奶
func (t *SalesterminalChaincode) addMilk(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	milkStr := args[0]

	var milk Milk
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(milkStr), &milk); err != nil {
		return shim.Error("json反序列化失败")
	}

	salesterminalId := milk.SaleId
	salesterminal,err:= t.getSalesterminalById(stub, salesterminalId)
	if err != nil {
		return shim.Error("销售终端不存在")
	}

	salesterminal.MaxMilkNo += 1

	milkIds := salesterminal.MilkIds
	salesterminal.MilkIds = append(milkIds, milk.Id)

	//跟新销售
	jsonString, err := json.Marshal(salesterminal)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	err = stub.PutState(salesterminalId, []byte(jsonString))
	if err != nil {
		shim.Error(err.Error())
	}

	fmt.Println("json:" + string(jsonString))
	//添加牛奶
	err = stub.PutState(milk.Id, []byte(milkStr))
	if err != nil {
		shim.Error(err.Error())
	}

	fmt.Println("json:" + string(milkStr))
	return shim.Success([]byte(milkStr))

}

//添加奶桶操作
func (t *SalesterminalChaincode) addOperation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	operationStr := args[0]

	var operation SaleOperation
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(operationStr), &operation); err != nil {
		return shim.Error("json反序列化失败")
	}
	operation.Time = time.Now().Format(dateFomat)
	milkId := operation.MilkId

	operationJson, err := json.Marshal(operation)
	if err != nil {
		return shim.Error("json序列化失败")
	}

	fmt.Println("json:" + string(operationJson))
	err = stub.PutState(milkId+operationSuffix, []byte(operationJson))
	if err != nil {
		return shim.Error(err.Error())
	}
	//若售出
	if operation.Operation == 1 {

		milkStr, err := stub.GetState(milkId)
		if err != nil {
			return shim.Error(err.Error())
		}
		var milk Milk
		//这里就是实际的解码和相关的错误检查
		if err := json.Unmarshal([]byte(milkStr), &milk); err != nil {
			return shim.Error("json反序列化失败")
		}

		milk.SaledTime = time.Now().Format(dateFomat)

		//跟新牛奶
		jsonString, err := json.Marshal(milk)
		if err != nil {
			return shim.Error("json序列化失败")
		}
		err = stub.PutState(milkId, []byte(jsonString))
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println("json:" + string(jsonString))

	}

	return shim.Success(nil)

}

func (t *SalesterminalChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	keys := strings.Split(args[0],",")

	var dmap = make(map[string]map[string]interface{})

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		//读出
		value, err := stub.GetState(key)
		if err != nil {
			return shim.Error(err.Error())
		}

		if value == nil {
			dmap[key] = nil
		} else {
			var imap map[string]interface{}
			//这里就是实际的解码和相关的错误检查
			if err := json.Unmarshal([]byte(value), &imap); err != nil {
				return shim.Error("json反序列化失败")
			}
			dmap[key] = imap
		}

	}

	jsonStr, err := json.Marshal(dmap)
	if err != nil {
		return shim.Error("json序列化失败")
	}

	fmt.Println("json:", string(jsonStr))
	return shim.Success(jsonStr)
}

func (t *SalesterminalChaincode) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("要输入一个键")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//获取操作记录
func (t *SalesterminalChaincode) getOperationHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("parm error")
	}
	id := args[0]
	keysIter, err := stub.GetHistoryForKey(id + operationSuffix)

	if err != nil {
		return shim.Error(fmt.Sprintf("GetHistoryForKey failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var historys = make([]SaleOperation, 0)

	for keysIter.HasNext() {

		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return shim.Error(fmt.Sprintf("GetHistoryForKey operation failed. Error accessing state: %s", err))
		}
		//交易的值
		txvalue := response.Value

		var operation SaleOperation
		//这里就是实际的解码和相关的错误检查
		if err := json.Unmarshal(txvalue, &operation); err != nil {
			return shim.Error("json反序列化失败")
		}

		historys = append(historys, operation)

	}

	jsonKeys, err := json.Marshal(historys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
	}

	fmt.Println("json:", string(jsonKeys))
	return shim.Success(jsonKeys)
}

//获取历史记录
func (t *SalesterminalChaincode) getMilkHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parm error")
	}
	milkId := args[0]
	//获取桶编号
	bucketId := string(milkId[0:16])
	//获取奶牛编号
	cowId := string(milkId[0:10])

	keysIter, err := stub.GetHistoryForKey(milkId + operationSuffix)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetHistoryForKey failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var milkHisrtoy = MilkHistory{}

	var saleHistory = make([]map[string]interface{}, 0)

	for keysIter.HasNext() {

		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return shim.Error(fmt.Sprintf("GetHistoryForKey operation failed. Error accessing state: %s", err))
		}
		//交易的值
		txvalue := response.Value

		var history map[string]interface{}
		//这里就是实际的解码和相关的错误检查
		if err := json.Unmarshal([]byte(txvalue), &history); err != nil {
			return shim.Error("json反序列化失败")
		}

		saleHistory = append(saleHistory, history)

	}

	milkHisrtoy.SaleHistory = saleHistory
	//获取牛奶信息
	milkStr, err := stub.GetState(milkId)
	if err != nil {
		return shim.Error(err.Error())
	}
	var milkInfo map[string]interface{}
	if err := json.Unmarshal(milkStr, &milkInfo); err != nil {
		return shim.Error("json反序列化失败")
	}
	milkHisrtoy.MilkInfo = milkInfo


	//获取加工厂的溯源
	response := stub.InvokeChaincode(machiningChaincodeName, [][]byte{[]byte("getOperationHistory"), []byte(bucketId)}, channelName)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}

	result := string(response.Payload)
	fmt.Println("桶操作：",result)

	var machiningHistorys []map[string]interface{}

	if err := json.Unmarshal([]byte(result), &machiningHistorys); err != nil {
		return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
	}

	milkHisrtoy.MachHistory = machiningHistorys

	//获取桶信息
	response = stub.InvokeChaincode(machiningChaincodeName, [][]byte{[]byte("get"), []byte(bucketId)}, channelName)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}

	result = string(response.Payload)

	fmt.Println("桶：",result)

	var bucketInfo map[string]interface{}
	if err := json.Unmarshal([]byte(result), &bucketInfo); err != nil {
		return shim.Error("json反序列化失败")
	}
	milkHisrtoy.BucketInfo = bucketInfo[bucketId].(map[string]interface{})


	//调用奶牛场的溯源信息
	response = stub.InvokeChaincode(dairyfarmChaincodeName, [][]byte{[]byte("getOperationHistory"), []byte(cowId)}, channelName)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}

	result = string(response.Payload)
	fmt.Println("奶牛操作：",result)
	var dairHistorys []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &dairHistorys); err != nil {
		return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
	}

	milkHisrtoy.DairHistory = dairHistorys

	//获取奶牛信息
	response = stub.InvokeChaincode(dairyfarmChaincodeName, [][]byte{[]byte("get"), []byte(cowId)}, channelName)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}

	result = string(response.Payload)
	fmt.Println("奶牛：",result)

	var cowInfo map[string]interface{}
	if err := json.Unmarshal([]byte(result), &cowInfo); err != nil {
		return shim.Error("json反序列化失败")
	}
	milkHisrtoy.CowInfo = cowInfo[cowId].(map[string]interface{})

	jsons, err := json.Marshal(milkHisrtoy)
	if err != nil {
		return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
	}

	return shim.Success(jsons)
}

func isExisted(stub shim.ChaincodeStubInterface, key string) bool {
	val, err := stub.GetState(key)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	if len(val) == 0 {
		return false
	}

	return true
}

func main() {
	err := shim.Start(new(SalesterminalChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
