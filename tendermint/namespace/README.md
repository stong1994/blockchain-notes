## 日志
- 本项目参照官方文档构建,文档地址:https://cosmos.network/docs/tutorial/
- 项目功能类似于买卖域名.

### 启动流程
1. 安装dep和依赖包, `get_tools`和`get_vendor_deps` 在Makefile中有指定安装命令
```
make get_tools && make get_vendor_deps
```
2. 安装项目cmd, `install` 同样在Makefile中有指定安装命令
```
make install
```
此时可以查看命令
```
nsd help
nscli help
```
3. 初始化应用文件
```
nsd init --chain-id testchain
```
如果已经运行过,需要重置
```
nsd unsafe-reset-all
```
或
```
rm -rf ~/.ns*
```
4. 创建账户jack和alice, 记住得到的地址
```
nscli keys add jack
nscli keys add alicecd
```
5. 给两个账户创建钱包和代币
```
nsd add-genesis-account $(nscli keys show jack -a) 1000nametoken,1000jackcoin
nsd add-genesis-account $(nscli keys show alice -a) 1000nametoken,1000alicecoin
```
6. 配置CLI
```
nscli config chain-id testchain
nscli config output json
nscli config indent true
nscli config trust-node true
```

#### 用CLI工具进行交互
7. 启动项目
```
nsd start
```
8. 查看两个账户的余额
```
nscli query account $(nscli keys show jack -a)
nscli query account $(nscli keys show alice -a)

```
9. 从初始化的文件中用代币购买域名
```
nscli tx nameservice buy-name jack.id 5nametoken --from jack
```
10. 分配域名
```
nscli tx nameservice set-name jack.id 8.8.8.8 --from jack
```
11. 解析域名
```
nscli query nameservice resolve jack.id
```
得到
> 8.8.8.8
12. 查询whois
```
nscli query nameservice whois jack.id
```
13. alice从jack手上购买域名
```
nscli tx nameservice buy-name jack.id 10nametoken --from alice
```

#### 用REST路由进行交互
7. 启动项目
```
nscli rest-server --chain-id testchain --trust-node
```
8. 查看两个用户的余额
```
curl -s http://localhost:1317/auth/accounts/$(nscli keys show jack -a)
curl -s http://localhost:1317/auth/accounts/$(nscli keys show alice -a)
```
9. jack购买域名
```
curl -XPOST -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"jack","password":"foobarbaz","chain_id":"testchain","sequence":"2","account_number":"0"},"name":"jack1.id","amount":"5nametoken","buyer":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6"}'
```
结果
> {"check_tx":{"gasWanted":"200000","gasUsed":"1242"},"deliver_tx":{"log":"Msg 0: ","gasWanted":"200000","gasUsed":"2986","tags":[{"key":"YWN0aW9u","value":"YnV5X25hbWU="}]},"hash":"098996CD7ED4323561AC9011DEA24C70C8FAED2A4A10BC8DE2CE35C1977C3B7A","height":"23"}
10. 分配域名
```
curl -XPUT -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"jack","password":"foobarbaz","chain_id":"testchain","sequence":"3","account_number":"0"},"name":"jack1.id","value":"8.8.4.4","owner":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6"}'
```
> {"check_tx":{"gasWanted":"200000","gasUsed":"1242"},"deliver_tx":{"log":"Msg 0: ","gasWanted":"200000","gasUsed":"1352","tags":[{"key":"YWN0aW9u","value":"c2V0X25hbWU="}]},"hash":"B4DF0105D57380D60524664A2E818428321A0DCA1B6B2F091FB3BEC54D68FAD7","height":"26"}
11. 解析域名
```
curl -s http://localhost:1317/nameservice/names/jack1.id
```
结果
> 8.8.4.4
12. 查询whois
```
curl -s http://localhost:1317/nameservice/names/jack1.id/whois
> {"value":"8.8.8.8","owner":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6","price":[{"denom":"STAKE","amount":"10"}]}
```
13. alice从jack手上购买域名
```
curl -XPOST -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"alice","password":"foobarbaz","chain_id":"testchain","sequence":"1","account_number":"1"},"name":"jack1.id","amount":"10nametoken","buyer":"cosmos1h7ztnf2zkf4558hdxv5kpemdrg3tf94hnpvgsl"}'
```
> {"check_tx":{"gasWanted":"200000","gasUsed":"1264"},"deliver_tx":{"log":"Msg 0: ","gasWanted":"200000","gasUsed":"4509","tags":[{"key":"YWN0aW9u","value":"YnV5X25hbWU="}]},"hash":"81A371392B52F703266257D524538085F8C749EE3CBC1C579873632EFBAFA40C","height":"70"}
