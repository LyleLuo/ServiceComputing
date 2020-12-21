package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	var err error
	fmt.Println("connecting to mysql")
	// Db, err = sql.Open("mysql", "mysql@sc-database:ServiceComputing2020@tcp(sc-database.mysql.database.azure.com:3306)/go")
	Db, err = sql.Open("mysql", "root:111111@tcp(172.26.28.10:3306)/go")

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
		v1.POST("/login", Login)
		v1.POST("/register", Register)
		v1.GET("/self", Self)
		v1.POST("/logout", Logout)
		v1.POST("/post", Post)
		v1.GET("/tags", Tags)
		v1.POST("/portal", MY)
		// v1.GET("/getTages", GetTages)
	}

	r.GET("/page/:id", Page)
	//启动
	r.Run() // listen and serve on 0.0.0.0:8080

}

var (
	Secret     = "service-computing-blogs"
	ExpireTime = 3600
)

type JWTClaims struct {
	jwt.StandardClaims
	UserId   int64  `json:"userid"`
	Password string `json:"password"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func getToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func verifyToken(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, err
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}

type registerModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(c *gin.Context) {
	var registerInfo registerModel
	c.Bind(&registerInfo)

	result, err := Db.Exec("insert into user (username, password, email) values (?,?,?);", registerInfo.UserName, registerInfo.Password, registerInfo.Email)
	var id int64
	if err != nil {
		fmt.Println("err:%s", err)
	} else {
		id, _ = result.LastInsertId()
	}

	claims := &JWTClaims{
		UserId:   id,
		UserName: registerInfo.UserName,
		Password: registerInfo.Password,
		Email:    registerInfo.Email,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	signedToken, err := getToken(claims)
	c.SetCookie("jwt-token", signedToken, 3600, "/", ".", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type loginModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt-token", "", 0, "/", ".", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type blogModel struct {
	title     string   `json:"title"`
	author_id int      `json:"author_id"`
	tags      []string `json:"tags"`
	text      string   `json:text`
}

func Post(c *gin.Context) {
	// var blogInfo blogModel
	// c.ShouldBind(&blogInfo)
	// fmt.Println("title:", blogInfo.title)
	// fmt.Println("author:", blogInfo.author)
	// fmt.Println("tags:", blogInfo.tags)
	// fmt.Println("text:", blogInfo.text)

	data, _ := ioutil.ReadAll(c.Request.Body)
	blogInfo := make(map[string]interface{})
	if err := json.Unmarshal(data, &blogInfo); err != nil {
		fmt.Println("error")
	}

	title := blogInfo["title"]
	author_id := blogInfo["author_id"]
	tags := blogInfo["tags"].([]interface{})
	text := blogInfo["text"]

	fmt.Println("title:", title)
	fmt.Println("author_id:", author_id)
	fmt.Print("tags:")

	for i := 0; i < len(tags); i++ {
		fmt.Print(" ", tags[i])
	}
	fmt.Println()
	fmt.Println("text:", text)

	if author_id.(float64) <= 0 {
		fmt.Println("not login")
		c.JSON(http.StatusOK, gin.H{
			"status": "not login",
			"Data":   nil,
		})
		return
	}

	result, err := Db.Exec("insert into blog (author_id, title, text) values (?,?,?);", author_id, title, text)
	var blog_id int64
	if err != nil {
		fmt.Println("err:%s", err)
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"Data":   nil,
		})
		return
	} else {
		blog_id, _ = result.LastInsertId()
	}

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		var tag_id int64
		err = Db.QueryRow("select tag_id from tag where tag_name = ?", tag).Scan(&tag_id)
		if err != nil {
			if err == sql.ErrNoRows { //如果未查询到对应字段则...
				result, err = Db.Exec("insert into tag (tag_name) values (?);", tag)
				fmt.Println("tag not found, insert into databse")
				if err != nil {
					fmt.Println("err:%s", err)
					c.JSON(http.StatusOK, gin.H{
						"status": "failure",
						"Data":   nil,
					})
					return
				} else {
					tag_id, _ = result.LastInsertId()
				}
			} else {
				fmt.Println("failue")
				log.Fatal(err)
				c.JSON(http.StatusOK, gin.H{
					"status": "failure",
					"Data":   nil,
				})
				return
			}
		}
		result, err = Db.Exec("insert into tag_blog (tag_id, blog_id) values (?,?);", tag_id, blog_id)
		if err != nil {
			fmt.Println("err%s", err)
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"Data":   nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"Data": map[string]interface{}{
			"blog_id": blog_id,
		},
	})

}

func Login(c *gin.Context) {
	var loginInfo loginModel
	c.Bind(&loginInfo)

	var id int64
	var email string
	status := "not defined"
	err := Db.QueryRow("select user_id, email from user where username = ? and password = ?", loginInfo.UserName, loginInfo.Password).Scan(&id, &email)
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
	}

	claims := &JWTClaims{
		UserId:   id,
		UserName: loginInfo.UserName,
		Password: loginInfo.Password,
		Email:    email,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	signedToken, _ := getToken(claims)
	c.SetCookie("jwt-token", signedToken, 3600, "/", ".", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func Self(c *gin.Context) {
	strToken, err := c.Cookie("jwt-token")
	claims, err := verifyToken(strToken)
	if err != nil {
		c.String(401, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.SetCookie("jwt-token", signedToken, 3600, "/", ".", false, true)

	c.JSON(http.StatusOK, gin.H{
		"name":  claims.UserName,
		"email": claims.Email,
		"id":    claims.UserId,
	})
}

type Tag struct {
	Id      int    `json:"id"`
	Tagname string `json:"tagname"`
}

func getTags() (tags []Tag, err error) {
	rows, err := Db.Query("SELECT tag_id,tag_name FROM tag")
	for rows.Next() {
		var tag1 Tag
		//遍历表中所有行的信息
		rows.Scan(&tag1.Id, &tag1.Tagname)
		//将user添加到users中
		tags = append(tags, tag1)
	}
	return
}

func Tags(c *gin.Context) {
	tags, err := getTags()
	if err != nil {
		log.Fatal(err)
	}
	//H is a shortcut for map[string]interface{}
	c.JSON(http.StatusOK, gin.H{
		"result": tags,
		"count":  len(tags),
	})
}

type Blog struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func getMY(user_id int) (blogs []Blog, err error) {
	print("%d", user_id)
	rows, err := Db.Query("SELECT blog_id,title FROM blog WHERE author_id = ?", 2)
	//print("%d", len(rows))
	for rows.Next() {
		var blog Blog
		//遍历表中所有行的信息
		rows.Scan(&blog.Id, &blog.Title)
		//将user添加到users中
		blogs = append(blogs, blog)
	}
	return
}

type my_id struct {
	Author_id int `json:"author_id"`
}

func MY(c *gin.Context) {
	author_id := c.PostForm("author_id")
	print(author_id)
	i, err := strconv.Atoi(author_id)
	blogs, err := getMY(i)
	if err != nil {
		log.Fatal(err)
	}
	//H is a shortcut for map[string]interface{}
	c.JSON(http.StatusOK, gin.H{
		"result": blogs,
		"count":  len(blogs),
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

func getJSON(sqlString string) ([]map[string]interface{}, error) {
	rows, err := Db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func Page(c *gin.Context) {
	id, covErr := strconv.Atoi(c.Param("id"))
	result, err := getJSON("select blog_id, username, title from blog, user where blog.author_id=user.user_id;")
	status := "undefined"
	itemPerPage := 3
	start := (id - 1) * itemPerPage
	last := id * itemPerPage
	if id*itemPerPage > len(result) {
		last = len(result)
	}

	if err != nil || covErr != nil || id < 1 || (id-1)*itemPerPage >= len(result) {
		status = "fail"
		start, last = 0, 0
	} else {
		status = "success"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   result[start:last],
	})
}
