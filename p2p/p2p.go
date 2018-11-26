package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
)

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
}



//一个Block的数组切片用于代表区块链
var Blockchain []Block


var mutex = &sync.Mutex{}

//要对数据进行哈希化，有两个主要原因
//节省空间，哈希值由区块中所有数据计算而来
//保护区块完整性,通过前一区块哈希检验
func calculateHash(block Block) string {
	//使用到了Index、Timestamp、Datas、PrevHash 字段用于计算当前块的哈希值
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Datas + block.PrevHash
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
	newBlock.Hash = calculateHash(newBlock)

	log.Println("created new block")
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

// 创建一个具有随机对等体ID的LIPP2P主机,给定多地址
// listenPort是主机监听的端口，其它节点会连接该端口
// secio表明是否开启数据流的安全选项，最好开启，因此它代表了”安全输入／输出”
// randSeed是一个可选的命令行标识，可以允许我们提供一个随机数种子来为我们的主机生成随机的地址,这里我们不会使用
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	//如果种子为零，则使用真正的密码随机性。否则，使用确定性随机源使生成的密钥保持不变,跨越多次运行
	var r io.Reader
	//对随机种子生成随机key
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	//为该主机生成密钥对,我们将使用它来获得有效的主机ID
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	log.Println("private key:",priv)
	if err != nil {
		return nil, err
	}

	//opts部分开始构建网络地址部分，这样其它节点就可以连接进来
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}


	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	//建立主机多地址
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	//现在我们可以通过封装两个地址来建立一个完整的多地址来到达这个主机
	addr := basicHost.Addrs()[0]
	log.Println("local address:",basicHost.Addrs())

	fullAddr := addr.Encapsulate(hostAddr)

	if secio {
		log.Printf("Now run \"go run p2p.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run p2p.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

func handleStream(s net.Stream) {

	log.Println("got a new stream!")

	//创建非阻塞读写缓冲流
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	//流会一直打开直到你关闭它（或者另一边关闭它）
	go readData(rw)
	go writeData(rw)


}

//读数据
//永不停歇的去读取外面进来的数据。首先，我们使用ReadString解析从其它节点发送过来的新的区块链（JSON字符串)
//然后检查进来的区块链的长度是否比我们本地的要长，如果进来的链更长，那么我们就接受新的链为最新的网络状态（最新的区块链)
func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {

			chain := make([]Block, 0)
			//这里就是实际的解码和相关的错误检查,并将byte数组转为Block切片
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				// 控制台绿色: 	\x1b[32m
				// 默认: 	    \x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

//写数据
//如果在我们主机的本地添加了新的区块到区块链上，那就需要把本地最新的区块链广播给其它相连的节点知道
//这些节点机会接受并更新到我们的区块链版本
func writeData(rw *bufio.ReadWriter) {

	//每5秒钟会将的最新的区块链状态广播给其它相连的节点
	go func() {
		//循环
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	//这里需要一个方法来创建一个新的Block区块
	//为了简化实现，直接在终端控制台上输入一个业务数据（Datas),然后使用generateBlock来生成区块
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)

		newBlock := generateBlock(Blockchain[len(Blockchain)-1], sendData)

		//验证区块
		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(Blockchain)

		mutex.Lock()
		//使用WriteString方法把最新的区块链广播给相连的其它节点
		//这样节点之间可以互相同步最新状态
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}

}

func main() {
	t := time.Now()
	//初始化区块链
	genesisBlock := Block{}
	//初始化首块
	genesisBlock = Block{0, t.String(),calculateHash(genesisBlock),"", "first block"}

	log.Println("First biock created ")

	Blockchain = append(Blockchain, genesisBlock)

	//日志级别
	//将debug改为info
	golog.SetAllLoggers(gologging.INFO)

	//从命令行解析参数
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")

	flag.Parse()

	log.Println("端口:",*listenF,",地址:",*target,",secio:",*secio)
    //验证端口
	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	//在给定的多地址上监听主机
	var seed int64 = 0
	ha, err := makeBasicHost(*listenF, *secio, seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		//在主机a /p2p/1.0.0 上设置流处理程序
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		//永远悬挂
		select {}
		//监听结束
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		//从给定的多地址中提取目标的对等ID
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		//从target上除掉/ipfs/<peerID>部分,使/ip4/<a.b.c.d>/ipfs/<peer>变成/ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		log.Println("target address:",targetAddr)

		//有一个对等ID和一个目标地址，所以我们把它添加到peerstore中,所以LibP2P知道如何联系它
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")

		//从主机B到主机A的新流,它应该由上面设置的处理程序在主机A上处理
		//我们使用相同的/p2p/1.0.0协议
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		//创建缓冲流，以便读取和写入是非阻塞的
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		///创建一个线程来读取和写入数据
		go writeData(rw)
		go readData(rw)
		//永远悬挂
		select {}

	}

}















