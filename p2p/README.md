#### p2p
* cd 到p2p目录下
* 第一个终端运行 `go run p2p.go -secio -l 10000`
* 第一个终端会打印一个链接该终端的命令,复制并粘贴到你的第二个终端,例如"go run p2p.go -l 10001 -d /ip4/127.0.0.1/tcp/10000/ipfs/QmVGUPLa6NxUWHxL6nriebF4CTXWKZo5giVm7eRbdyxtvd -secio"
* 第二终端也有类似的命令,复制和粘贴到第三个终端运行
* 三个终端都可以在控制台输入业务数据生成区块,并将最新区块链广播到所有终端
