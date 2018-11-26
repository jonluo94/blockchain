pragma solidity ^0.4.22;

//角力事件
contract  WrestlingEvent {
    //开始，玩家注册
    event WrestlingStartsEvent(address wrestler1, address wrestler2);
    //游戏期间，登记每一轮赛果
    event EndOfRoundEvent(uint wrestler1Deposit, uint wrestler2Deposit);
    //最后，其中一位玩家获胜
    event EndOfWrestlingEvent(address winner, uint gains);
}

//角力游戏
//规则;每一轮的比赛，玩家都可以投入一笔钱，如果一个玩家投入的钱是另一个玩家的两倍(总计)，那他就赢了
contract Wrestling is WrestlingEvent {
    //两个玩家
    address public wrestler1;
    address public wrestler2;

    //是否已投注
    bool public wrestler1Played;
    bool public wrestler2Played;

    //玩家总投注
    uint private wrestler1Deposit;
    uint private wrestler2Deposit;

    //游戏结束
    bool public gameFinished;
    //赢家
    address public theWinner;
    //回报
    uint public gains;

    // 在这里，第一个玩家是创造合约的人
    constructor() public {
        wrestler1 = msg.sender;
    }

    //另一个玩家进行注册
    function registerAsAnOpponent() public {
        require(
            wrestler2 == address(0),
            "玩家2已经注册"
        );
        wrestler2 = msg.sender;
        //两个玩家都注册，开始游戏
        emit WrestlingStartsEvent(wrestler1, wrestler2);
    }

    //游戏玩法：其中一方投入的资金必须是另一方的双倍
    function wrestle() public payable {
        require(
            !gameFinished,
            "游戏结束"
        );
        require(
            msg.sender == wrestler1 || msg.sender == wrestler2,
            "玩家不存在"
        );

        if(msg.sender == wrestler1) {
            require(
                wrestler1Played == false,
                "玩家已投注"
            );
            wrestler1Played = true;
            //总投注
            wrestler1Deposit += msg.value;
        } else {
            require(
                wrestler2Played == false,
                 "玩家已投注"
            );
            wrestler2Played = true;
            //总投注
            wrestler2Deposit += msg.value;
        }

        //两个玩家都投注了
        if(wrestler1Played && wrestler2Played) {
            if(wrestler1Deposit == wrestler2Deposit * 2) {
                //玩家1赢
                endOfGame(wrestler1);
            } else if (wrestler2Deposit == wrestler1Deposit * 2) {
                //玩家2赢
                endOfGame(wrestler2);
            } else {
                //进入下一轮
                endOfRound();
            }
        }
    }

    //重置投注状态，进入下一轮
    function endOfRound() internal {
        wrestler1Played = false;
        wrestler2Played = false;
        //每一轮赛果
        emit EndOfRoundEvent(wrestler1Deposit, wrestler2Deposit);
    }

    //游戏结束，得出赢家
    function endOfGame(address winner) internal {
        gameFinished = true;
        theWinner = winner;
        //投注总和，赢家回报
        gains = wrestler1Deposit + wrestler2Deposit;
        emit EndOfWrestlingEvent(winner, gains);
    }

    //游戏结束，赢家取回回报
    function withdraw() public {
        require(
            gameFinished,
            "游戏未结束"
        );
        require(
            theWinner == msg.sender,
            "玩家不是赢家"
        );

        //检查-生效-交互
        uint amount = gains;
        gains = 0;
        msg.sender.transfer(amount);
    }
}