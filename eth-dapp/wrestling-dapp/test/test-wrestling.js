var Wrestling = artifacts.require("Wrestling");

contract('Wrestling', function(accounts) {

    // "it" is the block to run a single test
    it("should not be able to withdraw ether", function() {
        Wrestling.deployed().then(function (inst) {
            // We retrieve the instance of the deployed Wrestling contract
            wrestlingInstance = inst;

            var account0 = accounts[0];

            // how much ether the account has before running the following transaction
            var beforeWithdraw = web3.eth.getBalance(account0);

            // We try to use the function withdraw from the Wrestling contract
            // It should revert because the wrestling isn't finished
            wrestlingInstance.withdraw({from: account0}).then(function (val) {
                assert(false, "should revert");
            }).catch(function (err) {
                // We expect a "revert" exception from the VM, because the user
                // should not be able to withdraw ether
                console.log('Error: ' + err);

                // how much ether the account has after running the transaction
                var afterWithdraw = web3.eth.getBalance(account0);
                var diff = beforeWithdraw - afterWithdraw;

                // The account paid for gas to execute the transaction
                console.log('Difference: ' + web3.fromWei(diff, "ether"));
            })
        })
    })
});