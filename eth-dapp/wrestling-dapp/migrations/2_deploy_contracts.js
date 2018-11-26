//导入并输出合约
var Wrestling = artifacts.require("./Wrestling.sol");
module.exports = function(deployer) {
  deployer.deploy(Wrestling);
};
