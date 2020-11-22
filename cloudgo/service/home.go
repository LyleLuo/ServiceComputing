package service

import (
	"net/http"
	"text/template"

	"github.com/unrolled/render"
)

func homeHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "index", nil)
	}
}

func checkform(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username, password := template.HTMLEscapeString(r.Form.Get("username")), template.HTMLEscapeString(r.Form.Get("password"))
	template.Must(template.New("login.html").ParseFiles("templates/login.html")).Execute(w, struct {
		Username string
		Password string
	}{Username: username, Password: password})
}

func apiTestHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{Username: "用户名: ", Password: "密码: &emsp;"})
	}
}
