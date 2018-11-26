pragma solidity ^0.4.22;

//合约的拥有者才能操作的基础合约
contract Owned{
    //合约拥有者
    address public owner;
    //把当前交易的发送者（也就是合约的创建者）赋予合约拥有者
    constructor() public{
        owner = msg.sender;
    }

    //声明修改器，只有合约的拥有者才能操作
    modifier onlyOwner{
        require(
            msg.sender == owner,
            "只有合约拥有者才可以操作"
        );
        _;
    }
    //把该合约的拥有者转给其他人
    function transferOwner(address newOwner) public onlyOwner {
        owner = newOwner;
    }
}

contract MyCoin is Owned {
    string public name; //代币名字
    string public symbol;//代币符号
    uint8 public decimals; //代币小数位
    uint  public totalSupply; // 代币总量

    uint public sellPrice = 1 ether; //设置卖代币的价格 默认一个以太币
    uint public buyPrice = 1 ether; //设置买代币的价格 默认一个以太币

    //记录所有账户的代币的余额
    mapping(address => uint) public balanceOf;

    //用一个映射类型的变量，来记录被冻结的账户
    mapping(address=>bool) public frozenAccount;

    constructor(string _name, string _symbol,uint8 _decimals, uint _totalSupply,address _owner) public payable{
        //手动指定代币的拥有者，如果不填，则默认为合约的创建者
        if(_owner !=0){
            owner = _owner;
        }
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _totalSupply;
        //拥有者初始化拥有所有代币
        balanceOf[owner] = _totalSupply;
    }

    //设置代币的买卖价格(灵活变动价格)
    function setPrice(uint newSellPrice,uint newBuyPrice) public onlyOwner{
        sellPrice = newSellPrice;
        buyPrice = newBuyPrice;
    }

    //发行代币，向指定的目标账户添加代币
    function minterCoin(address target,uint mintedAmount) public onlyOwner{
        require(
          target != 0,
          "账户不存在"
        );
        //设置目标账户相应的代币余额
        balanceOf[target] = mintedAmount;
        //合约拥有者总量减少
        balanceOf[owner] -=mintedAmount;

    }

    //实现账户的冻结和解冻 (true 冻结，false 解冻)
    function freezeAccount(address target,bool isFreez) public onlyOwner{
        require(
          target != 0,
          "账户不存在"
        );
        frozenAccount[target] = isFreez;

    }
    //实现账户间，代币的转移
    function transfer(address _to, uint _value) public{
        //检测交易的发起者的账户是不是被冻结了
        require(
          !frozenAccount[msg.sender],
          "账户被冻结了"
        );
        //检测交易发起者的账户的代币余额是否足够
        require(
          balanceOf[msg.sender] >= _value,
          "账户代币余额不足够"
        );
        //检测溢出
        require(
          balanceOf[_to] + _value >= balanceOf[_to],
          "账户代币余额溢出"
        );

        //实现代币转移
        balanceOf[msg.sender] -=_value;
        balanceOf[_to] +=_value;
    }

    //实现代币的卖操作
    function sell(uint amount) public returns(uint revenue){
        //检测交易的发起者的账户是不是被冻结了
        require(
          !frozenAccount[msg.sender],
          "账户被冻结了"
        );
        //检测交易发起者的账户的代币余额是否足够
        require(
          balanceOf[msg.sender] >= amount,
          "账户代币余额不足够"
        );

        //把相应数量的代币给合约的拥有者
        balanceOf[owner] +=amount ;
        //卖家的账户减去相应的余额
        balanceOf[msg.sender] -=amount;
        //计算对应的以太币的价值
        revenue = amount * sellPrice;
        //向卖家的账户发送对应数量的以太币
        msg.sender.transfer(revenue);
        return;

    }

    //实现买操作
    function buy() public payable returns(uint amount) {
        //检测买家是不是大于0
        require(
          buyPrice > 0,
          "买的价格小于等于0 "
        );
        //根据用户发送的以太币的数量和代币的买价，计算出代币的数量
        amount = msg.value / buyPrice;
        //检测合约的拥有者是否有足够多的代币
        require(
          balanceOf[owner] >= amount,
          "合约的拥有者没有足够的代币"
        );
        //向合约的拥有者转移以太币
        owner.transfer(msg.value);
        //从拥有者的账户上减去相应的代币
        balanceOf[owner] -=amount ;
        //买家的账户增加相应的余额
        balanceOf[msg.sender] +=amount;

        return;
    }

}