# 博客后端

## 环境
- go

## 部署方式

先赋予安装脚本权限

```
sudo chmod +x install.sh
```

执行安装：

```
./install.sh
```

## 运行服务器

将该仓库放到`$GOPATH/src/github.com/user/`下，cd到`server`下执行`go run server.go`

## 注意事项

* 必须在`GO111MODULE="on"`模式下

## 数据库

* ubuntu上：
  * ip：172.26.114.137
  * 用户名：root
  * 密码：111111
* 云服务器上：
  * ip: sc-database.mysql.database.azure.com
  * 用户名：mysql
  * 密码：ServiceComputing2020

对数据库的一些说明：

* 数据库现在部署在组员的一个ubuntu上，开机时间为早上10点到晚上1点
* ubuntu数据库的IP可能会发生改变，如果连接不上可以连接云服务器上的数据库，但是云服务器在美国，可能访问速度会比较慢（直接在`server.go`文件中解注释掉对应的那一行）
* 如果实在连不上数据库，有以下两种方法：
  * 自己配置一个数据库
  * 联系组员开恢复ubuntu数据库
