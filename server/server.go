package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	var err error
	fmt.Println("connecting to mysql")
	Db, err = sql.Open("mysql", "root:ccj6685533@tcp(172.19.9.94:3306)/go")

	err = Db.Ping()
	if err != nil {
		fmt.Println("connect error")
	} else {
		fmt.Println("connected to mysql")
	}

	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)

}

func main() {

	// Init()
	r := gin.Default()

	//允许跨域访问
	r.Use(CrosHandler())

	//设置分组路由
	v1 := r.Group("/user")

	//根据分组设置路由
	{
		v1.GET("/name", GetUserName)
		v1.POST("/login", Login)
		v1.POST("/register", Register)
	}

	//启动
	r.Run() // listen and serve on 0.0.0.0:8080

}

func GetUserName(c *gin.Context) {
	c.String(http.StatusOK, "Faker")
}

type registerModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(c *gin.Context) {
	var registerInfo registerModel
	c.Bind(&registerInfo)

	result, err := Db.Exec("insert into user (username, password, email) values (?,?,?)", registerInfo.UserName, registerInfo.Password, registerInfo.Email)
	if err != nil {
		fmt.Println("err:%s", err)
	} else {
		fmt.Println("result:%s", result)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type loginModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var loginInfo loginModel
	c.Bind(&loginInfo)

	id := "not defined"
	status := "not defined"
	err := Db.QueryRow("select user_id from user where username = ? and password = ?", loginInfo.UserName, loginInfo.Password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows { //如果未查询到对应字段则...
			status = "not found"
			fmt.Println("not found")
		} else {
			status = "failure"
			fmt.Println("failue")
			log.Fatal(err)
		}
	} else {
		status = "success"
		fmt.Println("found")
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
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
