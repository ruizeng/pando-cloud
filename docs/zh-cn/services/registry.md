# registry

`registry`服务负责维护平台的全局配置。

## 工作原理
该服务采用MySQL存储平台的配置和设备注册信息。

`registry`服务是无状态的，在名为`PandoCloud`的库中维护了`application`, `vendor`, `product`以及`device`表，并为了加快访问速度提供了缓存机制。

## 启动参数列表
* `-etcd` etcd服务的访问地址，必需参数。如`http://localhost:2379`,如果etcd是多副本部署，可以用分号隔开访问地址，如`http://192.168.0.2:2379;http://192.168.0.3:2379`。
* `-rpchost` rpc服务访问地址，必需参数。该参数定义该服务对其他服务提供的rpc服务的监听地址。如`http://localhost:20034`。
* `-aeskey` 用来生成KEY的aes加密密钥串，必须参数。该参数为32位的任意字符串，如`ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP`。该参数越随机越好。
* `-dbhost` MySQL数据库的访问地址，选填，如果没有填写默认为`localhost`。
* `-dbport` MySQL数据库的访问端口，选填，如果没有填写默认为`3306`。
* `-dbuser` MySQL数据库的访问用户名，选填，如果没有填写默认为`root`。
* `-dbpass` MySQL数据库的访问用户密码，选填，如果没有填写默认为空。
* `-loglevel` 服务打印日志的级别，选填，如果没有指定则默认为`info`级别。
