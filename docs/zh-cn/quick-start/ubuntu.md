# Ubuntu
Ubuntu是基于Linux的操作系统发行版。

主流的云服务提供商如阿里云，aws等都提供安装Ubuntu操作系统的服务器。

本文档基于Ubuntu 14.04 64位操作系统（可上网），其他版本参考本文档也可。

## 目录
* [依赖](#依赖)
* [编译](#编译)
* [部署](#部署)

## 依赖

### apt-get
打开终端，通过ubuntu自带的包管理工具`apt-get`安装依赖的编译工具

``` sh
sudo apt-get update
```

### vcs
首先安装git和bzr，用来下载的源码，执行命令

``` sh
sudo apt-get install git -y
sudo apt-get install bzr
```

### Go
安装Go语言编译环境，具体可以参考[这里](../environment/golang.md)


## 编译
### 1. 下载源码
首先通过`go get`命令下载并编译pando cloud代码：

```sh
go get github.com/PandoCloud/pando-cloud
```

源代码就下载在Go工作目录的`$GOPATH/src/github.com/PandoCloud/pando-cloud`下。

> 下载完成后可能会提示类似`no buildable Go source files`错误，是正常的，请无视。

进入目录，执行对应的编译脚本：

```sh
cd $GOPATH/src/github.com/PandoCloud/pando-cloud
sh -x ./build/local/linux/install.sh
```

编译好的二进制文件在`$GOPATH/bin`下，可以进入该目录确认是否编译成功：

```sh
cd $GOPATH/bin | ls
```

看到如下文件，表示所有模块编译成功：

```

```

### 2. 编译
进入该目录，执行build/local/linux目录下的`install.sh`脚本：

```sh
cd $GOPATH/src/github.com/PandoCloud/pando-cloud
./build/local/linux/install.sh
```



## 部署

### 物理机部署
#### 依赖

#### 准备工作