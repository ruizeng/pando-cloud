# 应用服务器API接口

## 概述
1. api遵循REST规范
2. 所有数据内容均通过JSON格式编码
3. 文档中**请求头**为需要带在HTTP请求头中的字段；**参数**为URL查询字符串(?后的键值对如a=1&b=2)；**请求内容**为HTTP请求正文（仅post,put请求）；**请求回复**为HTTP返回正文。

### 术语
* **identifier**: 设备的唯一标识符。由平台生成,设备标识符是标示设备的全局唯一字符串。
* **status**： 设备的状态,表示设备的当前信息的数据值。
* **command**: 特定产品定义的自定义指令,由厂商配置,同一产品内唯一。
* **event**: 特定产品定义的自定义事件,由厂商配置,同一产品内唯一。
* **webhook**: 应用服务器提供给平台的上报url，平台通过该url向应用服务器上报设备动态。
* **App-Key**: 平台分配给应用后台的访问Key，用以验证应用服务器的身份。
* **App-Token**: 应用后台向平台设置的回调token，平台回调应用服务器接口时会携带该token，应用服务器通过验证该token来验证调用者的合法性。

### 应用服务器工作流程
1. 获取device key。 客户端手机app通过pando sdk进行设备配置，配置完成后可得到device key. 手机app在得到device key后应将该key上传至app应用服务器并记录该key. device key是应用服务器通过平台访问设备的唯一凭据和标识。
2. 通过device key查询device indentifier并记录。
3. 应用服务器通过调用平台的api接口对设备进行操作和交互。
4. 应用服务器准备回调接口（web hook）并按照平台文档解析回调数据. 平台会将设备状态变化及设备事件通过回调接口传给应用服务器。平台通过identifier标识设备。

### API接口列表
#### 通过devicekey获取设备信息
*请求方式：*

```
GET
```
*请求URL*

```
/application/v1/device/info
```
*参数*

```
device_key=fffffffffaaaa......  
```
> 说明：device_key是app sdk完成设备配置后拿到的device key

*请求头*

```
App-Key: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*请求内容*

```
无
```
*返回JSON示例*

``` javascript
{
  // 返回码
  "code": 0, 
  // 正确或错误信息
  "message": "", 
  "data": {
     "identifier": "设备描述符",
     "name": "设备名称",
     "description": "设备介绍",
     "version": "设备版本"
  }
}
```

#### 通过identifier查询设备信息
*请求方式：*

```
GET
```
*请求URL*

```
/application/v1/devices/{identifier}/info
```
> 说明{identifier}替换为设备标识符

*参数*

```
无
```

*请求头*

```
App-Key: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*请求内容*

```
无
```
*返回JSON示例*

``` javascript
{
  // 返回码
  "code": 0, 
  // 正确或错误信息
  "message": "", 
  "data": {
     "identifier": "设备描述符",
     "name": "设备名称",
     "description": "设备介绍",
     "version": "设备版本"
  }
}
```

#### 查询设备的状态
*请求方式：*

```
GET
```
*请求URL*

```
/application/v1/devices/{identifier}/status/current
```
> 说明{identifier}替换为设备标识符

*请求头*

```
App-Key: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*参数*

```
无
```*返回JSON示例*

``` javascript
{
  "code": 0, //返回码
  "message": "OK", //正确或错误信息
  "data": {
    "object1":[1, 23, 12, 9], //对象中所有属性值
    "object2":[1, 23, 12, 9],
    "object3":[0]
  }
}
```


#### 设置设备状态
*请求方式：*

```
PUT
```
*请求URL*

```
/application/v1/devices/{identifier}/status
```
> 说明{identifier}替换为设备标识符

*请求头*

```
App-Key: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*参数*

```
无
```
*请求内容*
``` javascript
{
  "object1":[1, 23, 12, 9], //对象中数据值，设置一个或多个同时设置
  "object2":[1, 23, 12, 9],
  "object3":[0]
}
```
*返回JSON示例*
``` javascript
{
  "code": 0, //返回码
  "message": "OK", //正确或错误信息
}
```


### webhook回调接口列表
#### 设备事件上报
*请求方式*
```
POST
```
*请求头*
```
App-Token: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*参数*
```
无
```
*请求内容JSON示例*
``` javascript
{
  "tag": "event",
  "identifier": "fffff",
  "timestamp": 12312312312312,
  "data": {
    "event事件名": [参数1,参数2...]
  }
}
```

#### 设备状态上报
*请求方式：*
```
POST
```
*请求头*
```
App-Token: 4d3e3e2e4d3e3e2e4d3e3e2e4d3e3e2e...
```
*参数*
```
无
```
*请求内容JSON示例*
``` javascript
{
  "tag": "status",
  "identifier": "fffff",
  "timestamp": 12312312312312,
  "data": {
    "object1":[1, 23, 12, 9], //对象数据值
    "object2":[1, 23, 12, 9],
    "object3":[0]
  }
}
```

### 附录
#### 返回码说明
| 返回码    | 说明     |
| -------- | -------- |
| 0        |   正常   |
| 10001        |   系统错误   |
| 10002        |   产品不存在   |
| 10003        |   设备不存在   |
| 10004        |   设备当前不在线   |
| 10005        |   错误的设备状态请求格式   |
| 10006        |   错误的产品配置   |
