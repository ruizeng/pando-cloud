# Ubuntu
Ubuntu是基于Linux的操作系统发行版。

主流的云服务提供商如阿里云，aws等都提供安装Ubuntu操作系统的服务器。

本文档基于Ubuntu 14.04 64位操作系统（可上网），其他版本参考本文档也可。

## 目录
* [编译](#编译)
* [部署](#部署)

## 说明
如果不想自己编译二进制程序，可以使用我们预先编译好的二进制包。

所有版本的二进制包在**[这里](https://github.com/PandoCloud/pando-cloud/releases)**发布和维护，选择适合自己系统的二进制包下载并解压，**不再需要编译**，直接参考[部署](#部署)说明进行部署。

## 编译

### 1. 安装依赖

#### vcs
首先安装git和bzr，用来下载的源码，执行命令

``` sh
sudo apt-get install git -y
sudo apt-get install bzr
```

#### Go
安装Go语言编译环境，具体可以参考[这里](../environment/golang.md)
### 2. 下载源码
首先通过`go get`命令下载并编译pando cloud代码：

```sh
go get github.com/PandoCloud/pando-cloud
```

源代码就下载在Go工作目录的`$GOPATH/src/github.com/PandoCloud/pando-cloud`下。

> 下载完成后可能会提示类似`no buildable Go source files`错误，是正常的，请无视。

### 3. 编译
进入目录，执行对应的编译脚本：

```sh
cd $GOPATH/src/github.com/PandoCloud/pando-cloud
sh -x ./build/local/linux/install.sh
```

编译好的二进制文件在`$GOPATH/bin`下，可以查看该目录确认是否编译成功：

```sh
ls $GOPATH/bin
```

看到如下文件，表示所有模块编译成功：

```
apiprovider  controller  devicemanager  mqttaccess  httpaccess  registry
```

## 部署
