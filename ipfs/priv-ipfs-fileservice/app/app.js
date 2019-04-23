App = {
    peers: ["192.168.1.210", "192.168.1.213"],
    node: "",
    nodeIndex: 0,
    ipfs: null,
    //初始化
    init: function () {
        App.changeNode();
    },
    //获取文本内容
    getData: function () {

        let hash = $("#hash").text();

        App.ipfs.cat(hash, (err, data) => {
            if (err) {
                console.error(err)
            }
            console.log(data.toString());
            $("#outDate").text(data.toString());
        });

    },
    //上传文本
    upData: function () {

        let datas = $("#inputData").val();

        App.ipfs.add(App.ipfs.types.Buffer.from(datas), (err, res) => {
            if (err) {
                console.error(err)
            }
            var hash = res[0].hash
            App.ipfs.cat(hash, (err, data) => {
                if (err) {
                    console.error(err);
                }
                console.log("文本hash:" + hash);
                $("#hash").text(hash);
            })
        })
    },
    //上传文件
    upFile: function () {

        let file = $("#fileData")[0].files[0];
        let filename = file.name;
        console.log(filename);

        var reader = new FileReader();//新建一个FileReader
        reader.readAsArrayBuffer(file);//读取文件

        reader.onloadend = function (evt) { //读取完文件之后会回来这里
            var fileStream = evt.target.result; // 读取文件内容

            var fileDetails = {
                path: filename,
                content: App.ipfs.types.Buffer.from(fileStream)
            }

            var options = {
                wrapWithDirectory: true,
                progress: (prog) => console.log(`received: ${prog}`)
            }
            App.ipfs.add(fileDetails, options)
                .then((response) => {
                    console.log(response)
                    // CID of wrapping directory is returned last
                    ipfsId = response[response.length - 1].hash
                    console.log(ipfsId)
                    App.getFileData(ipfsId);


                }).catch((err) => {
                console.error(err)
            })

        }

    },
    getFileData: function (hash) {

        let name = "";
        App.ipfs.dag.get(hash, function (err, result) {
            if (err) {
                console.error('error: ' + err);
            }

            if (result.value.links.length == 1) {
                console.log(result.value.links[0].name);
                name = result.value.links[0].name;
            }
            $("#fileHash").text(hash + "/" + name);
            $("#down").attr("href", "http://" + App.node + ":8080/ipfs/" + hash + "/" + name);

        })
    },
    checkNode: function () {
        let ipfs = window.IpfsHttpClient('/ip4/' + App.node + '/tcp/5001');
        ipfs.id()
            .then(res => {
                console.log("daemon active id：" + res.id);
            })
            .catch(err => {
                if (App.nodeIndex == App.peers.length - 1) {
                    App.nodeIndex = 0
                } else {
                    App.nodeIndex++
                }
                console.error("daemon inactive ");
                App.changeNode();

            });
    },
    changeNode: function () {
        App.node = App.peers[App.nodeIndex]
        let ipfs = window.IpfsHttpClient('/ip4/' + App.node + '/tcp/5001');
        ipfs.id()
            .then(res => {
                console.log("daemon active id：" + res.id);
                App.ipfs = ipfs
            })
            .catch(err => {
                console.error("daemon inactive ");
            });
    },
};

$(function () {
    $(window).load(function () {
        App.init();
    });
    //循环执行，每隔10秒钟执行一次
    window.setInterval(function () {
        console.log("check");
        App.checkNode()
    }, 10000);

});