# 基于Docker快速体验

[Docker](https://www.docker.com)是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中，然后发布到运行任何流行操作系统的机器上。

我们将二进制程序打包到了一个docker镜像中，方便在基于Docker的容器平台中部署和测试。

## 依赖
Docker是跨平台的容器管理工具，在Linux，OSX以及Windows上运行。您需要先再自己的系统中[安装Docker引擎](https://docs.docker.com/engine/installation/)。安装好Docker后不再需要安装其他环境。

## 执行启动脚本

我们提供了脚本在单机环境启动平台：

```sh
sudo sh -x ./build/local/docker/run.sh

```

首次运行需要下载镜像，需要较长时间，请耐心等待。

通过docker命令查看所有容器是否正常启动：

```
sudo docker ps -a
```