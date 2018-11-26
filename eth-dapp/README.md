# 以太坊DAPP开发
## truffle + testrpc/ganache-cli 实践
 * 描述：
   * Truffle：是以太坊的开发环境、测试框架和资产通道。可以帮助你开发、发布和测试智能合约
   * Ganache：以前叫作 TestRPC，在 TestRPC 和 Truffle 的集成后被重新命名为 Ganache  
     Ganache 的工作很简单：创建一个虚拟的以太坊区块链，并生成一些我们将在开发过程中用到的虚拟账号  
 * 1.创建项目文件夹
 * 2.在项目文件夹下，运行 truffle unbox ，初始化一个空的 truffle 项目，目录：
   ```
   .
   ├── contracts                       (合约目录)
   │   ├── ConvertLib.sol              (合约库，可删)
   │   ├── MetaCoin.sol                (代币合约，可删)
   │   └── Migrations.sol              (固件truffle迁移文件，管理和更新部署的智能合同的状态)
   ├── migrations                      (迁移目录)
   │   ├── 1_initial_migration.js      (truflle迁移部署脚本)
   │   └── 2_deploy_contracts.js       (合约迁移部署脚本)
   ├── test                            (测试目录)
   │   ├── metacoin.js                 (用JavaScript编写的测试文件，可删)
   │   └── TestMetacoin.sol            (用Solidity编写的测试文件，可删)
   └── truffle.js                      (truffle配置文件，为建立网络信息和其他项目有关的设置)
   ```
 * 3.由于创建自己的合约 DApp，所以删除 ConvertLib.sol，MetaCoin.sol，metacoin.js，TestMetacoin.sol 文件
 * 4.在 contracts 创建并编写自己的合约，我以在 remix 编写好的角力游戏合约为例:Wrestling.sol
 * 5.编写好合约之后在将 2_deploy_contracts.js 中修改成
   ```
   //导入并输出合约
   var Wrestling = artifacts.require("./Wrestling.sol");
   module.exports = function(deployer) {
     deployer.deploy(Wrestling);
   };
   ```
 * 6.在根目录下创建app文件夹，在其下面编写web端应用的源代码，目录结构为：
   ```
   .
   ├── css                           (样式)
   │   └── dialog.css                (弹窗样式)
   ├── favicon.ico                   (图标)
   ├── index.html                    (简单html也页面)
   └── js                            (js文件夹)
       ├── app.js                    (入口js文件，合约业务代码)
       └── lib                       (js库)
           ├── dialog.js             (弹窗js库)
           ├── jquery.min.js         (为了方便可以用jquery减少原生js编写)
           ├── truffle-contract.js   (必不可少的js库)
           └── web3.min.js           (必不可少的js库)
   ```
   在实际 DApp 开发中，当然可以使用诸如 React、Angular 或是 Vue 等等顺手的框架来进行开发
 * 7.web端写好后，需要使用一个小型的本地 http 服务器来为文件提供服务。采用 lite-server：
   ```
   npm init -y
   npm install lite-server --save-dev
   ```
   在项目的根目录下为 lite-server 创建一个名为 “bs-config.json” 的配置文件       
   然后添加内容：
   ```
   {
     "port": 8080,
     "server": {
       "baseDir": ["app", "build/contracts"] 
     }
   }
   ```
   接着进入 “package.json” 中的 “scripts” 节点下添加`"dev": "lite-server"`
 * 8.源代码编写好了，编译合约，会生成 build/contracts 文件夹，里面的是编译后合约ABI文件
   ```
   truffle compile
   ```
 * 9.用 ganache-cli 来模拟以太坊区块链，创建一个虚拟的以太坊区块链，并生成一些开发过程中用到的虚拟账号，每个帐号1000以太币
   ```
   ganache-cli -e 1000
   ```
 * 10.在test文件夹下编写好测试代码，运行测试
    ```
    truffle test --network development
    ```
 * 11.检查 “truffle.js” 文件设置，正确运行以下的 truffle 迁移指令来部署智能合约
   ```
   truffle migrate --network development
   ```
 * 12.在浏览器安装好 Metamask 之后，点击那个小狐狸的标志，然后点击左上角的弹出式下拉菜单，你会看到不同的接入网络，选择 “http://127.0.0.1:8545”  
   现在回到 Metamsk 主界面，点击 “restore from seed phrase” ，将 ganache 中产生的12个记忆词复制粘贴到 Wallet Seed 中，然后在下面一栏中输入你自定义的密码，这一步会解锁的第一个账户  
   接下来点击 Metamask 右上角的用户图标，然后选择 “import account” ，粘贴你从 ganache-cli 中拷贝得来的私钥，导入帐号
 * 13.运行 DApp
   ```
   npm run dev
   ```
 * 14.访问 localhost:8080 操作 DApp
 * 15.附上源代码(wrestling-dapp，pet-dapp)

## embark
 * 描述：
   * embark： 对智能合约和IPFS封装的框架 
 * 安装 npm -g install embark
 * 运行 embark demo 生成 embark_demo 项目
 * cd embark_demo
 * 打开新的控制台，运行 embark simulator
 * 运行 embark run
 * 访问 http://localhost:8000

#### 相关文档
- [Web3.JS接口文档](http://web3.tryblockchain.org/) 接口中文手册
- [Truffle框架文档](http://truffle.tryblockchain.org/) 框架中文手册
- [Ganache-cli](https://github.com/trufflesuite/ganache-cli/) 
- [Emback框架文档](https://embark.readthedocs.io/en/2.6.6/)

   