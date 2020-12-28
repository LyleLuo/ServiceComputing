# docker 实践
## docker 安装

### 卸载旧版本
```
yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
```
### 安装依赖

```
yum install -y yum-utils
```
### 换源
```
yum-config-manager \
    --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

sed -i 's/download.docker.com/mirrors.aliyun.com\/docker-ce/g' /etc/yum.repos.d/docker-ce.repo
```
### 安装 docker
```
yum install docker-ce docker-ce-cli containerd.io
```
### 启动 docker
```
systemctl enable docker
systemctl start docker
```

### 测试是否安装成功
```
[root@localhost lyle]# docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
0e03bdcc26d7: Pull complete 
Digest: sha256:1a523af650137b8accdaed439c17d684df61ee4d74feac151b5b337bd29e7eec
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

### 建立 docker 用户组
为了避免使用root用户，将需要使用 docker 的用户加入 docker 用户组。
```
usermod -aG docker lyle
```



## Docker 基本操作

### 显示本地镜像库内容

```
[lyle@localhost ~]$ docker images
REPOSITORY    TAG       IMAGE ID       CREATED         SIZE
hello-world   latest    bf756fb1ae65   11 months ago   13.3kB
```

### 获得帮助

```
[lyle@localhost ~]$  docker --help

Usage:  docker [OPTIONS] COMMAND

A self-sufficient runtime for containers

Options:
      --config string      Location of client config files (default "/home/lyle/.docker")
  -c, --context string     Name of the context to use to connect to the daemon (overrides DOCKER_HOST env var and default context set with "docker context use")
  -D, --debug              Enable debug mode
  -H, --host list          Daemon socket(s) to connect to
  -l, --log-level string   Set the logging level ("debug"|"info"|"warn"|"error"|"fatal") (default "info")
      --tls                Use TLS; implied by --tlsverify
      --tlscacert string   Trust certs signed only by this CA (default "/home/lyle/.docker/ca.pem")
      --tlscert string     Path to TLS certificate file (default "/home/lyle/.docker/cert.pem")
      --tlskey string      Path to TLS key file (default "/home/lyle/.docker/key.pem")
      --tlsverify          Use TLS and verify the remote
  -v, --version            Print version information and quit

Management Commands:
  app*        Docker App (Docker Inc., v0.9.1-beta3)
  builder     Manage builds
  buildx*     Build with BuildKit (Docker Inc., v0.5.0-docker)
  config      Manage Docker configs
  container   Manage containers
  context     Manage contexts
  image       Manage images
  manifest    Manage Docker image manifests and manifest lists
  network     Manage networks
  node        Manage Swarm nodes
  plugin      Manage plugins
  secret      Manage Docker secrets
  service     Manage services
  stack       Manage Docker stacks
  swarm       Manage Swarm
  system      Manage Docker
  trust       Manage trust on Docker images
  volume      Manage volumes

Commands:
  attach      Attach local standard input, output, and error streams to a running container
  build       Build an image from a Dockerfile
  commit      Create a new image from a container's changes
  cp          Copy files/folders between a container and the local filesystem
  create      Create a new container
  diff        Inspect changes to files or directories on a container's filesystem
  events      Get real time events from the server
  exec        Run a command in a running container
  export      Export a container's filesystem as a tar archive
  history     Show the history of an image
  images      List images
  import      Import the contents from a tarball to create a filesystem image
  info        Display system-wide information
  inspect     Return low-level information on Docker objects
  kill        Kill one or more running containers
  load        Load an image from a tar archive or STDIN
  login       Log in to a Docker registry
  logout      Log out from a Docker registry
  logs        Fetch the logs of a container
  pause       Pause all processes within one or more containers
  port        List port mappings or a specific mapping for the container
  ps          List containers
  pull        Pull an image or a repository from a registry
  push        Push an image or a repository to a registry
  rename      Rename a container
  restart     Restart one or more containers
  rm          Remove one or more containers
  rmi         Remove one or more images
  run         Run a command in a new container
  save        Save one or more images to a tar archive (streamed to STDOUT by default)
  search      Search the Docker Hub for images
  start       Start one or more stopped containers
  stats       Display a live stream of container(s) resource usage statistics
  stop        Stop one or more running containers
  tag         Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
  top         Display the running processes of a container
  unpause     Unpause all processes within one or more containers
  update      Update configuration of one or more containers
  version     Show the Docker version information
  wait        Block until one or more containers stop, then print their exit codes

Run 'docker COMMAND --help' for more information on a command.

To get more help with docker, check out our guides at https://docs.docker.com/go/guides/
```



### 显示运行中容器

```
[lyle@localhost ~]$ docker ps
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
```

### 显示所有容器（包含已中止）

```
[lyle@localhost ~]$ docker ps -a
CONTAINER ID   IMAGE         COMMAND    CREATED             STATUS                         PORTS     NAMES
b272f629b4ef   hello-world   "/hello"   57 minutes ago      Exited (0) 57 minutes ago                gifted_elion
b13d474f5649   hello-world   "/hello"   About an hour ago   Exited (0) About an hour ago             recursing_feynman
```



## 构建镜像练习

```
[lyle@localhost Desktop]$ mkdir mydock && cd mydock
[lyle@localhost mydock]$ vi dockerfile
[lyle@localhost mydock]$ ls
dockerfile
[lyle@localhost mydock]$ cat dockerfile 
FROM ubuntu
ENTRYPOINT ["top", "-b"]
CMD ["-c"]
```

构建镜像

```
[lyle@localhost mydock]$ docker build . -t hello
Sending build context to Docker daemon  2.048kB
Step 1/3 : FROM ubuntu
latest: Pulling from library/ubuntu
da7391352a9b: Pull complete 
14428a6d4bcd: Pull complete 
2c2d948710f2: Pull complete 
Digest: sha256:c95a8e48bf88e9849f3e0f723d9f49fa12c5a00cfc6e60d2bc99d87555295e4c
Status: Downloaded newer image for ubuntu:latest
 ---> f643c72bc252
Step 2/3 : ENTRYPOINT ["top", "-b"]
 ---> Running in bd0d94a9b2ff
Removing intermediate container bd0d94a9b2ff
 ---> 5f0353974902
Step 3/3 : CMD ["-c"]
 ---> Running in c263b7927469
Removing intermediate container c263b7927469
 ---> f77936ca864a
Successfully built f77936ca864a
Successfully tagged hello:latest

```

运行镜像

```
[lyle@localhost mydock]$ docker run -it --rm hello -H 
top - 08:44:40 up 57 min,  0 users,  load average: 0.05, 0.35, 0.39
Threads:   1 total,   1 running,   0 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.0 us, 50.0 sy,  0.0 ni, 50.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :   7802.9 total,   5511.6 free,   1070.2 used,   1221.0 buff/cache
MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   6458.7 avail Mem 

   PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
     1 root      20   0    5960   1708   1280 R   0.0   0.0   0:00.21 top

top - 08:44:43 up 57 min,  0 users,  load average: 0.04, 0.35, 0.38
Threads:   1 total,   1 running,   0 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.3 us,  0.7 sy,  0.0 ni, 99.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :   7802.9 total,   5512.5 free,   1069.3 used,   1221.1 buff/cache
MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   6459.6 avail Mem 

   PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
     1 root      20   0    5960   1708   1280 R   0.0   0.0   0:00.21 top

top - 08:44:46 up 57 min,  0 users,  load average: 0.04, 0.35, 0.38
Threads:   1 total,   1 running,   0 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.4 us,  0.7 sy,  0.0 ni, 98.8 id,  0.0 wa,  0.0 hi,  0.1 si,  0.0 st
MiB Mem :   7802.9 total,   5512.6 free,   1069.1 used,   1221.1 buff/cache
MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   6459.8 avail Mem 

   PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
     1 root      20   0    5960   1708   1280 R   0.0   0.0   0:00.21 top^C

```



## 镜像源修改

```
tc/docker
[root@localhost docker]# vi daemon.json
```

添加以下内容

```
{
"registry-mirrors": ["http://hub-mirror.c.163.com"]
}
```



## MySQL与容器化

### 拉取 MySQL 镜像

```
[lyle@localhost Desktop]$ docker pull mysql:5.7
5.7: Pulling from library/mysql
6ec7b7d162b2: Pull complete 
fedd960d3481: Pull complete 
7ab947313861: Pull complete 
64f92f19e638: Pull complete 
3e80b17bff96: Pull complete 
014e976799f9: Pull complete 
59ae84fee1b3: Pull complete 
7d1da2a18e2e: Pull complete 
301a28b700b9: Pull complete 
529dc8dbeaf3: Pull complete 
bc9d021dc13f: Pull complete 
Digest: sha256:c3a567d3e3ad8b05dfce401ed08f0f6bf3f3b64cc17694979d5f2e5d78e10173
Status: Downloaded newer image for mysql:5.7
docker.io/library/mysql:5.7
```

## 使用MySQL容器

### 启动服务器

```
[lyle@localhost Desktop]$  docker run -p 3306:3306 --name mysql2 -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7
713af3fb742bced496b9b826d033f412760d8499ab8b937b49a6e92d2673fb28
```

```
[lyle@localhost Desktop]$ docker ps
CONTAINER ID   IMAGE       COMMAND                  CREATED         STATUS         PORTS                               NAMES
713af3fb742b   mysql:5.7   "docker-entrypoint.s…"   2 minutes ago   Up 2 minutes   0.0.0.0:3306->3306/tcp, 33060/tcp   mysql2

```

### 启动 MySQL 客户端

```
[lyle@localhost Desktop]$ docker run -it --net host mysql:5.7 "sh"
# mysql -h127.0.0.1 -P3306 -uroot -proot
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.32 MySQL Community Server (GPL)

Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> 

```

### 挂载卷保存db

#### 数据库文件在哪里？

```
[lyle@localhost Desktop]$ docker exec -it mysql2 bash
root@713af3fb742b:/#  ls /var/lib/mysql
auto.cnf	 client-key.pem  ibdata1	     private_key.pem  sys
ca-key.pem	 ib_buffer_pool  ibtmp1		     public_key.pem
ca.pem		 ib_logfile0	 mysql		     server-cert.pem
client-cert.pem  ib_logfile1	 performance_schema  server-key.pem
```



#### Dockerfile 的 VOLUME /var/lib/mysql 的含义

```
[lyle@localhost Desktop]$ docker volume ls
DRIVER    VOLUME NAME
local     062d7f8fc113076572c74bb22e3963c485785942b8374c40c46031364093e548
local     f1791196472dc6f7e22f704912cef0b4dc4705a3e651837d28dea39fe9683cb7
```

- 每次启动 mysql 容器，docker 创建一个文件卷挂载在容器内/var/lib/mysql位置 

- 这个卷在主机（host）的 /var/lib/docker/volumes/ 目录下

### 创建卷并挂载

```
[lyle@localhost Desktop]$  docker rm $(docker ps -a -q) -f -v
713af3fb742b
[lyle@localhost Desktop]$  docker volume create mydb
mydb
[lyle@localhost Desktop]$  docker run --name mysql2 -e MYSQL_ROOT_PASSWORD=root -v mydb:/var/lib/mysql -d mysql:5.7
167d05e1f1c81831ee09ec0b917a643a26544d99dbf28e2137f69ac8f9e1951e
```

### 启动客户端容器链接服务器

```
[lyle@localhost Desktop]$  docker run --name myclient --link mysql2:mysql -it mysql:5.7 bash
root@2199fbd8a57d:/# env
MYSQL_PORT_33060_TCP_ADDR=172.17.0.2
MYSQL_PORT=tcp://172.17.0.2:3306
MYSQL_PORT_3306_TCP_ADDR=172.17.0.2
MYSQL_NAME=/myclient/mysql
MYSQL_ENV_MYSQL_ROOT_PASSWORD=root
MYSQL_MAJOR=5.7
MYSQL_PORT_3306_TCP_PORT=3306
HOSTNAME=2199fbd8a57d
MYSQL_ENV_MYSQL_MAJOR=5.7
MYSQL_PORT_3306_TCP=tcp://172.17.0.2:3306
PWD=/
HOME=/root
MYSQL_ENV_GOSU_VERSION=1.12
MYSQL_PORT_33060_TCP_PROTO=tcp
MYSQL_VERSION=5.7.32-1debian10
GOSU_VERSION=1.12
TERM=xterm
MYSQL_PORT_33060_TCP_PORT=33060
MYSQL_PORT_3306_TCP_PROTO=tcp
SHLVL=1
MYSQL_PORT_33060_TCP=tcp://172.17.0.2:33060
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
MYSQL_ENV_MYSQL_VERSION=5.7.32-1debian10
_=/usr/bin/env
root@2199fbd8a57d:/#  mysql -hmysql -P3306 -uroot -proot
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.32 MySQL Community Server (GPL)

Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> 
```

### 挂载现有数据库

```
[lyle@localhost Desktop]$ docker run -v "$PWD/data":/var/lib/mysql --user 1000:1000 --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:5.7
e2dfe013f2d36d1cde445b10020b0bd1386b95ab359763626479a7f4491fef2f
```

### 修改容器配置

```
[lyle@localhost Desktop]$ docker run --name some-mysql2 -v /my/custom:/etc/mysql/conf.d -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:5.7
ce346b623e9bdec90e7a3602284429edff074f95b7ebdb48d38d302cb0f4271a
```

## docker 网络

### 容器网络内容

```
[lyle@localhost Desktop]$ docker network ls
NETWORK ID     NAME      DRIVER    SCOPE
869a959200ea   bridge    bridge    local
7151c5809a16   host      host      local
30b2026d0988   none      null      local

```

### 备制支持 ifconfig 和 ping 命令的 ubuntu 容器

```
[lyle@localhost Desktop]$  docker run --name unet -it --rm ubuntu bash
root@27c6a5d2ab45:/#  apt-get update
Get:1 http://archive.ubuntu.com/ubuntu focal InRelease [265 kB]
Get:2 http://security.ubuntu.com/ubuntu focal-security InRelease [109 kB]
Get:3 http://archive.ubuntu.com/ubuntu focal-updates InRelease [114 kB]
Get:4 http://archive.ubuntu.com/ubuntu focal-backports InRelease [101 kB]
Get:5 http://archive.ubuntu.com/ubuntu focal/restricted amd64 Packages [33.4 kB]
Get:6 http://archive.ubuntu.com/ubuntu focal/multiverse amd64 Packages [177 kB]
Get:7 http://archive.ubuntu.com/ubuntu focal/universe amd64 Packages [11.3 MB] 
Get:8 http://security.ubuntu.com/ubuntu focal-security/restricted amd64 Packages [103 kB]
Get:9 http://security.ubuntu.com/ubuntu focal-security/multiverse amd64 Packages [1167 B]
Get:10 http://security.ubuntu.com/ubuntu focal-security/universe amd64 Packages [645 kB]
Get:11 http://archive.ubuntu.com/ubuntu focal/main amd64 Packages [1275 kB]    
Get:12 http://archive.ubuntu.com/ubuntu focal-updates/multiverse amd64 Packages [30.4 kB]
Get:13 http://archive.ubuntu.com/ubuntu focal-updates/restricted amd64 Packages [136 kB]
Get:14 http://archive.ubuntu.com/ubuntu focal-updates/universe amd64 Packages [885 kB]
Get:15 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 Packages [885 kB]
Get:16 http://archive.ubuntu.com/ubuntu focal-backports/universe amd64 Packages [4250 B]
Get:17 http://security.ubuntu.com/ubuntu focal-security/main amd64 Packages [495 kB]
Fetched 16.6 MB in 19s (881 kB/s)                                              
Reading package lists... Done
root@27c6a5d2ab45:/# apt-get install net-tools
Reading package lists... Done
Building dependency tree       
Reading state information... Done
The following NEW packages will be installed:
  net-tools
0 upgraded, 1 newly installed, 0 to remove and 2 not upgraded.
Need to get 196 kB of archives.
After this operation, 864 kB of additional disk space will be used.
Get:1 http://archive.ubuntu.com/ubuntu focal/main amd64 net-tools amd64 1.60+git20180626.aebd88e-1ubuntu1 [196 kB]
Fetched 196 kB in 1s (159 kB/s)     
debconf: delaying package configuration, since apt-utils is not installed
Selecting previously unselected package net-tools.
(Reading database ... 4121 files and directories currently installed.)
Preparing to unpack .../net-tools_1.60+git20180626.aebd88e-1ubuntu1_amd64.deb ...
Unpacking net-tools (1.60+git20180626.aebd88e-1ubuntu1) ...
Setting up net-tools (1.60+git20180626.aebd88e-1ubuntu1) ...
root@27c6a5d2ab45:/# apt-get install iputils-ping -y
Reading package lists... Done
Building dependency tree       
Reading state information... Done
The following additional packages will be installed:
  libcap2 libcap2-bin libpam-cap
The following NEW packages will be installed:
  iputils-ping libcap2 libcap2-bin libpam-cap
0 upgraded, 4 newly installed, 0 to remove and 2 not upgraded.
Need to get 90.5 kB of archives.
After this operation, 333 kB of additional disk space will be used.
Get:1 http://archive.ubuntu.com/ubuntu focal/main amd64 libcap2 amd64 1:2.32-1 [15.9 kB]
Get:2 http://archive.ubuntu.com/ubuntu focal/main amd64 libcap2-bin amd64 1:2.32-1 [26.2 kB]
Get:3 http://archive.ubuntu.com/ubuntu focal/main amd64 iputils-ping amd64 3:20190709-3 [40.1 kB]
Get:4 http://archive.ubuntu.com/ubuntu focal/main amd64 libpam-cap amd64 1:2.32-1 [8352 B]
Fetched 90.5 kB in 2s (39.5 kB/s)    
debconf: delaying package configuration, since apt-utils is not installed
Selecting previously unselected package libcap2:amd64.
(Reading database ... 4170 files and directories currently installed.)
Preparing to unpack .../libcap2_1%3a2.32-1_amd64.deb ...
Unpacking libcap2:amd64 (1:2.32-1) ...
Selecting previously unselected package libcap2-bin.
Preparing to unpack .../libcap2-bin_1%3a2.32-1_amd64.deb ...
Unpacking libcap2-bin (1:2.32-1) ...
Selecting previously unselected package iputils-ping.
Preparing to unpack .../iputils-ping_3%3a20190709-3_amd64.deb ...
Unpacking iputils-ping (3:20190709-3) ...
Selecting previously unselected package libpam-cap:amd64.
Preparing to unpack .../libpam-cap_1%3a2.32-1_amd64.deb ...
Unpacking libpam-cap:amd64 (1:2.32-1) ...
Setting up libcap2:amd64 (1:2.32-1) ...
Setting up libcap2-bin (1:2.32-1) ...
Setting up libpam-cap:amd64 (1:2.32-1) ...
debconf: unable to initialize frontend: Dialog
debconf: (No usable dialog-like program is installed, so the dialog based frontend cannot be used. at /usr/share/perl5/Debconf/FrontEnd/Dialog.pm line 76.)
debconf: falling back to frontend: Readline
debconf: unable to initialize frontend: Readline
debconf: (Can't locate Term/ReadLine.pm in @INC (you may need to install the Term::ReadLine module) (@INC contains: /etc/perl /usr/local/lib/x86_64-linux-gnu/perl/5.30.0 /usr/local/share/perl/5.30.0 /usr/lib/x86_64-linux-gnu/perl5/5.30 /usr/share/perl5 /usr/lib/x86_64-linux-gnu/perl/5.30 /usr/share/perl/5.30 /usr/local/lib/site_perl /usr/lib/x86_64-linux-gnu/perl-base) at /usr/share/perl5/Debconf/FrontEnd/Readline.pm line 7.)
debconf: falling back to frontend: Teletype
Setting up iputils-ping (3:20190709-3) ...
Processing triggers for libc-bin (2.31-0ubuntu9.1) ...

```

```
root@27c6a5d2ab45:/# ifconfig
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 172.17.0.4  netmask 255.255.0.0  broadcast 172.17.255.255
        ether 02:42:ac:11:00:04  txqueuelen 0  (Ethernet)
        RX packets 3340  bytes 17075898 (17.0 MB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 2079  bytes 117352 (117.3 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        loop  txqueuelen 1000  (Local Loopback)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

root@27c6a5d2ab45:/# ping localhost
PING localhost (127.0.0.1) 56(84) bytes of data.
64 bytes from localhost (127.0.0.1): icmp_seq=1 ttl=64 time=0.065 ms
64 bytes from localhost (127.0.0.1): icmp_seq=2 ttl=64 time=0.110 ms
^C
--- localhost ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 0.065/0.087/0.110/0.022 ms
root@27c6a5d2ab45:/# ping 192.168.150.134
PING 192.168.150.134 (192.168.150.134) 56(84) bytes of data.
64 bytes from 192.168.150.134: icmp_seq=1 ttl=64 time=0.252 ms
64 bytes from 192.168.150.134: icmp_seq=2 ttl=64 time=0.179 ms
64 bytes from 192.168.150.134: icmp_seq=3 ttl=64 time=0.173 ms
^C
--- 192.168.150.134 ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2000ms
rtt min/avg/max/mdev = 0.173/0.201/0.252/0.035 ms

```

### 启动另一个命令窗口，由容器制作镜像

```
[lyle@localhost ~]$  docker commit unet ubuntu:net
sha256:7f0075d6104ff9d260b9cc07883cafcc3311c2e9a1af3500bc222c22e103a038
```

### 创建自定义网络

```
[lyle@localhost Desktop]$  docker network create mynet
0fa7b38a868b3f8790a9ee607740add2a0e2956213825f635a8fb74c5c256d80
[lyle@localhost Desktop]$ docker network create u1
f2ef9862cf35bdd6f81cbd97159e237388d3fb4a3ce74363ff6bdff6c41ef98d
[lyle@localhost Desktop]$ docker network create u2
f5b73aeb508204e682dccf48d11f6a163abee81f0a1860cfc12dd66b9503c90a
```

```
[lyle@localhost Desktop]$ docker run --name u1 -it -p 8080:80 --net mynet --rm ubuntu:net bash
root@76f0fb3e84f2:/# 
```

```
le@localhost ~]$ docker run --name u2 --net mynet -it --rm ubuntu:net bash
root@7691906f6eb6:/#
```

```
[lyle@localhost ~]$ docker inspect u1
[
    {
        "Id": "76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679",
        "Created": "2020-12-28T15:05:35.811004541Z",
        "Path": "bash",
        "Args": [],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false,
            "Restarting": false,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 24999,
            "ExitCode": 0,
            "Error": "",
            "StartedAt": "2020-12-28T15:05:37.023891145Z",
            "FinishedAt": "0001-01-01T00:00:00Z"
        },
        "Image": "sha256:7f0075d6104ff9d260b9cc07883cafcc3311c2e9a1af3500bc222c22e103a038",
        "ResolvConfPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/hostname",
        "HostsPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/hosts",
        "LogPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679-json.log",
        "Name": "/u1",
        "RestartCount": 0,
        "Driver": "overlay2",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": null,
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {}
            },
            "NetworkMode": "mynet",
            "PortBindings": {
                "80/tcp": [
                    {
                        "HostIp": "",
                        "HostPort": "8080"
                    }
                ]
            },
            "RestartPolicy": {
                "Name": "no",
                "MaximumRetryCount": 0
            },
            "AutoRemove": true,
            "VolumeDriver": "",
            "VolumesFrom": null,
            "CapAdd": null,
            "CapDrop": null,
            "CgroupnsMode": "host",
            "Dns": [],
            "DnsOptions": [],
            "DnsSearch": [],
            "ExtraHosts": null,
            "GroupAdd": null,
            "IpcMode": "private",
            "Cgroup": "",
            "Links": null,
            "OomScoreAdj": 0,
            "PidMode": "",
            "Privileged": false,
            "PublishAllPorts": false,
            "ReadonlyRootfs": false,
            "SecurityOpt": null,
            "UTSMode": "",
            "UsernsMode": "",
            "ShmSize": 67108864,
            "Runtime": "runc",
            "ConsoleSize": [
                0,
                0
            ],
            "Isolation": "",
            "CpuShares": 0,
            "Memory": 0,
            "NanoCpus": 0,
            "CgroupParent": "",
            "BlkioWeight": 0,
            "BlkioWeightDevice": [],
            "BlkioDeviceReadBps": null,
            "BlkioDeviceWriteBps": null,
            "BlkioDeviceReadIOps": null,
            "BlkioDeviceWriteIOps": null,
            "CpuPeriod": 0,
            "CpuQuota": 0,
            "CpuRealtimePeriod": 0,
            "CpuRealtimeRuntime": 0,
            "CpusetCpus": "",
            "CpusetMems": "",
            "Devices": [],
            "DeviceCgroupRules": null,
            "DeviceRequests": null,
            "KernelMemory": 0,
            "KernelMemoryTCP": 0,
            "MemoryReservation": 0,
            "MemorySwap": 0,
            "MemorySwappiness": null,
            "OomKillDisable": false,
            "PidsLimit": null,
            "Ulimits": null,
            "CpuCount": 0,
            "CpuPercent": 0,
            "IOMaximumIOps": 0,
            "IOMaximumBandwidth": 0,
            "MaskedPaths": [
                "/proc/asound",
                "/proc/acpi",
                "/proc/kcore",
                "/proc/keys",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/proc/scsi",
                "/sys/firmware"
            ],
            "ReadonlyPaths": [
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger"
            ]
        },
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566-init/diff:/var/lib/docker/overlay2/a91889001beac8656f7cb41c959c91a5f63f1e3d2099eca45e94c64e14c83955/diff:/var/lib/docker/overlay2/f6de4c86c92476c2fa932f956d6cca0aca6b7e4ae5fbb4e08ec4ba8ee45a13cd/diff:/var/lib/docker/overlay2/c59d93e90e49ca7af1ff7787c49659ddf87b9bba4debfb6cea3b0301b286f1b6/diff:/var/lib/docker/overlay2/263c4c68745a5d2668f97d7a14ddc52e5c7a30e1bb5dfdccb6c210ccd28fcb4d/diff",
                "MergedDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/merged",
                "UpperDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/diff",
                "WorkDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/work"
            },
            "Name": "overlay2"
        },
        "Mounts": [],
        "Config": {
            "Hostname": "76f0fb3e84f2",
            "Domainname": "",
            "User": "",
            "AttachStdin": true,
            "AttachStdout": true,
            "AttachStderr": true,
            "ExposedPorts": {
                "80/tcp": {}
            },
            "Tty": true,
            "OpenStdin": true,
            "StdinOnce": true,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            ],
            "Cmd": [
                "bash"
            ],
            "Image": "ubuntu:net",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": null,
            "OnBuild": null,
            "Labels": {}
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "03bd28a2ebaac66afcf92ce6000143ef135af49e64104660384af7b1a1c11a3e",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {
                "80/tcp": [
                    {
                        "HostIp": "0.0.0.0",
                        "HostPort": "8080"
                    }
                ]
            },
            "SandboxKey": "/var/run/docker/netns/03bd28a2ebaa",
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "",
            "Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "",
            "IPPrefixLen": 0,
            "IPv6Gateway": "",
            "MacAddress": "",
            "Networks": {
                "mynet": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": [
                        "76f0fb3e84f2"
                    ],
                    "NetworkID": "0fa7b38a868b3f8790a9ee607740add2a0e2956213825f635a8fb74c5c256d80",
                    "EndpointID": "8df7abc0a436fc105116191eeb5d86aa24d980cf7399e24e20d1705828a2d884",
                    "Gateway": "172.20.0.1",
                    "IPAddress": "172.20.0.2",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:14:00:02",
                    "DriverOpts": null
                }
            }
        }
    }
]
[lyle@localhost ~]$ docker network connect bridge u1
[lyle@localhost ~]$ docker network disconnect mynet u1
```



## 容器监控与与日志

### 检查docker的状态

```
[lyle@localhost ~]$ docker info
Client:
 Context:    default
 Debug Mode: false
 Plugins:
  app: Docker App (Docker Inc., v0.9.1-beta3)
  buildx: Build with BuildKit (Docker Inc., v0.5.0-docker)

Server:
 Containers: 6
  Running: 4
  Paused: 0
  Stopped: 2
 Images: 6
 Server Version: 20.10.1
 Storage Driver: overlay2
  Backing Filesystem: xfs
  Supports d_type: true
  Native Overlay Diff: true
 Logging Driver: json-file
 Cgroup Driver: cgroupfs
 Cgroup Version: 1
 Plugins:
  Volume: local
  Network: bridge host ipvlan macvlan null overlay
  Log: awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog
 Swarm: inactive
 Runtimes: io.containerd.runtime.v1.linux runc io.containerd.runc.v2
 Default Runtime: runc
 Init Binary: docker-init
 containerd version: 269548fa27e0089a8b8278fc4fc781d7f65a939b
 runc version: ff819c7e9184c13b7c2607fe6c30ae19403a7aff
 init version: de40ad0
 Security Options:
  seccomp
   Profile: default
 Kernel Version: 3.10.0-1062.1.2.el7.x86_64
 Operating System: CentOS Linux 7 (Core)
 OSType: linux
 Architecture: x86_64
 CPUs: 4
 Total Memory: 7.62GiB
 Name: localhost.localdomain
 ID: 4BUE:AP4Y:MMAX:X66J:XSOI:2FTK:PK43:ESF7:BYUP:VD4N:A4XD:23N7
 Docker Root Dir: /var/lib/docker
 Debug Mode: false
 Registry: https://index.docker.io/v1/
 Labels:
 Experimental: false
 Insecure Registries:
  127.0.0.0/8
 Registry Mirrors:
  http://hub-mirror.c.163.com/
 Live Restore Enabled: false

```



### 查看容器内进程

```
[lyle@localhost ~]$ docker top 76f0
UID                 PID                 PPID                C                   STIME               TTY                 TIME                CMD
root                24999               24978               0                   23:05               pts/0    
```



### 容器详细信息

```
[lyle@localhost ~]$ docker inspect 76f0
[
    {
        "Id": "76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679",
        "Created": "2020-12-28T15:05:35.811004541Z",
        "Path": "bash",
        "Args": [],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false,
            "Restarting": false,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 24999,
            "ExitCode": 0,
            "Error": "",
            "StartedAt": "2020-12-28T15:05:37.023891145Z",
            "FinishedAt": "0001-01-01T00:00:00Z"
        },
        "Image": "sha256:7f0075d6104ff9d260b9cc07883cafcc3311c2e9a1af3500bc222c22e103a038",
        "ResolvConfPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/hostname",
        "HostsPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/hosts",
        "LogPath": "/var/lib/docker/containers/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679/76f0fb3e84f2fa9f306a5bbfe83a3ff9d52040b5991412d474298d23e2401679-json.log",
        "Name": "/u1",
        "RestartCount": 0,
        "Driver": "overlay2",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": null,
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {}
            },
            "NetworkMode": "mynet",
            "PortBindings": {
                "80/tcp": [
                    {
                        "HostIp": "",
                        "HostPort": "8080"
                    }
                ]
            },
            "RestartPolicy": {
                "Name": "no",
                "MaximumRetryCount": 0
            },
            "AutoRemove": true,
            "VolumeDriver": "",
            "VolumesFrom": null,
            "CapAdd": null,
            "CapDrop": null,
            "CgroupnsMode": "host",
            "Dns": [],
            "DnsOptions": [],
            "DnsSearch": [],
            "ExtraHosts": null,
            "GroupAdd": null,
            "IpcMode": "private",
            "Cgroup": "",
            "Links": null,
            "OomScoreAdj": 0,
            "PidMode": "",
            "Privileged": false,
            "PublishAllPorts": false,
            "ReadonlyRootfs": false,
            "SecurityOpt": null,
            "UTSMode": "",
            "UsernsMode": "",
            "ShmSize": 67108864,
            "Runtime": "runc",
            "ConsoleSize": [
                0,
                0
            ],
            "Isolation": "",
            "CpuShares": 0,
            "Memory": 0,
            "NanoCpus": 0,
            "CgroupParent": "",
            "BlkioWeight": 0,
            "BlkioWeightDevice": [],
            "BlkioDeviceReadBps": null,
            "BlkioDeviceWriteBps": null,
            "BlkioDeviceReadIOps": null,
            "BlkioDeviceWriteIOps": null,
            "CpuPeriod": 0,
            "CpuQuota": 0,
            "CpuRealtimePeriod": 0,
            "CpuRealtimeRuntime": 0,
            "CpusetCpus": "",
            "CpusetMems": "",
            "Devices": [],
            "DeviceCgroupRules": null,
            "DeviceRequests": null,
            "KernelMemory": 0,
            "KernelMemoryTCP": 0,
            "MemoryReservation": 0,
            "MemorySwap": 0,
            "MemorySwappiness": null,
            "OomKillDisable": false,
            "PidsLimit": null,
            "Ulimits": null,
            "CpuCount": 0,
            "CpuPercent": 0,
            "IOMaximumIOps": 0,
            "IOMaximumBandwidth": 0,
            "MaskedPaths": [
                "/proc/asound",
                "/proc/acpi",
                "/proc/kcore",
                "/proc/keys",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/proc/scsi",
                "/sys/firmware"
            ],
            "ReadonlyPaths": [
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger"
            ]
        },
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566-init/diff:/var/lib/docker/overlay2/a91889001beac8656f7cb41c959c91a5f63f1e3d2099eca45e94c64e14c83955/diff:/var/lib/docker/overlay2/f6de4c86c92476c2fa932f956d6cca0aca6b7e4ae5fbb4e08ec4ba8ee45a13cd/diff:/var/lib/docker/overlay2/c59d93e90e49ca7af1ff7787c49659ddf87b9bba4debfb6cea3b0301b286f1b6/diff:/var/lib/docker/overlay2/263c4c68745a5d2668f97d7a14ddc52e5c7a30e1bb5dfdccb6c210ccd28fcb4d/diff",
                "MergedDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/merged",
                "UpperDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/diff",
                "WorkDir": "/var/lib/docker/overlay2/f53df6396d2bd0d51c99ea53002a5dc057ec5fc8bad08e32b0cd661ed12e5566/work"
            },
            "Name": "overlay2"
        },
        "Mounts": [],
        "Config": {
            "Hostname": "76f0fb3e84f2",
            "Domainname": "",
            "User": "",
            "AttachStdin": true,
            "AttachStdout": true,
            "AttachStderr": true,
            "ExposedPorts": {
                "80/tcp": {}
            },
            "Tty": true,
            "OpenStdin": true,
            "StdinOnce": true,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
            ],
            "Cmd": [
                "bash"
            ],
            "Image": "ubuntu:net",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": null,
            "OnBuild": null,
            "Labels": {}
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "03bd28a2ebaac66afcf92ce6000143ef135af49e64104660384af7b1a1c11a3e",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {
                "80/tcp": [
                    {
                        "HostIp": "0.0.0.0",
                        "HostPort": "8080"
                    }
                ]
            },
            "SandboxKey": "/var/run/docker/netns/03bd28a2ebaa",
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "70c00876efc3a88aa4a2d0876c0159138f4fa79c3f4b73c2b12742ca17325bd9",
            "Gateway": "172.17.0.1",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "172.17.0.4",
            "IPPrefixLen": 16,
            "IPv6Gateway": "",
            "MacAddress": "02:42:ac:11:00:04",
            "Networks": {
                "bridge": {
                    "IPAMConfig": {},
                    "Links": null,
                    "Aliases": [],
                    "NetworkID": "869a959200eafc36ecd669ba5e06e73050252e44bdfb4490cf8e3fbcef888973",
                    "EndpointID": "70c00876efc3a88aa4a2d0876c0159138f4fa79c3f4b73c2b12742ca17325bd9",
                    "Gateway": "172.17.0.1",
                    "IPAddress": "172.17.0.4",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:11:00:04",
                    "DriverOpts": {}
                }
            }
        }
    }
]

```



### 容器日志查看

```
[lyle@localhost ~]$ docker logs 7691906f6eb6
root@7691906f6eb6:/# ^C
root@7691906f6eb6:/# ls
bin   dev  home  lib32  libx32  mnt  proc  run   srv  tmp  var
boot  etc  lib   lib64  media   opt  root  sbin  sys  usr

```

## docker 仓库

登陆

```
[root@localhost lyle]# docker login --username=小罗在养生 registry.cn-shenzhen.aliyuncs.com
Password: 
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded

```

标签

```
[root@localhost lyle]# docker tag hello-world registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:helloworld
```

下载

```
[root@localhost lyle]# docker pull registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:hello-world
hello-world: Pulling from pmlpml/repo
983bfa07a342: Pull complete 
Digest: sha256:2075ac87b043415d35bb6351b4a59df19b8ad154e578f7048335feeb02d0f759
Status: Downloaded newer image for registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:hello-world
registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:hello-world
```

删除

```
[root@localhost lyle]#  docker rmi registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:hello-world
Untagged: registry.cn-shenzhen.aliyuncs.com/pmlpml/repo:hello-world
Untagged: registry.cn-shenzhen.aliyuncs.com/pmlpml/repo@sha256:2075ac87b043415d35bb6351b4a59df19b8ad154e578f7048335feeb02d0f759
Deleted: sha256:48b5124b2768d2b917edcb640435044a97967015485e812545546cbed5cf0233
Deleted: sha256:98c944e98de8d35097100ff70a31083ec57704be0991a92c51700465e4544d08
```

运行

```
[root@localhost lyle]#  docker run --rm hello-world

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

退出

```
[root@localhost lyle]#  docker logout registry.cn-shenzhen.aliyuncs.com
Removing login credentials for registry.cn-shenzhen.aliyuncs.com
```



