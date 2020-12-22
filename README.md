# 服务计算小组项目文档
## 部署和运行
### 数据库部署
数据库采用 MySQL 8.0，建表脚本：

```sql
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `blog`
--

DROP TABLE IF EXISTS `blog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog` (
  `blog_id` int(11) NOT NULL AUTO_INCREMENT,
  `author_id` int(11) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `text` text,
  PRIMARY KEY (`blog_id`)
) /*!50100 TABLESPACE `innodb_system` */ ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tag` (
  `tag_id` int(11) NOT NULL AUTO_INCREMENT,
  `tag_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`tag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag_blog`
--

DROP TABLE IF EXISTS `tag_blog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tag_blog` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tag_id` int(11) DEFAULT NULL,
  `blog_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
```

利用上述 SQL 即可完成数据库表结构的创建。

### 前端
前端利用 yarn 作为包管理器，首先确保安装了 node 14 或以上版本，然后安装：

```bash
npm install -g yarn
yarn
```

由于网络问题，可能无法成功还原，此时可以解压本仓库中的 node_modules 压缩包到项目目录中，然后对 `node_modules/.bin` 中的所有文件设置可执行权限：

```bash
chmod +x node_modules/.bin/*
```

最后运行前端项目：

```bash
yarn start
```

即可。

### 后端
首先运行 install 安装依赖项：
```bash
./install.sh
```

然后进入 `server` 目录，修改数据库连接字符串后运行即可：

```bash
cd server
nano server.go # 此步修改源代码中数据库的连接字符串：Db, err = sql.Open("mysql", "......")
go run server.go
```
如果不想使用nano命令的话，去`server.go`找对应的数据库配置代码（23行）用下面代码改变连接的数据库：

```
Db, err = sql.Open("数据库类型", "用户名:密码@tcp(数据库IP:端口)/数据库名称")
```

## API 原型说明
### Login `POST /user/login`
使用该 API 可以登录用户账户。

```
Content-Type: application/json
Method: POST
Body: 
{ 
    "username": string, 
    "password": string
}
Response: 
{ 
    "status": string
}
```
### Register `POST /user/register`
使用该 API 可以注册用户账户。

```
Content-Type: application/json
Method: POST
Body: 
{ 
    "username": string, 
    "password": string, 
    "email": string
} 
Response: 
{ 
    "status": string 
}
```
### Self `GET /user/self`
使用该 API 可以获取当前已登录用户的信息。

```
Content-Type: application/json
Method: GET
Body: Empty
Response: 
{ 
    "name": string, 
    "email": string, 
    "id": number 
}
```
### Logout `POST /user/logout`
使用该 API 可以获取退出当前用户。

```
Content-Type: application/json
Method: POST
Body: Empty
Response: 
{ 
    "status": string
}
```
### Post `POST /user/post`
使用该 API 可以发布博客。

```
Content-Type: application/json
Method: POST
Body: 
{ 
    "author_id": string, 
    "tags": string[], 
    "text": string, 
    "title": string
}
Response: 
{ 
    "status": string,
    "Data": { 
        "blog_id": number
    }
}
```
### Tags `GET /user/tags`
使用该 API 可以获取所有的标签和博客列表。

```
Content-Type: application/json
Method: GET
Body: Empty
Response: 
{ 
    "status": string, 
    "tags": { 
        "id": number, 
        "tagname": string 
    }[], 
    "blogs": { 
        "id": number, 
        "author": string, 
        "title": string, 
        "tags": string[]
    }[], 
    "count": number
}
```
### Portal `POST /user/portal`
使用该 API 可以获取指定用户的博客列表。

```
Content-Type: application/json
Method: POST
Body: 
{
    "author_id": number
}
Response:
{
    "result": {
        "id": number,
        "title": string
    }[],
    "count": number
}
```
### Page `GET /page/:id`
使用该 API 可以获取指定页的博客列表。

```
Content-Type: application/json
Method: GET
Body: Empty
Response:
{
    "status": string,
    "data": {
        "blog_id": string,
        "title": string,
        "username": string
    }[]
}
```
### Details `GET /details/:id`
使用该 API 可以获取指定 ID 的博客。

```
Content-Type: application/json
Method: GET
Body: Empty
Response:
{
    "status": string,
    "data": {
        "author": string,
        "authorEmail": string,
        "authorId": number,
        "id": number,
        "tags": string[],
        "text": string,
        "title": string
    }
}
```

## API 请求实例
### Register
使用该 API 注册一个用户。

请求：
```
curl 'http://localhost:8080/user/register' -H 'Content-Type: application/json' --data-binary '{"username":"testuser","password":"123456","email":"someone@example.com"}'
```

响应：
```
{"status":"success"}
```

### Page
使用该 API 获取指定页的博客列表。

请求：
```
curl 'http://localhost:8080/page/1'
```

响应：
```
{"data":[{"blog_id":"12","title":"test","username":"test"},{"blog_id":"11","title":"Markdown: Syntax","username":"test"},{"blog_id":"10","title":"blablablabla","username":"test"}],"status":"success"}
```
