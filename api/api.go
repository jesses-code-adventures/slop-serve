package api

import (
	"log/slog"
	"net/http"

	"github.com/jesses-code-adventures/slop/server"
)

func BindRoutes(a Api) {
	a.Server().RegisterHandler("/register", a.UserRegister)
	a.Server().RegisterHandler("/login", a.UserLogin)
	a.Server().RegisterHandler("/image", a.UserAuthenticate(a.ImageGenerate))
}

type Api interface {
	Server() *server.Server
	ImageGenerate(http.ResponseWriter, *http.Request)
	UserRegister(http.ResponseWriter, *http.Request)
	UserAuthenticate(http.HandlerFunc) http.HandlerFunc // middleware
	UserLogin(http.ResponseWriter, *http.Request)
}

type TestApi struct {
	server *server.Server
	logger *slog.Logger
}

func NewTestApi(s *server.Server, l *slog.Logger) TestApi {
	return TestApi{server: s, logger: l}
}

func (t TestApi) Server() *server.Server {
	return t.server
}

func (t TestApi) UserAuthenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.logger.Info("Would authenticate a reqest here")
		next(w, r)
	}
}

func (t TestApi) UserRegister(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("would register a user here")
}

func (t TestApi) UserLogin(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("would log a user in here")
}

func (t TestApi) ImageGenerate(w http.ResponseWriter, r *http.Request) {
	t.logger.Info("would generate an image here")
}
