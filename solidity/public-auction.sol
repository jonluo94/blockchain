pragma solidity ^0.4.22;

//拍卖事件
contract AuctionEvent {
     //出现更高投标价格时，触发该事件，公开投标人和投标金额
    event HighestBidIncreased(address bidder, uint amount);
    //投标结束，通过事件公开投标获胜者和投标金额
    event AuctionEnded(address winner, uint amount);
    //退回未中标的投标者的钱
    event DoWithdraw(bool _b,address _addr, uint amount);
}

//公开拍卖
contract PublicAuction is AuctionEvent{
    address public beneficiary; //拍卖受益人
    uint public auctionEnd; //拍卖结束时间(时间是unix的绝对时间戳)

    address public highestBidder; //最高投标金额的投标人
    uint public highestBid; //最高投标金额

    mapping(address => uint) public pendingReturns;//每个人的投标总金额(可以取回的之前的出价)
    bool public ended;//拍卖结束

    //合约创建时启动拍卖，初始化时长和受益人
    constructor(uint _biddingTime, address _beneficiary) public{
        beneficiary = _beneficiary;
        auctionEnd = now + _biddingTime;
    }

    //对拍卖进行出价，具体的出价随交易一起发送。
    //如果没有在拍卖中胜出，则返还出价。
    function bid() public payable {

        // 对于能接收以太币的函数，关键字 payable 是必须的。
        // 如果拍卖已结束，撤销函数的调用。
        require(
            now <= auctionEnd,
            "拍卖已结束"
        );

        // 如果出价不够高，返还你的钱
        require(
            msg.value > highestBid,
            "出价不够高"
        );

        if (highestBid != 0) {
            // 返还出价时，简单地直接调用 highestBidder.send(highestBid) 函数，
            // 是有安全风险的，因为它有可能执行一个非信任合约。
            // 更为安全的做法是让接收方自己提取金钱。
            pendingReturns[highestBidder] += highestBid;
        }
        highestBidder = msg.sender;
        highestBid = msg.value;
        emit HighestBidIncreased(msg.sender, msg.value);
    }

    //退回未中标的资金
    function withdraw() public returns(bool) {

        uint amount = pendingReturns[msg.sender];
        if (amount > 0) {
            // 这里很重要，首先要设零值。
            // 因为，作为接收调用的一部分，
            // 接收者可以在 `send` 返回之前，重新调用该函数。
            pendingReturns[msg.sender] = 0;

            if (!msg.sender.send(amount)) {
                //如果以太币发送失败，则重置需要返还的钱
                pendingReturns[msg.sender] = amount;
                emit DoWithdraw(false,msg.sender,amount);
                return false;
            }else{
                emit DoWithdraw(true,msg.sender,amount);
            }

        }
        return true;
    }


     //结束此次拍卖，并把最高的出价发送给受益人
    function auctionEnd() public {
        //如果时间还没到，则终止程序
        require(now >= auctionEnd, "时间还没到");
        //已经被结束过了
        require(!ended, "已经被结束过了");

        //结束生效
        ended = true;
        //将获胜者的信息显示
        emit AuctionEnded(highestBidder, highestBid);

        //将中标的金额发送给拍卖受益人
        beneficiary.transfer(highestBid);
    }


}
