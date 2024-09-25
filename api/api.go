package api

import (
	"net/http"

	"github.com/jesses-code-adventures/slop/server"
)

func BindRoutes(a Api) {
	a.Server().RegisterHandler("/register", a.UserRegister)
	a.Server().RegisterHandler("/login", a.UserLogin)
	a.Server().RegisterHandler("/image", a.UserAuthorize(a.ImageGenerate))
}

type Api interface {
	Server() *server.Server
	ImageGenerate(http.ResponseWriter, *http.Request)
	UserRegister(http.ResponseWriter, *http.Request)
	UserAuthorize(http.HandlerFunc) http.HandlerFunc // middleware
	UserLogin(http.ResponseWriter, *http.Request)
}
