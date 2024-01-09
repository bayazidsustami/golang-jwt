package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/novalagung/gubrak/v2"
)

type CustomMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

type MyClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Group    string `group:"group"`
}

func (c *CustomMux) RegisterMiddleware(next func(next http.Handler) http.Handler) {
	c.middlewares = append(c.middlewares, next)
}

func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &c.ServeMux

	for _, next := range c.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

type M map[string]interface{}

var APPLICATION_NAME = "My Simple JWT APP"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("rahasia banget")

func main() {
	mux := new(CustomMux)
	mux.RegisterMiddleware(MiddlewareJWTAuthorization)

	mux.HandleFunc("/login", HandlerLogin)
	mux.HandleFunc("/index", HandlerIndex)

	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":8080"

	fmt.Println("Starting Server at", server.Addr)
	server.ListenAndServe()
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported http method", http.StatusBadRequest)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	ok, userInfo := authenticateuser(username, password)
	if !ok {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(LOGIN_EXPIRATION_DURATION)),
		},
		Username: userInfo["username"].(string),
		Email:    userInfo["email"].(string),
		Group:    userInfo["group"].(string),
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, _ := json.Marshal(M{"token": signedToken})
	w.Write([]byte(tokenString))

}

func authenticateuser(username, password string) (bool, M) {
	basePath, _ := os.Getwd()
	dbPath := filepath.Join(basePath, "users.json")
	buf, _ := os.ReadFile(dbPath)

	data := make([]M, 0)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return false, nil
	}

	res := gubrak.From(data).Find(func(each M) bool {
		return each["username"] == username && each["password"] == password
	}).Result()

	if res != nil {
		resM := res.(M)
		delete(resM, "password")
		return true, resM
	}

	return false, nil
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	message := fmt.Sprintf("hello %s (%s)", userInfo["username"], userInfo["email"])
	w.Write([]byte(message))
}
