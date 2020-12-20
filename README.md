# 博客后端

##   环境
- go

## 安装

先赋予安装脚本权限

```
sudo chmod +x install.sh
```

执行安装：

```
./install.sh
```

## 部署方式

将该仓库放到`$GOPATH/src/github.com/user/`下，cd到`server`下执行`go run server.go`

## 注意事项

* 必须在`GO111MODULE="on"`模式下

## 数据库

* ubuntu上：
  * ip：172.26.28.10
  * 用户名：root
  * 密码：111111
* 云服务器上：
  * ip: sc-database.mysql.database.azure.com
  * 用户名：mysql
  * 密码：ServiceComputing2020

## 部署方式

* 
