package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//允许跨域访问
	r.Use(CrosHandler())

	//设置分组路由
	v1 := r.Group("/user")

	//根据分组设置路由
	{
		v1.GET("/name", GetUserName)
		v1.POST("/login", Login)
	}

	//启动
	r.Run() // listen and serve on 0.0.0.0:8080

}

func GetUserName(c *gin.Context) {
	c.String(http.StatusOK, "Faker")
}

func Login(c *gin.Context) {
	fmt.Println(c.FullPath())
	username, nameValid := c.GetPostForm("username")
	password, passwordValid := c.GetPostForm("password")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	if !nameValid || !passwordValid {
		fmt.Println("get param fail")
	} else {
		fmt.Println(username)
		fmt.Println(password)
	}
}

//跨域访问：cross  origin resource share
func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json") //设置返回格式是json

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "OK")
		}

		//处理请求
		context.Next()
	}
}
