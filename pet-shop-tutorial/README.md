## 日志
* 2019.3.27
- [项目参考地址](https://learnblockchain.cn/2018/01/12/first-dapp/#more)
- [ganache](https://truffleframework.com/docs/ganache/quickstart),代替testprc,能够方便的运行测试环境.
> 在ubuntu中下载的执行文件为`AppImage`格式,需要`chmod a+x *.AppImage`来添加权限,打开后需要添加`workspace`,不然metamask无法连接.
- [truffle文档](https://truffleframework.com/docs/truffle/getting-started/creating-a-project)
- [truffle boxes](https://truffleframework.com/docs/truffle/getting-started/creating-a-project),也就是一些truffle的开源项目
- 如果需要重新migrate合约, 需要添加reset命令,即`truffle migrate --reset`,否则会报`returned values aren't valid did it run out of gas`


* 2019.3.28
解决页面空白问题（`<div id="petTemplate" style="display: none;">` => `<div id="petTemplate" style="display: block;">`）

* 2019.3.30
按照代码逻辑,页面应该显示16个宠物,28号解决问题中的代码是前端的tmplate模板,dsiplay应该为none,而页面空白的问题是js代码有问题,有个地方少了个逗号.现在已能够正常运行.

### 遗留问题
- 2019.3.27 宠物商店的页面显示异常.其他正常.[已解决]
