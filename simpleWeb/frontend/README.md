# 博客前端
## 环境
- nodejs
- yarn

## 部署方式

前端利用 yarn 作为包管理器，首先确保安装了 node 14 或以上版本，然后安装：

```bash
npm install -g yarn
yarn
```

由于网络问题，可能无法成功还原，此时可以解压doc仓库中的 node_modules 压缩包到项目目录中，然后对 `node_modules/.bin` 中的所有文件设置可执行权限：

```bash
chmod +x node_modules/.bin/*
```

最后运行前端项目：

```bash
yarn start
```

然后访问 `http://localhost:3000`即可。

## 编写指南
[Guide.md](Guide.md)