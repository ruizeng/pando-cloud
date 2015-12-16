# Ubuntu
Ubuntu是基于Linux的操作系统发行版。

主流的云服务提供商如阿里云，aws等都提供安装Ubuntu操作系统的服务器。

本文档基于Ubuntu 14.04 64位操作系统（可上网），其他版本参考本文档也可。

## 目录
* [编译](#编译)
* [部署](#部署)

## 说明
如果不想自己编译二进制程序，可以使用我们预先编译好的二进制包。

所有版本的二进制包在**[GitHub Release](https://github.com/PandoCloud/pando-cloud/releases)**或者[百度网盘（国内推荐）](http://pan.baidu.com/s/1hqLlLFA)发布和维护，选择适合自己系统的二进制包下载并解压，**不再需要编译**，直接参考[部署](#部署)说明进行部署。

## 编译

### 1. 安装依赖

#### vcs
首先安装git和bzr，用来下载的源码，执行命令

``` sh
sudo apt-get install git -y
sudo apt-get install bzr -y
```

#### Go
安装Go语言编译环境，具体可以参考[这里](../environment/golang.md)
### 2. 下载源码
首先通过`go get`命令下载并编译pando cloud代码：

```sh
go get -u github.com/PandoCloud/pando-cloud
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
### 1. 安装依赖服务
#### MySQL

```sh
sudo apt-get install mysql-server -y
```

> 安装过程中会弹出设置root账号密码，如果只是单机体验，可以不设置，直接回车。

mysql安装完成后，请登录数据库并创建名为PandoCloud的数据库，如（无root密码）:

```sh
mysql -uroot
```

如果安装mysql时设置了root密码，则还需增加密码参数，如：

```sh
mysql -uroot -pyourpassword
```

进入后通过SQL创建数据库：

```sql
CREATE DATABASE PandoCloud;
```
**注意**：请确保数据库创建成功，否则无法正常启动registry服务。

#### Redis

```sh
sudo apt-get install redis-server -y
```

#### MongoDB

```sh
sudo apt-get install mongodb -y
```

#### Etcd

```sh
wget  https://github.com/coreos/etcd/releases/download/v2.2.2/etcd-v2.2.2-linux-amd64.tar.gz
sudo tar -C /usr/local -xzf etcd-v2.2.2-linux-amd64.tar.gz
sudo mv /usr/local/etcd-v2.2.2-linux-amd64 /usr/local/etcd
```

然后将/usr/local/etcd加入PATH环境变量

```
export PATH=$PATH:/usr/local/etcd
```

> 这样修改的PATH重启后无效，想重启后有效可以将以上命令加入到`~/.bashrc`文件末尾。

### 2. 确保服务启动
默认情况下，Redis，MySQL，MongoDB安装好后都会自动启动。也可以通过`service`命令启动和重启。

etcd需要手动启动：

```sh
etcd &
```

### 3. 执行启动脚本

启动服务前请确保步骤1中依赖的服务已经成功安装并启动

```sh
cd $GOPATH/src/github.com/PandoCloud/pando-cloud
sudo sh -x ./build/local/linux/run.sh
```

查询进程，如果6个服务都启动，则证明顺利运行。

```sh
ps aux | grep $GOPATH/bin
```

> 说明：`run.sh`脚本默认将本服务器的内外网ip的443端口（https默认端口）作为设备登陆的监听端口，将本地（`localhost`）地址的8888端口作为api接口服务的监听端口（默认情况api服务只能本机访问，如果需要通过外网访问，需要去掉httphost参数中的`localhost`）。

如果需要自定义各服务的启动参数，请参考[服务部署指南](../services/README.md)

### 4. 配置
平台提供命令行配置工具进行配置，详细使用方法请参考[配置工具](../tools/pdcfg.md)


### 5. 测试
运行测试程序`device`，可以模拟一个普通设备接入云平台，进行测试：
```sh
go install github.com/PandoCloud/pando-cloud/tests/device
$GOPATH/bin/device
```

也可以使用任何实现了[Pando嵌入式框架](https://github.com/PandoCloud/pando-embeded-framework)的设备进行测试。(登录地址改为本地http接入服务器地址)

通过HTTP请求工具（如`curl`）向api服务器发送请求就可以测试读取、设置设备状态。具体可参考[应用服务器API文档](../api-doc/application.md)
