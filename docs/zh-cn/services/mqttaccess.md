# mqttaccess

mqttaccess是支持mqtt协议的接入服务器，平台目前默认采用该协议。提供了tcp长连接连接设备，并对设备交互数据进行转发服务。

## 启动参数

* `-etcd` etcd服务的访问地址，必需参数。如`http://localhost:2379`,如果etcd是多副本部署，可以用分号隔开访问地址，如`http://192.168.0.2:2379;http://192.168.0.3:2379`。
* `-rpchost` rpc服务访问地址，必需参数。该参数定义该服务对其他服务提供的rpc服务的监听地址。如`http://localhost:20034`。
* `-tcphost` tcp服务地址，必须参数。格式为`ip:port`如`localhost:1883`，一般绑定为外网ip加1883端口（mqtt默认端口）。
* `-usetls` 是否启动ssl加密服务，默认不启用。如果启用，则必须提供以下`cafile`和`keyfile`两个参数。如果pando设备需接入，必须开启tls加密选项，否则无法接入。
* `-cafile` ssl加密证书的证书文件路径（pem格式）。
* `-keyfile` ssl加密证书的密钥文件路径（pem格式）。
* `-loglevel` 服务打印日志的级别，选填，如果没有指定则默认为`info`级别。

> 说明：ssl证书和密钥的pem文件生成方法可以参考[这里](http://killeraction.iteye.com/blog/858325)