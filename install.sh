go env -w GO111MODULE=on 
go env -w GOPROXY=https://goproxy.cn,direct
go get -u github.com/gin-gonic/gin
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/go-sql-driver/mysql
