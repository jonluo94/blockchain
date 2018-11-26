pragma solidity ^0.4.17;

//宠物收养合约
contract Adoption {

    address[16] public adopters;
    // 收养宠物
    function adopt(uint petId) public returns (uint) {
      require(petId >= 0 && petId <= 15);

      adopters[petId] = msg.sender;

      return petId;
    }

    // 检索收养者
    function getAdopters() public view returns (address[16]) {
      return adopters;
    }

}