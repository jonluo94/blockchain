pragma solidity ^0.4.22;

/// @title 委托投票
contract Vote {

    //投票人
    struct Voter {
        uint weight; //权重
        bool voted; //是否已投票
        address delegate; //委托代理人投票
        uint voteNo; //投票提案的索引
    }

    //提案
    struct Proposal {
        bytes32 name; //简称
        uint voteCount; //得票数
    }

    //投票的发起人
    address public chairperson;

    //所有投票人
    mapping(address => Voter) public voters;

    //所有具体投票的提案
    Proposal[] public proposals;

    //构造器并为 `proposalNames` 中的每个提案，创建一个新的（投票）表决
    constructor(bytes32[] proposalNames) public {
        //发起人为合约创建人
        chairperson = msg.sender;
        //发起人的权重为1
        voters[chairperson].weight = 1;
        //初始化投票的提案
        for(uint i = 0; i < proposalNames.length; i++) {
            // `Proposal({...})` 创建一个临时 Proposal 对象，
            // `proposals.push(...)` 将其添加到 `proposals` 的末尾
            proposals.push(Proposal({
                name:proposalNames[i],
                voteCount:0
            }));
        }
    }

    //添加投票人
    function giveRightToVote(address voter) public {
        require(
            msg.sender == chairperson,
            "只有合约创建人可以添加投票人"
        );
        require(
            !voters[voter].voted,
            "投票人已经投票了"
        );
        require(
            voters[voter].weight == 0,
            "投票人没有投票资格"
        );
        //赋予投票资格
        voters[voter].weight = 1;
    }

    //把你的投票委托到他人投票
    function delegate(address to) public {
        Voter storage sender = voters[msg.sender];
        require(
            !sender.voted,
            "委托人已经投票了"
        );
        require(
            to != msg.sender,
            "代理人为自身是不允许的"
        );

        //委托是可以传递的，找到最终没有代理人的投票人
        while(voters[to].delegate != address(0)){
            //一般来说，这种循环委托是危险的。
            //因为，如果传递的链条太长,则可能需消耗的gas要多于区块中剩余的（大于区块设置的gasLimit）
            //这种情况下，委托不会被执行。
            //另一些情况下，如果形成闭环，则会让合约完全卡住。
            to = voters[to].delegate;

            require(
                to != msg.sender,
                "形成闭环，不允许闭环委托"
            );
        }

        //已投票
        sender.voted = true;
        //设置委托人
        sender.delegate = to;

        //找到代理人
        Voter storage agent = voters[to];
        //检查代理人
        if(agent.voted){
            //若代理人已经投过票了，直接增加得票数
            proposals[agent.voteNo].voteCount += sender.weight;
        }else{
            //若代理人还没投票，增加代理人的权重
            agent.weight += sender.weight;
        }

    }

    //投票(包括委托给你的票)投给提案
    function vote(uint proposalNo) public {
        //投票人
        Voter storage sender = voters[msg.sender];
        require(
            !sender.voted,
            "投票人已经投票了"
        );
        sender.voted = true;
        sender.voteNo = proposalNo;
        //如果超过了数组的范围，则会自动抛出异常，并恢复所有的改动
        proposals[proposalNo].voteCount += sender.weight;
    }

    //结合之前所有的投票，计算出最终胜出的提案索引
    function winningProposal() public view returns (uint winningProposalNo){
        uint winningVoteCount = 0;
        for (uint j = 0; j < proposals.length; j++) {
            if (proposals[j].voteCount > winningVoteCount) {
                winningVoteCount = proposals[j].voteCount;
                winningProposalNo = j;
            }
        }
    }

    // 调用函数以获取提案数组中获胜者的索引，并以此返回获胜者的名称
    function winnerName() public view returns (bytes32 winName){
        winName = proposals[winningProposal()].name;
    }



}