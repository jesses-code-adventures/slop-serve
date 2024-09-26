package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jesses-code-adventures/slop/auth"
	"github.com/jesses-code-adventures/slop/db"
	"github.com/jesses-code-adventures/slop/jwt"
	"github.com/jesses-code-adventures/slop/server"
)

type AppApi struct {
	server *server.Server
	db     db.Database
	ctx    context.Context
	auth   auth.Authorizer
}

func NewAppApi(s *server.Server, ctx context.Context, d db.Database) AppApi {
	return AppApi{server: s, ctx: ctx, db: d, auth: auth.NewPasswordHandler()}
}

func (a *AppApi) logger() *slog.Logger {
	return a.ctx.Value("logger").(*slog.Logger)
}

func (a AppApi) Server() *server.Server {
	return a.server
}

func (a AppApi) UserAuthorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger().Info("Authenticating request")
		token := a.preferCookieOverHeader(r, "Authorization")
		if token == "" {
			a.logger().Error("Found no token")
			http.Error(w, a.jsonErrorString("Unauthorized"), http.StatusUnauthorized)
			return
		}
		a.logger().Debug(fmt.Sprintf("got token %s", token))
		userId := a.preferCookieOverHeader(r, "x-slop-user-id")
		if userId == "" {
			a.logger().Error("couldn't get x-slop-user-id")
			http.Error(w, a.jsonErrorString("Unauthorized"), http.StatusUnauthorized)
			return
		}
		newToken, err := jwt.ValidateToken(userId, token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Error(w, a.jsonErrorString("Invalid Auth Token"), http.StatusUnauthorized)
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

func (u userNewFromRequest) Print() {
	fmt.Printf("FirstName: %s\n", u.FirstName)
	fmt.Printf("LastName: %s\n", u.LastName)
	fmt.Printf("Email: %s\n", u.Email)
	fmt.Printf("RawPassword: %s\n", u.RawPassword)
}

func (a AppApi) UserRegister(w http.ResponseWriter, r *http.Request) {
	var user userNewFromRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, a.jsonErrorString("Invalid input"), http.StatusBadRequest)
		return
	}
	hashed, err := a.auth.Hash(user.RawPassword)
	if err != nil {
		a.logger().Error("got error hashing password")
		http.Error(w, a.jsonErrorString("Error processing password"), http.StatusInternalServerError)
		return
	}
	userId, err := a.db.UserCreate(user.FirstName, user.LastName, user.Email, hashed)
	if err != nil {
		a.logger().Error("got error registering user %s", err.Error())
		http.Error(w, a.jsonErrorString("Error registering user"), http.StatusInternalServerError)
		return
	}
	a.logger().Info(fmt.Sprintf("User registered successfully"))
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
		http.Error(w, a.jsonErrorString("Invalid input"), http.StatusBadRequest)
		return
	}
	userId, hashedPassword, err := a.db.HashedPasswordGet(user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, a.jsonErrorString("Email Not Found"), http.StatusNotFound)
			return
		} else {
			http.Error(w, a.jsonErrorString("Internal Server Error"), http.StatusInternalServerError)
			return
		}
	}
	authorized, err := a.auth.Verify(user.Password, hashedPassword)
	if err != nil || !authorized {
		http.Error(w, a.jsonErrorString("Not Authorized"), http.StatusUnauthorized)
		return
	}
	token, err := jwt.CreateToken(userId)
	if err != nil {
		http.Error(w, a.jsonErrorString("Error creating auth token"), http.StatusInternalServerError)
		return
	}
	a.setAuthorization(w, r, token)
	a.logger().Info("User logged in successfully")
	a.jsonResponse(w, r, "token", token)
}

type ImageRequestMetadata struct {
	Character string
}

func (a AppApi) ImageGenerate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, a.jsonErrorString("Unable to parse form"), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, a.jsonErrorString("Error retrieving file"), http.StatusBadRequest)
		return
	}
	defer file.Close()
	metadata := ImageRequestMetadata{
		Character: r.FormValue("character"),
	}
	a.logger().Info("Image generation process started", "character", metadata.Character)
	// Simulate image processing and generation
	a.jsonResponse(w, r, "message", "Image generated successfully")
}
