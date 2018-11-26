### api
#### Dairyfarm Chaincode
* addDairyFarm 方法  
  方法描述:添加奶牛场  
  参数:(string,string)   
  参数描述：第一个奶牛场id，第二个为奶牛场名  
  例子：F001,奶牛场1 
  
* addCow 方法  
  方法描述:添加奶牛   
  参数:(string)  
  参数描述：为奶牛结构体json字符串  
  例子：{"farmId":"F001","healthy":true,"quarantine":true,"feedSource":"国产","stutas":0}  
  
| 字段 | 类型 |含义|  
|:----:|:----:|:----:|  
| farmId | string | 奶牛场id |   
| healthy | bool | 健康状态 |   
| quarantine | bool | 检疫状态 | 
| feedSource | string | 食物来源 |
| stutas | int | 状态：0正常，1死亡 |

  
* delCow 方法   
  方法描述:删除奶牛   
  参数:(string)  
  参数描述：奶牛id  
  例子：F001000001  
  
* addCowOperate 方法    
  方法描述:添加操作   
  参数:(string)  
  参数描述：为操作结构体json字符串  
  例子：{"cowId":"F001000001","operation":1,"consumptionOrOutput":"b1"}  
  
| 字段 | 类型 |含义|  
|:----:|:----:|:----:|  
| cowId | string | 奶牛id |   
| operation | int |  操作类型  1为喂养 2为检疫 |   
| consumptionOrOutput | string | 消耗或产出即额外数据，建议存json字符串 | 

* addCowMilking 方法    
  方法描述:添加操作   
  参数:(string)  
  参数描述：为操作结构体json字符串  
  例子：{"cowId":"F001000001","consumptionOrOutput":"b1"}  
  
| 字段 | 类型 |含义|  
|:----:|:----:|:----:|  
| cowId | string | 奶牛id |     
| consumptionOrOutput | string | 消耗或产出即额外数据，建议存json字符串|  
    
* sentProcess 方法    
  方法描述:桶发送到工厂   
  参数:(string，string)  
  参数描述：第一个为桶id，第二个为工厂id  
  例子：F001000001000001，M001 
  
* confirmBucket 方法    
  方法描述:工厂确认奶桶是否接受，如果是确认接受会立即向工厂同步奶桶信息
  参数:(string，string，string)  
  参数描述：第一个为桶id，第二个为工厂id，第三个为 "1"为确认接受，"2"为拒绝  
  例子：F001000001000001，M001，1
  
* checkBucketForMachining 方法    
  方法描述:查看工厂待收的奶桶   
  参数:(string)  
  参数描述：工厂id  
  例子：M001 
  
* getOperationHistory 方法    
  方法描述:获取操作历史   
  参数:(string)  
  参数描述：奶牛id或者奶桶id  
  例子：F001000001 
  
* get 方法    
  方法描述:获取数据  
  参数:(string)  
  参数描述：多个key以逗号分割  
  例子：F001000001,F001000002
  
* set 方法    
  方法描述:存入数据  
  参数:(string，string)  
  参数描述：key，value  
  例子：F001000001000001，{"id":"F001000001000001","machiningId":"","time":"2018-10-04 08:40:38","stutas":0}
    
#### Machining Chaincode
* addMachining 方法  
  方法描述:添加加工厂  
  参数:(string,string)   
  参数描述：第一个为加工厂id，第二个为加工厂名  
  例子：M001,加工厂1 
  
* addBucket 方法  
  方法描述:添加奶桶   
  参数:(string)  
  参数描述：奶桶结构体json字符串  
  例子：{"id":"F001000001000001","machiningId":"M001","time":"2018-10-08 15:26:37","inMachiningTime":"2018-10-08 15:26:37","stutas":0}
  
| 字段 | 类型 |含义|
|:----:|:----:|:----:|  
| id | string | 桶id | 
| machiningId | string | 加工厂id | 
| time | string | 装桶时间 |
| inMachiningTime | string |进入加工厂时间 |
| stutas | int | 状态 |
    
* addMilkPack 方法    
  方法描述:打包牛奶  
  参数:(string)  
  参数描述：为操作结构体json字符串  
  例子：{"bucketId":"F001000001000001","consumptionOrOutput":"打包牛奶"}  
  
| 字段 | 类型 |含义|
|:----:|:----:|:----:|  
| bucketId | string | 桶id | 
| consumptionOrOutput | string | 消耗或产出即额外数据，建议存json字符串 |
 
* addMilkOperation 方法    
  方法描述:添加奶桶操作   
  参数:(string)  
  参数描述：为操作结构体json字符串  
  例子：{"bucketId":"F001000001000001","operation":1,"consumptionOrOutput":"灌装"}  
  
| 字段 | 类型 |含义|
|:----:|:----:|:----:|  
| bucketId | string | 桶id | 
| operation | int |  操作类型 0为消毒，1为灌装 | 
| consumptionOrOutput | string | 消耗或产出即额外数据，建议存json字符串 |


* sentSale 方法    
  方法描述:牛奶送到销售   
  参数:(string，string)  
  参数描述：第一个为牛奶id，第二个为销售id  
  例子：F0010000010000010001，S001 

* confirmMilk 方法    
  方法描述:销售终端确认牛奶是否接受，如果是确认接受会立即向销售终端同步牛奶信息
  参数:(string，string，string)  
  参数描述：第一个为牛奶id，第二个为销售终端id，第三个为 "1"为确认接受，"2"为拒绝  
  例子：F0010000010000010001，S001，1
  
* checkMilkForSaleterminal 方法    
  方法描述:查看销售终端待收的牛奶   
  参数:(string)  
  参数描述：终端id  
  例子：S001 

* getOperationHistory 方法    
  方法描述:获取操作历史   
  参数:(string)  
  参数描述：奶牛id或者奶桶id  
  例子：F0010000010000010001
  
* get 方法    
  方法描述:获取数据  
  参数:(string)  
  参数描述：多个key以逗号分割
  例子：F0010000010000010001,F0010000010000010002
  
* set 方法    
  方法描述:存入数据  
  参数:(string，string)  
  参数描述：key，value  
  例子：F001000001000001，{"id":"F001000001000001","machiningId":"","time":"2018-10-04 08:40:38","stutas":0}  

#### Salesterminal Chaincode
* addSalesterminal 方法  
  方法描述:添加销售终端  
  参数:(string,string)   
  参数描述：第一个为销售终端id， 第二个为销售终端名  
  例子：S001,销售终端1
  
* addMilk 方法  
  方法描述:添加奶桶   
  参数:(string)  
  参数描述：奶桶结构体json字符串  
  例子：{"id":"F0010000010000010001","time":"2018-10-08 15:48:27","InSaleTime":"2018-10-08 15:48:27","saleId":"S001","stutas":0}
  
| 字段 | 类型 |含义|
|:----:|:----:|:----:|  
| id | string | 牛奶id | 
| saleId | string | 销售端id | 
| time | string | 生产日期 |
| inSaleTime | string | 进入售买日期 ||
| stutas | int | 状态 |
  
* addOperation 方法    
  方法描述:添加操作   
  参数:(string)  
  参数描述：为操作结构体json字符串  
  例子：{"milkId":"F0010000010000010001","operation":1,"consumptionOrOutput":"售出"}  

| 字段 | 类型 |含义|
|:----:|:----:|:----:|  
| milkId | string | 牛奶id | 
| operation | int | 操作类型 0为上架，1为售出，2为下架 | 
| consumptionOrOutput | string | 消耗或产出 |
  
* getOperationHistory 方法    
  方法描述:获取操作历史   
  参数:(string)  
  参数描述：奶牛id 
  例子：F0010000010000010001
  
* getMilkHistory 方法    
  方法描述:获取溯源操作历史   
  参数:(string)  
  参数描述：奶牛id 
  例子：F0010000010000010001
  
* get 方法    
  方法描述:获取数据  
  参数:(string)  
  参数描述：多个key以逗号分割  
  例子：F0010000010000010001,F0010000010000010002
  
* set 方法    
  方法描述:存入数据  
  参数:(string，string)  
  参数描述：key，value  
  例子：F0010000010000010001，{"id":"F0010000010000010001","time":"2018-10-08 15:48:27","InSaleTime":"2018-10-08 15:48:27","saleId":"S001","stutas":0}
   