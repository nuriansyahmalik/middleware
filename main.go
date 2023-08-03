package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(LoggingMiddleware)
	r.Group(func(r chi.Router) {
		r.Route("/middleware", func(r chi.Router) {
			r.Use(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Println("received request", r.Method, r.URL.Path)
					h.ServeHTTP(w, r)
				})
			})
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("middleware"))
			})
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Custom Middleware"))
	})

	jwtTOkens, _ := generateJwt()
	fmt.Println(jwtTOkens)
	fmt.Println(validateJwt(jwtTOkens))

	http.ListenAndServe(":8080", r)
}

var xApiKey = "apikey"
var secret = "secret"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request", r.Method, r.URL.Path)
		apikey := r.Header.Get("X-API-KEY")
		if apikey != xApiKey {
			fmt.Println("failed x api key")
		}
		w.Header().Set("X-API-KEY", xApiKey)
		next.ServeHTTP(w, r)
	})
}

type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateJwt() (string, error) {
	claim := Claims{
		UserID:   123,
		Username: "dadang ben",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "evermos",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJwt(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("JWT validation failed, %s", err)
	}
	if claim, ok := token.Claims.(*Claims); ok != token.Valid {
		return claim, nil
	}
	return nil, fmt.Errorf("JWT not valid")
}
