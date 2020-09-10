# 作业一
## 系统硬件与操作系统
由于虚拟机和云服务器都可以胜任本课程的任务，感觉使用云服务器更贴近本课程的内容，加上我本来就有一个云服务器，因此我选择了使用云服务器。
+ 服务器来源：腾讯云服务器（学生机）
+ Intel(R) Xeon(R) CPU E5-26xx v4
+ Memory 2GB
+ Disk 50G
+ CentOS Linux release 7.6.1810 (Core)


## 安装过程

### 首先配置舒服的终端环境（Visual Studio Code）
原本我是使用[Xshell](https://www.netsarang.com/zh/xshell/)+[Xftp](https://www.netsarang.com/zh/xftp/)。这套组合管理多个服务器十分方便，个人十分安利这套组合。但考虑到本课程会涉及到大量的编码，使用VS Code无论是在编程上还是远程调试上都方便许多，于是选择使用VS Code作为以后实验的主要工具。基本人人都在Windows上安装过VS Code，因此不再赘述安装过程。

1. 安装Remote SSH插件  
![SSH_r](figure/remote_ssh.png)
2. 新增配置文件  
![SSH_conf](figure/conf_ssh.png)
配置文件内容如下所示，由于懒得输密码，我使用了密钥登陆（[生成密钥教程](https://www.cnblogs.com/henkeyi/p/10487553.html)）。其中IdentityFile是本地私钥的路径，公钥放到云服务器的 ~/.ssh/authorized_keys 上。
```
    Host luowle  
        HostName luowle.cn  
        User luowle  
        IdentityFile "C:\Users\31290\.ssh\luowle_id_rsa"
```
3. 容易进行开发的终端环境就这样建立好了
![SSH](figure/ssh.png)

### golang配置
#### 安装
1. 使用yum进行安装
```
yum install golang
```
2. 查看安装目录，可见安装在 /usr/lib/golang 上
```
[root@VM_0_4_centos luowle]# rpm -ql golang |more
/etc/gdbinit.d
/etc/gdbinit.d/golang.gdb
/etc/prelink.conf.d
/etc/prelink.conf.d/golang.conf
/usr/lib/golang
/usr/lib/golang/VERSION
/usr/lib/golang/api
/usr/lib/golang/api/README
/usr/lib/golang/api/except.txt
……
……
```
3. 查看版本
```
[root@VM_0_4_centos luowle]# go version
go version go1.13.14 linux/amd64
```
#### 设置环境变量
1. 从root切换回平常使用的账户（luowle），创建工作目录
```sh
[luowle@VM_0_4_centos ~]$ mkdir $HOME/gowork
[luowle@VM_0_4_centos ~]$ ls
download  gowork  log  login.log  software  src
```
2. 配置环境变量，对于 centos 在 ~/.bashrc 文件中添加以下环境变量后重新载入
```sh
export GOPATH=$HOME/gowork
export PATH=$PATH:$GOPATH/bin
```

```sh
$ source $HOME/.bashrc
```
3. 检查配置
```
[luowle@VM_0_4_centos ~]$ go env
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOCACHE="/home/luowle/.cache/go-build"
GOENV="/home/luowle/.config/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
GOPATH="/home/luowle/gowork"
```
#### 创建Hello, world!
1. 重新登陆后创建源代码目录
```sh
$ mkdir $GOPATH/src/github.com/github-user/hello -p
```
2. 使用 VS Code 创建 hello.go
```go
package main

import "fmt"

func main() {
    fmt.Printf("hello, world\n")
}
```
3. 在终端运行
```sh
[luowle@VM_0_4_centos hello]$ go run hello.go 
hello, world
```

### 安装必要的工具和插件
#### GIT
git在yum安装的时候已经自动安装了依赖，无需再安装

#### 安装 go 的一些工具
进入 VS Code 就提示需要安装相关工具，如下图所示。但是由于众所周知的原因安装不成功。因此需要手动安装！
![go_tool](figure/go_tool.png)

1. 下载源代码到本地
```sh
# 创建文件夹
mkdir $GOPATH/src/golang.org/x/ -p
# 下载源码
go get -d github.com/golang/tools
# copy 
cp $GOPATH/src/github.com/golang/tools $GOPATH/src/golang.org/x/ -rf
```






## 问题或要点小结

