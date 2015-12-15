# httpaccess

`httpaccess`服务提供了供设备注册、登录的HTTP[接口](../api-doc/device.md)。

## 启动参数
* `-etcd` etcd服务的访问地址，必需参数。如`http://localhost:2379`,如果etcd是多副本部署，可以用分号隔开访问地址，如`http://192.168.0.2:2379;http://192.168.0.3:2379`。
* `-httphost` HTTP服务地址，必须参数。格式为`ip:port`如`localhost:443`，强烈建议开启https选项。一般情况下绑定到外网的443端口（https默认端口）。
* `-usehttps` 是否启动https服务，默认不启用。如果启用，则必须提供以下`cafile`和`keyfile`两个参数。如果需要pando设备连接，则必须启用，否则无法连接。
* `-cafile` ssl加密证书的证书文件路径（pem格式）。
* `-keyfile` ssl加密证书的密钥文件路径（pem格式）。
* `-loglevel` 服务打印日志的级别，选填，如果没有指定则默认为`info`级别。

> 说明：ssl证书和密钥的pem文件生成方法可以参考[这里](http://killeraction.iteye.com/blog/858325)