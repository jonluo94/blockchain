App = {
    node: null,
    init: function () {
        let node = new Ipfs({repo:'ipfs-'+ Math.random()});
        node.once('ready', () => {
            console.log('Online status: ', node.isOnline() ? 'online' : 'offline')
        });
        App.node = node;
    },
    getData: function () {
        let hash = $("#hash").text();

        App.node.files.cat(hash, function (err, data) {
            if (err) {
                console.error(err)
            }
            console.log(data.toString());
            $("#outDate").text(data.toString());
        })
    },
    //上传数据
    upData: function () {
        let datas = $("#inputData").val();
        App.node.files.add(new App.node.types.Buffer(datas), (err, file) => {
            if (err) {
                console.error(err)
            }
            console.log(file[0].hash);
            $("#hash").text(file[0].hash);
        });

    },
    upFile : function () {
        let fileList = $("#fileData")[0].files;
        var reader = new FileReader();//新建一个FileReader
        reader.readAsArrayBuffer(fileList[0]);//读取文件

        reader.onloadend = function(evt){ //读取完文件之后会回来这里
            var fileString = evt.target.result; // 读取文件内容
            console.log(fileString);

            App.node.files.add(new App.node.types.Buffer(fileString), (err, file) => {
                if (err) {
                    console.error(err)
                }
                $("#fileHash").text(file[0].hash);
                $("#down").attr("href","https://ipfs.io/ipfs/"+file[0].hash);
            });
        }

    },
    getFileData: function () {
        let hash = $("#hash").text();

    },

};

$(function () {
    $(window).load(function () {
        App.init();
    });
});