package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jesses-code-adventures/slop/db"
	"github.com/jesses-code-adventures/slop/jwt"
	"github.com/jesses-code-adventures/slop/server"
)

type AppApi struct {
	server *server.Server
	logger *slog.Logger
	db     db.Database
}

func NewAppApi(s *server.Server, l *slog.Logger, d db.Database) AppApi {
	return AppApi{server: s, logger: l, db: d}
}

func (a AppApi) Server() *server.Server {
	return a.server
}

func (a AppApi) UserAuthorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Authenticating request")
		token := a.preferCookieOverHeader(r, "Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userId := a.preferCookieOverHeader(r, "x-slop-user-id")
		if userId == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		newToken, err := jwt.ValidateToken(userId, token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Error(w, "Invalid Auth Token", http.StatusUnauthorized)
			return
		}
		a.setAuthorization(w, r, newToken)
		next(w, r)
	}
}

type userNewFromRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}

func (a AppApi) UserRegister(w http.ResponseWriter, r *http.Request) {
	var user userNewFromRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	userId, err := a.db.UserRegister(user.FirstName, user.LastName, user.Email, user.RawPassword)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}
	a.logger.Info(fmt.Sprintf("User registered successfully"))
	w.WriteHeader(http.StatusCreated)
	a.jsonResponse(w, r, "id", userId)
}

type userLoginFromRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a AppApi) UserLogin(w http.ResponseWriter, r *http.Request) {
	var user userLoginFromRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// TODO: password comparison, argon2 etc
	userId, err := a.db.UserLogin(user.Email, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := jwt.CreateToken(userId)
	if err != nil {
		http.Error(w, "Error creating auth token", http.StatusInternalServerError)
		return
	}
	a.logger.Info("User logged in successfully")
	a.jsonResponse(w, r, "token", token)
}

type ImageRequestMetadata struct {
	Character string
}

func (a AppApi) ImageGenerate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	metadata := ImageRequestMetadata{
		Character: r.FormValue("character"),
	}
	// Process the image and metadata as needed
	a.logger.Info("Image generation process started", "character", metadata.Character)
	// Simulate image processing and generation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image generated successfully"))
}
