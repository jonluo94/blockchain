package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"time"
	"encoding/json"
	"strconv"
	"strings"
)

type MachiningChaincode struct {
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

//加工厂
type Machining struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`        //名称
	CreateTime  string   `json:"createTime"`  //创建时间
	MaxMilkNo   uint     `json:"maxMilkNo"`   //最大牛奶编号
	MaxBucketNo int      `json:"maxBucketNo"` //最大桶
	BucketIds   []string `json:"bucketIds"`   //状态
	MilkIds     []string `json:"milkIds"`     //牛奶批次
}

//奶牛桶
type Bucket struct {
	Id              string `json:"id"`
	MachiningId     string `json:"machiningId"`     //加工场id
	Time            string `json:"time"`            //装桶时间
	InMachiningTime string `json:"inMachiningTime"` //进入加工场时间
	Stutas          int8   `json:"stutas"`          //状态
}

//加工操作
type ProcessOperation struct {
	BucketId            string `json:"bucketId"`            //装满奶牛的桶id
	Time                string `json:"time"`                //时间
	Operation           int8   `json:"operation"`           // 操作类型 0为消毒，1为灌装，2为包装
	ConsumptionOrOutput string `json:"consumptionOrOutput"` //消耗或产出
}

const (
	channelName               = "mychannel" //通道名
	saletenminalChaincodeName = "salesterminal"         // 销售终端chaincode名
	operationSuffix           = "-opr"                 //奶牛操作后缀
	intPrefix                 = "%04d"                 //4位数值，不足前面补0
	dateFomat                 = "2006-01-02 15:04:05"  //日期格式
	combinedConstruction      = "saleId~milkId"        //组合键
)

func (t *MachiningChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (t *MachiningChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	fn, args := stub.GetFunctionAndParameters()

	fmt.Printf("方法: %s  参数 ： %s \n", fn, args)

	if fn == "addMachining" {
		return t.addMachining(stub, args)
	} else if fn == "addBucket" {
		return t.addBucket(stub, args)
	} else if fn == "addMilkOperation" {
		return t.addMilkOperation(stub, args)
	} else if fn == "addMilkPack" {
		return t.addMilkPack(stub, args)
	} else if fn == "sentSale" {
		return t.sentSale(stub, args)
	} else if fn == "checkMilkForSaleterminal" {
		return t.checkMilkForSaleterminal(stub, args)
	} else if fn == "confirmMilk" {
		return t.confirmMilk(stub, args)
	} else if fn == "getOperationHistory" {
		return t.getOperationHistory(stub, args)
	} else if fn == "get" {
		return t.get(stub, args)
	} else if fn == "set" {
		return t.set(stub, args)
	}

	return shim.Error("Machining No operation:" + fn)

}

//添加加工场
func (t *MachiningChaincode) addMachining(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("参数出错")
	}
	id := args[0]
	name := args[1]

	if isExisted(stub, id) {
		return shim.Error("已存在")
	}

	machining := Machining{}
	machining.Id = id
	machining.Name = name
	machining.CreateTime = time.Now().Format(dateFomat)
	machining.MaxMilkNo = 0
	machining.MaxBucketNo = 0
	machining.MilkIds = []string{}
	machining.BucketIds = []string{}

	jsonString, err := json.Marshal(machining)

	fmt.Println("json:" + string(jsonString))

	if err != nil {
		return shim.Error("json序列化失败")
	}

	err = stub.PutState(id, []byte(jsonString))

	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(jsonString)

}

//根据id获取加工厂
func (t *MachiningChaincode) getMachiningById(stub shim.ChaincodeStubInterface, id string) (Machining, error) {

	dvalue, err := stub.GetState(id)
	var d Machining
	err = json.Unmarshal([]byte(dvalue), &d)
	return d, err

}

//添加奶桶
func (t *MachiningChaincode) addBucket(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	bucketStr := args[0]

	var bucket Bucket
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(bucketStr), &bucket); err != nil {
		return shim.Error("json反序列化失败")
	}

	machiningId := bucket.MachiningId
	machining, err := t.getMachiningById(stub, machiningId)
	if err != nil {
		return shim.Error("加工厂不存在")
	}

	machining.MaxBucketNo += 1

	bucketIds := machining.BucketIds
	machining.BucketIds = append(bucketIds, bucket.Id)

	//跟新牛奶场
	jsonString, err := json.Marshal(machining)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	err = stub.PutState(machiningId, []byte(jsonString))
	if err != nil {
		shim.Error(err.Error())
	}

	fmt.Println("json:" + string(jsonString))
	//添加奶牛

	err = stub.PutState(bucket.Id, []byte(bucketStr))
	if err != nil {
		shim.Error(err.Error())
	}

	fmt.Println("json:" + string(bucketStr))
	return shim.Success([]byte(bucketStr))

}

//添加奶桶操作
func (t *MachiningChaincode) addMilkOperation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	operationStr := args[0]

	var operation ProcessOperation
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(operationStr), &operation); err != nil {
		return shim.Error("json反序列化失败")
	}
	operation.Time = time.Now().Format(dateFomat)
	bucketId := operation.BucketId

	operationJson, err := json.Marshal(operation)
	if err != nil {
		return shim.Error("json序列化失败")
	}

	fmt.Println("json:" + string(operationJson))
	err = stub.PutState(bucketId+operationSuffix, []byte(operationJson))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(operationJson)

}

//添加打包牛奶
func (t *MachiningChaincode) addMilkPack(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	operationStr := args[0]

	var operation ProcessOperation
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(operationStr), &operation); err != nil {
		return shim.Error("json反序列化失败")
	}
	operation.Time = time.Now().Format(dateFomat)
	operation.Operation = 2

	bucketId := operation.BucketId

	operationJson, err := json.Marshal(operation)
	if err != nil {
		return shim.Error("json序列化失败")
	}

	fmt.Println("json:" + string(operationJson))
	err = stub.PutState(bucketId+operationSuffix, []byte(operationJson))
	if err != nil {
		return shim.Error(err.Error())
	}

	bucketStr, err := stub.GetState(bucketId)
	if err != nil {
		return shim.Error(err.Error())
	}
	var bucket Bucket
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(bucketStr), &bucket); err != nil {
		return shim.Error("json反序列化失败")
	}

	machiningId := bucket.MachiningId
	machining, err := t.getMachiningById(stub, machiningId)
	if err != nil {
		return shim.Error("加工厂不存在")
	}

	machining.MaxMilkNo += 1

	milkId := bucket.Id + fmt.Sprintf(intPrefix, machining.MaxMilkNo)

	milkIds := machining.MilkIds
	machining.MilkIds = append(milkIds, milkId)

	//跟新牛奶场
	jsonString, err := json.Marshal(machining)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	err = stub.PutState(machiningId, []byte(jsonString))
	if err != nil {
		shim.Error(err.Error())
	}
	fmt.Println("json:" + string(jsonString))

	//添加奶牛
	milk := Milk{}
	milk.Id = milkId
	milk.Time = time.Now().Format(dateFomat)
	milk.Stutas = 0

	jsonString, err = json.Marshal(milk)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	err = stub.PutState(milkId, []byte(jsonString))
	if err != nil {
		shim.Error(err.Error())
	}

	fmt.Println("json:" + string(jsonString))

	return shim.Success(jsonString)

}

//发送到销售终端
func (t *MachiningChaincode) sentSale(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("参数出错")
	}

	milkId := args[0]
	saleId := args[1]

	milkStr, err := stub.GetState(milkId)
	if err != nil {
		return shim.Error(err.Error())
	}
	var milk Milk
	//这里就是实际的解码和相关的错误检查
	if err := json.Unmarshal([]byte(milkStr), &milk); err != nil {
		return shim.Error("json反序列化失败")
	}
	milk.SaleId = saleId
	milk.InSaleTime = time.Now().Format(dateFomat)

	milkStr, err = json.Marshal(milk)
	if err != nil {
		return shim.Error("json序列化失败")
	}

	err = stub.PutState(milkId, milkStr)

	if err != nil {
		return shim.Error(err.Error())
	}

	//添加销售终端查询牛奶临时组合建
	indexKey, err := stub.CreateCompositeKey(combinedConstruction, []string{saleId, milkId})
	err = stub.PutState(indexKey, milkStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("json:", string(milkStr))
	return shim.Success(nil)

}

//销售终端查询牛奶批次状态
func (t *MachiningChaincode) checkMilkForSaleterminal(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数出错")
	}

	saleId := args[0]

	var milks = make([]Milk, 0)

	resultIterator, err := stub.GetStateByPartialCompositeKey(combinedConstruction, []string{saleId})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()
	for resultIterator.HasNext() {
		item, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var milk Milk
		//这里就是实际的解码和相关的错误检查
		if err := json.Unmarshal(item.Value, &milk); err != nil {
			return shim.Error("json反序列化失败")
		}

		milks = append(milks, milk)

	}

	jsonStr, err := json.Marshal(milks)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	fmt.Println("json:", string(jsonStr))
	return shim.Success(jsonStr)

}

//确认奶状态
func (t *MachiningChaincode) confirmMilk(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("参数出错")
	}

	milkId := args[0]
	saleId := args[1]
	isConfirm := args[2]

	resultIterator, err := stub.GetStateByPartialCompositeKey(combinedConstruction, []string{saleId, milkId})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var milk Milk

	var indexKey string

	for resultIterator.HasNext() {
		item, err := resultIterator.Next()
		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}

		indexKey = item.Key
		//这里就是实际的解码和相关的错误检查
		if err := json.Unmarshal(item.Value, &milk); err != nil {
			return shim.Error("json反序列化失败")
		}

	}

	status, err := strconv.Atoi(isConfirm)
	if err != nil {
		return shim.Error(err.Error())
	}

	milk.Stutas = int8(status)

	jsonStr, err := json.Marshal(milk)
	if err != nil {
		return shim.Error("json序列化失败")
	}
	fmt.Println("json:", string(jsonStr))

	//跟新bucket
	err = stub.PutState(milkId, jsonStr)

	if err != nil {
		return shim.Error(err.Error())
	}

	//确认签收向工厂添加奶桶
	if isConfirm == "1" {
		response := stub.InvokeChaincode(saletenminalChaincodeName, [][]byte{[]byte("addMilk"), []byte(jsonStr)}, channelName)
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
			fmt.Printf(errStr)
			return shim.Error(errStr)
		}
	}

	//删除组合键
	err = stub.DelState(indexKey)
	fmt.Println("删除组合建：" + indexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsonStr)

}

func (t *MachiningChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {


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

func (t *MachiningChaincode) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
func (t *MachiningChaincode) getOperationHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("parm error")
	}
	id := args[0]
	buckerId := string(id[0:16])
	keysIter, err := stub.GetHistoryForKey(buckerId + operationSuffix)

	if err != nil {
		return shim.Error(fmt.Sprintf("GetHistoryForKey failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var historys = make([]ProcessOperation, 0)

	for keysIter.HasNext() {

		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return shim.Error(fmt.Sprintf("GetHistoryForKey operation failed. Error accessing state: %s", err))
		}
		//交易的值
		txvalue := response.Value

		var operation ProcessOperation
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
	err := shim.Start(new(MachiningChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
