# GoLang
本文档介绍如何安装Go语言环境。

## Ubuntu

### 1.下载Go二进制包：

* [64位](http://www.golangtc.com/static/go/go1.5.1/go1.5.1.linux-amd64.tar.gz)
* [32位](http://www.golangtc.com/static/go/go1.5.1/go1.5.1.linux-386.tar.gz)

在Ubuntu下，可以通过wget命令进行下载，如（64位）：

```sh
wget http://www.golangtc.com/static/go/go1.5.1/go1.5.1.linux-amd64.tar.gz
```

> 注意，不要使用apt-get命令安装Go环境，否则安装的Go版本过低。

### 2. 解压

将下载的二进制包解压到任何喜欢的目录下，如`/usr/local`：

```sh
sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
```

### 3. 创建Go工程目录

创建一个目录作为Go工作目录，如：

```sh
mkdir ~/golang
```

### 4. 配置环境变量

打开home下的bashrc,增加Go环境变量：

```sh
vi ~/.bashrc
```

在最末尾添加：

```sh
export GOROOT=/usr/local/go
export GOPATH=~/golang
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

### 5. 加载环境变量

```sh
source ~/.bashrc
```

### 6. 验证

```sh
go version
```
如果能正常显示go版本号，则证明安装成功。