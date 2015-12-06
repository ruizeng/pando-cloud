# 设备接入API

## 概述
需要接入pando云平台的设备，都需要遵循pando物联网平台的接入流程和接入协议。设备相关api由设备代理(网关)向服务器发起，设备在第一次成功联网后第一时间向服务器注册自己的信息，设备通过设备登陆接口获取设备token。

## API列表
### 设备注册
*请求方式*
```
POST
```
*请求URL*
```
/v1/devices/registration
```
*参数*
```
无
```
*请求头*
```
无
```
*请求内容*
```
{
  // 产品key,平台分配给产品
  "product_key": "3d6few3ac31w7a6d3f...",
  // 设备序列号,设备唯一硬件标识，如mac地址
  "device_code": "4d3e2a5d3fff",
  // 网关版本
  "产品版本": "0.1.0"
}
```
*返回内容*
```
{
  // 返回码
  "code": 0,
  // 正确或错误信息
  "message": "", 
  // 如果成功,将返回设备id及设备密码
  "data": { 
    // 设备id
    "device_id": 12324, 
    // 设备密码
    "device_secret": "3d6few3ac31w7a6d3f", 
    // 设备激活码,用来绑定设备
    "device_key": "34ffffffff",
    // 设备标识符
    "device_identifier": "64-64-fe4efe"
  }
}
```

### 设备登录
*请求方式*
```
POST
```
*请求URL*
```
/v1/devices/authentication
```
*参数*
```
无
```
*请求头*
```
无
```
*请求内容*
```
{
  // 设备id
  "device_id": 123,
  // 设备密码
  "device_secret": "fsfwefewf23r2r32r23rfs",
  // 协议类型,不填写表示默认协议（mqtt）
  "protocol": "mqtt"
}
```
*返回内容*
```
{
  //返回码
  "code": 0, 
  //正确或错误信息
  "message": "", 
  "data": { 
    // 设备token，16个字节，经过hex编码
    "access_token": "3sffefefefefsf...", 
    // 接入服务器地址+端口
    "access_addr": "ip:port"
  }
}
```