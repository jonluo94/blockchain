package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	//spew可以理解为一种结构化输出工具
	"github.com/davecgh/go-spew/spew"
	//mux 是一个用于Web开发的组件
	"github.com/gorilla/mux"
	//Gotdotenv 是一个读取在项目根目录的.env文件中的配置信息的组件
	"github.com/joho/godotenv"
	"strings"
	"fmt"
)

//困难系数
const DIFFICULTY = 3

// 每个区块的数据结构
type Block struct {
	//区块号
	Index     int
	//时间戳
	Timestamp string
	//本块的哈希值
	Hash      string
	//前一个块的哈希值
	PrevHash  string
	//业务数据
	Datas     string
	//困难级别
	Difficulty int
	//随机串
	Nonce     string
}


//一个Block的数组切片用于代表区块链
var Blockchain []Block

// 请求参数结构体
type Message struct {
	Datas string
}

var mutex = &sync.Mutex{}


//共识算法POW工作量证明的hash值验证(前缀以困难系数个零开始)
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}



//要对数据进行哈希化，有两个主要原因
//节省空间，哈希值由区块中所有数据计算而来
//保护区块完整性,通过前一区块哈希检验
func calculateHash(block Block) string {
	//使用到了Index、Timestamp、Datas、PrevHash,Nonce字段用于计算当前块的哈希值
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Datas + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//用于产生新块的方法
func generateBlock(oldBlock Block, datas string) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Datas = datas
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = DIFFICULTY


	start := time.Now().UnixNano()
	for i := 0; ; i++ {

		nonce := fmt.Sprintf("%x", i)
		newBlock.Nonce = nonce
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			//time.Sleep(time.Second)
			log.Println("开始工作! ",calculateHash(newBlock))
			continue
		} else {
			newBlock.Hash = calculateHash(newBlock)
			log.Println("工作完成! ",calculateHash(newBlock))
			calculateHash(newBlock)
			break
		}

	}
	end := time.Now().UnixNano()
	log.Println("耗时:",(end-start)/1000000," 毫秒")
	return newBlock

}

// 区块正确性验证
func isBlockValid(newBlock, oldBlock Block) bool {
	//前一个索引加1等于后一块索引
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	//前一个哈希值等于后一块前哈希值
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
    //检查本块的哈希值是否正确
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

//当两个节点上的区块链长度不同时，选择较长的那个
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// 创建mux处理器
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// 处理GET请求
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//整个Blokcchain转换为JSON串作为GET请求的返回值
	io.WriteString(w, string(bytes))
}

// 使用POST请求添加新块
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	//获取参数
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], m.Datas)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		//打印区块链
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}
//当出现错误的时候返回HTTP:500，成功的返回新区块
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// http web服务器
func run() error {
	//创建mux
	mux := makeMuxRouter()
	//获取端口
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//启动服务
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}


func main() {
	//将配置写入环境变量 .env文件一定要放在项目根目录
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		//初始化区块链
		genesisBlock := Block{}
		//初始化首块
		genesisBlock = Block{0, t.String(),calculateHash(genesisBlock),"","first block",DIFFICULTY, ""}

		log.Println("First biock created ")
		//互斥锁
		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()

	run()

}















