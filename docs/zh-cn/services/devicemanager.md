# devicemanager

`devicemanager`维护了设备的状态，如设备当前在哪台接入服务器接入，设备是否在线等。

## 工作原理
`devicemanager`是无状态的，采用redis存储了设备的实时状态信息，并对外提供了rpc接口供其他服务查询设备的当前信息。

## 启动参数
* `-etcd` etcd服务的访问地址，必需参数。如`http://localhost:2379`,如果etcd是多副本部署，可以用分号隔开访问地址，如`http://192.168.0.2:2379;http://192.168.0.3:2379`。
* `-rpchost` rpc服务访问地址，必需参数。该参数定义该服务对其他服务提供的rpc服务的监听地址。如`http://localhost:20034`。
* `-loglevel` 服务打印日志的级别，选填，如果没有指定则默认为`info`级别。