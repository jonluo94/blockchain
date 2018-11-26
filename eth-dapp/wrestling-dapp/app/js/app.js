//web3
var web3Provider = null;
//角力游戏合约
var WrestlingContract;

//初始化方法
function init() {
    //初始化web3方法用于接入区块链
    initWeb3();
    //初始化角力游戏合约
    initWrestlingContract();
}

function initWeb3() {
    //检查Web3是否已被浏览器注入 (Mist/MetaMask)
    if (typeof web3 !== 'undefined' && typeof web3.currentProvider !== 'undefined') {
        web3Provider = web3.currentProvider;
        web3 = new Web3(web3Provider);
    } else {
        console.error('Web3是没有被浏览器注入.请安装在浏览器 MetaMask 或者使用 Mist 浏览器');
        alert('Web3是没有被浏览器注入.请安装在浏览器 MetaMask 或者使用 Mist 浏览器');
    }
}


function initWrestlingContract() {
    $.getJSON('Wrestling.json', function (data) {
        //获取编译后的合约文件并用 truffle-contract 实例化得到合约实例
        WrestlingContract = TruffleContract(data);

        //为合约设置web3供应者
        WrestlingContract.setProvider(web3Provider);

        //获取玩家
        getFirstWrestlerAddress();
        getSecondWrestlerAddress();

        //获取合约发出的事件
        getEvents();

        getWinner();
        getGains();


    });
}

function getEvents() {
    WrestlingContract.deployed().then(function (instance) {

        var events = instance.allEvents(function (error, log) {
            if (!error) {
                //使用jQuery，在列表添加新事件
                switch (log.event) {
                    case 'WrestlingStartsEvent':
                        $("#eventsList").prepend('<li>事件：' + log.event +'，玩家1：'+log.args.wrestler1+'，玩家2：'+log.args.wrestler2+'</li>');
                        break;
                    case 'EndOfRoundEvent':
                        $("#eventsList").prepend('<li>事件：' + log.event +'，玩家1总投注：'+log.args.wrestler1Deposit+'，玩家1总投注：'+log.args.wrestler2Deposit+'</li>');
                        break;
                    case 'EndOfWrestlingEvent':
                        $("#eventsList").prepend('<li>事件：' + log.event +'，赢家：'+log.args.winner+'，回报：'+log.args.gains+'</li>');
                        getWinner();
                        getGains();
                        break;
                }

            }else {
                console.log(error);
            }


        });
    }).catch(function (err) {
        console.log(err.message);
    });
}

function getFirstWrestlerAddress() {
    WrestlingContract.deployed().then(function (instance) {
        //实例中属性 wrestler1 回调
        return instance.wrestler1.call();
    }).then(function (address) {
        //返回 wrestler1 的地址值
        $("#wrestler1").text(address);
    }).catch(function (err) {
        console.log(err.message);
    });
}

function getSecondWrestlerAddress() {
    WrestlingContract.deployed().then(function (instance) {
        //实例中属性 wrestler2 回调
        return instance.wrestler2.call();
    }).then(function (result) {
        //返回 wrestler2 的地址值
        if (result != "0x0000000000000000000000000000000000000000") {
            $("#wrestler2").text(result);
            $("#registerToFight").remove(); //删除按钮
        } else {
            $("#wrestler2").text("玩家2还没注册，你可以注册去参加游戏");
        }
    }).catch(function (err) {
        console.log(err.message);
    });
}


function getWinner() {
    WrestlingContract.deployed().then(function (instance) {
        //实例中属性 wrestler1 回调
        return instance.theWinner.call();
    }).then(function (address) {
        //返回 wrestler1 的地址值
        $("#winner").text(address);
    }).catch(function (err) {
        console.log(err.message);
    });
}

function getGains() {
    WrestlingContract.deployed().then(function (instance) {
        //实例中属性 wrestler1 回调
        return instance.gains.call();
    }).then(function (address) {
        //返回 wrestler1 的地址值
        $("#gains").text(address);
    }).catch(function (err) {
        console.log(err.message);
    });
}



function registerAsSecondWrestler() {
    //获取当前发起交易的帐号
    web3.eth.getAccounts(function (error, accounts) {
        if (!error) {
            if (accounts.length <= 0) {
                alert("没有帐号，请在Metamask创建帐号")
            } else {
                WrestlingContract.deployed().then(function (instance) {
                    //instance合约实例可以直接调用合约里的方法
                    //{from: accounts[0]}不属于方法的参数，但相当于solidity里的msg.sender交易的发起者
                    return instance.registerAsAnOpponent({from: accounts[0]});
                }).then(function (result) {
                    //玩家2注册后回显
                    getSecondWrestlerAddress();
                }).catch(function (err) {
                    console.log(err.message);
                });
            }

        } else {
            console.log(error);
        }
    });
}

function fight() {

    Dialog.init('<input type="text" placeholder="请输入投注的以太币"  style="margin:8px 0;width:100%;padding:11px 8px; border:1px solid #999;"/>',{
        title : '以太币数量',
        button : {
            确定 : function(){
                var amount = this.querySelector('input').value;
                Dialog.close(this);
                //获取当前发起交易的帐号
                web3.eth.getAccounts(function (error, accounts) {
                    if (!error) {
                        if (accounts.length <= 0) {
                            alert("没有帐号，请在Metamask创建帐号")
                        } else {
                            WrestlingContract.deployed().then(function (instance) {
                                //instance合约实例可以直接调用合约里的方法
                                //{from: accounts[0]}不属于方法的参数，但相当于solidity里的msg.sender交易的发起者
                                //调用payable方法要用 实例.方法.sendTransaction({from: accounts, value: web3.toWei(10, 'ether')})
                                return instance.wrestle.sendTransaction({from: accounts[0], value: web3.toWei(Number(amount), 'ether')});
                            }).then(function (result) {
                                console.log(result);
                            }).catch(function (err) {
                                console.log(err.message);
                            });
                        }

                    } else {
                        console.log(error);
                    }
                });
            },
            关闭 : function(){
                Dialog.close(this);
            }
        }
    });


}

function winnerGet() {
    //获取当前发起交易的帐号
    web3.eth.getAccounts(function (error, accounts) {
        if (!error) {
            if (accounts.length <= 0) {
                alert("没有帐号，请在Metamask创建帐号")
            } else {
                WrestlingContract.deployed().then(function (instance) {
                    //instance合约实例可以直接调用合约里的方法
                    //{from: accounts[0]}不属于方法的参数，但相当于solidity里的msg.sender交易的发起者
                    return instance.withdraw({from: accounts[0]});
                }).then(function (result) {
                    console.log(result)
                }).catch(function (err) {
                    console.log(err.message);
                });
            }

        } else {
            console.log(error);
        }
    });
}

// When the page loads, this will call the init() function
$(function () {
    $(window).load(function () {
        init();
    });
});