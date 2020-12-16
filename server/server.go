package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//设置分组路由
	v1 := r.Group("/user")

	//根据分组设置路由
	{
		v1.GET("/name", GetUserName)
	}

	//启动
	r.Run() // listen and serve on 0.0.0.0:8080

}

func GetUserName(c *gin.Context) {
	c.String(http.StatusOK, "Faker")
}
