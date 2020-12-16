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
	Db, err = sql.Open("mysql", "root:123@tcp(172.18.43.166:3306)/mysql")

	err = Db.Ping()
	if err != nil {
		print("connect error")
	}
	if err != nil {
		print("error")
	} else {
		print("success")
	}

	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
	var name string
	err = Db.QueryRow("select name from test where test_id = ?", "01").Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows { //如果未查询到对应字段则...
			print("no rows")
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(name)
	//insert
	stmt, err1 := Db.Prepare("INSERT INTO test SET name=?")
	res, err2 := stmt.Exec("test2")
	if err1 != nil {
		print("error1")
	} else {
		print("success")
	}
	if err2 != nil {
		print("error2", err2)
	} else {
		print("success")
	}
	if res != nil {
		print("error")
	} else {
		print("success")
	}
}

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
