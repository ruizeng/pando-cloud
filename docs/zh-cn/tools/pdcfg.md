# pdcfg

pdcfg是对平台进行配置的小工具。

## 概述
平台配置包含以下基本信息：

* **vendor（厂商）**：厂商是产品的所有者，平台应该至少有一个厂商。
* **product（产品）**：一个厂商会生产一个或多个产品，每个产品拥有不定数量的功能完全相同的设备。
* **application（应用）**：应用是基于平台的api接口开发的服务端程序，负责用户端业务逻辑的实现。

## 安装
```sh
go install github.com/PandoCloud/pando-cloud/tools/pdcfg
```

## 使用
根据配置的etcd访问地址，执行程序，如

```sh
$GOPATH/bin/pdcfg -etcd http://localhost:2379
```

成功连接后，会进入命令行界面，出现`>`符号并等待输入。如需退出，同时按下`ctrl`和`c`键即可。

### 添加厂商

```
vendor add
```
执行后，根据提示输入需要添加的厂商的名称和描述。成功后会打印新添加的厂商的id和key，示例:

```
> vendor add
vendor name:test12
vendor description:test fsfsf fsfs fs.
{"ip":"","level":"info","msg":"UpdateServerHosts is done: map[mqttaccess:map[rpchost:[localhost:20030] tcphost:[localhost:1883]] registry:map[rpchost:[localhost:20034]] test:map[httphost:[] rpchost:[] tcphost:[]] apiprovidor:map[httphost:[localhost:8888]] controller:map[rpchost:[localhost:20032]] devicemanager:map[rpchost:[localhost:20033]] httpaccess:map[httphost:[:443]]]","service":"pdcfg","time":"2015-12-14T16:54:12+08:00"}
=======> vendor created successfully:
ID: 39
VendorName: test12
VendorKey: 8480f0f959e788ba6918d589540f36421a9e430ee0dd6f9a50289b3bf881a1ba
VendorDescription: test fsfsf fsfs fs.
=======
```

### 添加产品
```
product add
```

执行后，根据提示输入产品所述的vendor的ID号，以及产品的名称和描述，以及产品配置的文件路径。示例：

```
> product add
vendor ID:39
product name:ffffff
product description:fffff
product config json file: /Users/ruizeng/pando.product.config.json
{"ip":"","level":"info","msg":"UpdateServerHosts is done: map[apiprovidor:map[httphost:[localhost:8888]] controller:map[rpchost:[localhost:20032]] devicemanager:map[rpchost:[localhost:20033]] httpaccess:map[httphost:[:443]] mqttaccess:map[rpchost:[localhost:20030] tcphost:[localhost:1883]] registry:map[rpchost:[localhost:20034]] test:map[rpchost:[] tcphost:[] httphost:[]]]","service":"pdcfg","time":"2015-12-14T19:16:58+08:00"}
=======> product created successfully:
ID: 20
VendorID: 0
ProductName: ffffff
ProductDescription: fffff
ProductKey: fc5f90f70fe051b4b9364c8f4b7eb060e04550fc0bfdbd313aebe56d5d107c6b
```

> 如果vendor ID不存在或为0，新增设备可以成功，vendor ID会被设置为0，表示暂无厂商，用于测试。

产品配置文件是描述产品信息的json文件，具体格式和规范请参考[产品配置](../config/product-json-config.md)

### 添加应用
```
application add
```

执行后，根据提示输入应用的名称、描述、权限、回调url以及AppToken。说明：

###### AppDomain
该字段配置了应用可管理的产品，配置完成后符合该规则的产品都可以执行应用发出的操作，并接收该产品设备的状态和事件推送。不可为空。示例：

* vendor/1  可以管理厂商ID为1的厂商名下的所有产品
* product/2  可以管理产品ID为2的产品
* *  可以管理所有设备

###### ReportURL
该字段设置了接收应用管理的设备发出的状态改变和事件通知所用的URL，平台会以HTTP post请求的方式向该URL发送状态和事件通知。外网通讯推荐采用https。

###### AppToken
该字段为应用方设置的任意字符串，平台向ReportURL发送通知时会携带该token供应用校验通知的合法性。

示例：

```
> application add
application name: app
application description: this is a test app
application domain: *
application report url: http://example.com/report
application token: test321
{"ip":"","level":"info","msg":"UpdateServerHosts is done: map[registry:map[rpchost:[localhost:20034]] test:map[rpchost:[] tcphost:[] httphost:[]] apiprovidor:map[httphost:[localhost:8888]] controller:map[rpchost:[localhost:20032]] devicemanager:map[rpchost:[localhost:20033]] httpaccess:map[httphost:[:443]] mqttaccess:map[rpchost:[localhost:20030] tcphost:[localhost:1883]]]","service":"pdcfg","time":"2015-12-14T19:36:03+08:00"}
=======> application created successfully:
ID: 4
AppKey: bca95c4973c38612c1fcec666ffc1caebd815296e9e2e34e8202ba36cecf1ce7
AppToken: test321
ReportUrl: http://example.com/report
AppName: app
AppDescription: this is a test app
AppDomain: *
=======
```


配置完成后，第三方应用就可以通过生成的AppKey调用pandocloud的RESTful API对设备进行管理。具体api接口参考[这里](../api-doc/application.md)