# 基于Docker快速体验

[Docker](https://www.docker.com)是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中，然后发布到运行任何流行操作系统的机器上。

我们将二进制程序打包到了一个docker镜像中，方便在基于Docker的容器平台中部署和测试。

## 依赖
Docker是跨平台的容器管理工具，在Linux，OSX以及Windows上运行。您需要先再自己的系统中[安装Docker引擎](https://docs.docker.com/engine/installation/)。安装好Docker后不再需要安装其他环境。

## 安装依赖环境
我们提供了脚本通过docker一键安装和启动平台依赖的环境如`mysql mongodb redis etd`等。下载并执行脚本：

```
wget https://raw.githubusercontent.com/PandoCloud/pando-cloud/master/build/local/docker/base.sh
sudo sh -x base.sh
```

首次运行需要下载镜像，需要较长时间，请耐心等待。

通过docker命令查看所有容器是否正常启动：

```sh
sudo docker ps -a
```

如果看到`mysql mongodb redis etd`等容器正常运行，则启动成功。

一般情况下核心服务组件不需要重新启动，服务器如果不重启只需安装一次。请不要随意删除容器，否则对应服务中的数据也会被清除。

## 创建数据库
执行以下命令创建数据库

```sh
sudo docker run --link mysql mysql sh -c 'exec mysql -hmysql -uroot -e"CREATE DATABASE PandoCloud"'
```

## 执行启动脚本
我们提供了脚本在单机环境启动平台。可以不下载代码库，直接下载脚本：

```sh
wget https://raw.githubusercontent.com/PandoCloud/pando-cloud/master/build/local/docker/run.sh
sudo sh -x run.sh
```

首次运行需要下载镜像，需要较长时间，请耐心等待。

通过docker命令查看所有容器是否正常启动：

```sh
sudo docker ps -a
```
如果看到`registry devicemanager mqttaccess controller httpaccess apiprovider`等容器正常运行，则启动成功。

## 配置
本机先编辑好[产品JSON配置文件](../config/product-json-config.md)。将配置文件保存到本机当前目录。

通过容器启动`pdcfg`配置工具，并导入配置文件到容器，对平台进行基本配置：

```sh
sudo docker run -it --name pdcfg -v `echo $(pwd)`:/root --link etcd pandocloud/pando-cloud pdcfg -etcd http://etcd:2379
```

提示：

在配置产品时，需要提供产品配置json文件路径，应该填写映射到docker容器内的路径，如`/root/product.json`而**不是**本机当前目录`./product.json`

## 测试
通过容器启动设备登陆部署的本地云平台进行测试(请使用自己配置好的product key)：

```sh
sudo docker run -it --name device --net host pandocloud/pando-cloud device -productkey=59362a15e27a0649149ff75cee1e7938f78c7cd2bb319f252694f01b7351a1
```