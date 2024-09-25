package api

import (
	"log/slog"
	"net/http"

	"github.com/jesses-code-adventures/slop/server"
)

type TestApi struct {
	server *server.Server
	logger *slog.Logger
}

func NewTestApi(s *server.Server, l *slog.Logger) TestApi {
	return TestApi{server: s, logger: l}
}

func (a TestApi) Server() *server.Server {
	return a.server
}

func (a TestApi) UserAuthorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Would authenticate a reqest here")
		next(w, r)
	}
}

func (a TestApi) UserRegister(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("would register a user here")
}

func (a TestApi) UserLogin(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("would log a user in here")
}

func (a TestApi) ImageGenerate(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("would generate an image here")
}
