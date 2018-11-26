package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"math/rand"
	"bufio"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"
	"os"
	"net"
)

//在PoS中，是基于每个节点(Node)愿意作为抵押的令牌(Token)数量，
//这些参与抵押的节点被称为验证者(Validator),令牌的含义对于不同
//的区块链平台是不同的，如果验证者愿意提供更多的令牌作为抵押品，
//他们就有更大的机会记账下一个区块并获得奖励

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
	//验证者
	Validator string

}


//一个Block的数组切片用于代表区块链
var Blockchain []Block
//临时存储单元，在区块被选出来并添加到BlockChain之前，临时存储在这里
var tempBlocks []Block

// 任何一个节点在提出一个新块时都将它发送到这个通道
var candidateBlocks = make(chan Block)

// TCP服务器将向所有节点广播最新的区块链的通道
var announcements = make(chan string)

// 会保存每个节点持有的令牌数
var validators = make(map[string]int)

// 请求参数结构体
type Message struct {
	Datas string
}

var mutex = &sync.Mutex{}



//要对数据进行哈希化，有两个主要原因
//节省空间，哈希值由区块中所有数据计算而来
//保护区块完整性,通过前一区块哈希检验
//计算区块哈希
func calculateBlockHash(block Block) string {
	//使用到了Index、Timestamp、Datas、PrevHash字段用于计算当前块的哈希值
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Datas + block.PrevHash
	return calculateHash(record)
}

// SHA256哈希
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}


//用于产生新块的方法
func generateBlock(oldBlock Block, datas string, address string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Datas = datas
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
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
	if calculateBlockHash(newBlock) != newBlock.Hash {
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

//选择获胜者
//这里是PoS的主要逻辑。我们需要编写代码以实现获胜验证者的选择;他们所持有的令牌数量越高，他们就越有可能被选为胜利者
//为了简化代码，我们只会让提出新块儿的验证者参与竞争。在传统的PoS，一个验证者即使没有提出一个新的区块，也可以被选为胜利者

//pickWinner创建一个验证者的彩票池，并选择通过从池中按令牌量加权随机选择新区块
func pickWinner() {
	//每隔30秒，选出一个胜利者，这样对于每个验证者来说，都有时间提议新的区块，参与到竞争中来
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	//接着创建一个lotteryPool，它会持有所有验证者的地址
	lotteryPool := []string{}
	//判断是否已经有了被提议的区块
	if len(temp) > 0 {

		//一种改进的传统证明股份算法
		//从提交一个块的所有验证者中，根据被标记的令牌的数量对它们进行加权
		//在传统的股份证明中，验证者可以参与而不提交区块
	OUTER:
		for _, block := range temp {
			//如果已经在彩票池中，跳过
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			//防止数据竞争的验证器锁定列表
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}
		log.Println("彩票池:",lotteryPool)

		//从彩票池随机抽取赢家
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)

		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		//将块的块添加到块链中，并让所有其他节点知道
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				for _ = range validators {
					log.Println("胜利者:",lotteryWinner)
					announcements <- "\n胜利者: " + lotteryWinner + "\n"
				}
				break
			}
		}
	}

	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	//广播通知,包含了获胜验证者的通知
	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()
	//验证者地址
	var address string

	//允许用户输入其令牌数,令牌的数量越大，锻造新的块的机会就越大
	io.WriteString(conn, "输入令牌数量:")

	//控制台输入
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v 不是一个数字: %v", scanBalance.Text(), err)
			return
		}
		//该验证者被分配一个SHA256地址，随后该验证者地址和验证者的令牌数被添加到验证者列表validators中
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		break
	}

	io.WriteString(conn, "\n输入新的业务数据:")
	scanDatas := bufio.NewScanner(conn)
	go func() {
		for {
			//从控制台中获取数据并在必要验证后将其添加到区块链
			for scanDatas.Scan() {
				datas, err := strconv.Atoi(scanDatas.Text())
				//如果验证者试图提议一个被污染（例如伪造）的block，例如包含一个不是整数的数据，
				//那么程序会抛出一个错误，我们会立即从我们的验证器列表validators中删除该验证者，
				//他们将不再有资格参与到新块的铸造过程同时丢失相应的抵押令牌
				if err != nil {
					log.Printf("%v 不是一个数字: %v", scanDatas.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				mutex.Lock()
				//最新区块
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				// 创建要锻造的新块
				newBlock, err := generateBlock(oldLastIndex, strconv.Itoa(datas), address)
				if err != nil {
					log.Println(err)
					continue
				}
				//验证区块合法
				if isBlockValid(newBlock, oldLastIndex) {
					//加入待铸造的区块
					log.Println("新块加入候选区块:",newBlock)
					candidateBlocks <- newBlock
				}
			}
		}
	}()

	// 模拟接收广播
	for {
		time.Sleep(time.Minute)
		log.Println("广播最新区块链")
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}

}



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}


	t := time.Now()
	//初始化区块链
	genesisBlock := Block{}
	//初始化首块
	genesisBlock = Block{0, t.String(),calculateBlockHash(genesisBlock), "", "first block", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	httpPort := os.Getenv("PORT")

	//开启TCP服务
	server, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port :", httpPort)
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			log.Println("接收候选区块",candidate)
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	//抽取胜利者
	go func() {
		for {
			pickWinner()
		}
	}()

	//接受客户端连接
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}





