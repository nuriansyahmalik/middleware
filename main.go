package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(LoggingMiddleware)
	r.Route("/middleware", func(r chi.Router) {
		r.Use(LoggingMiddleware, func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("received request", r.Method, r.URL.Path)
				h.ServeHTTP(w, r)
			})
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("middleware"))
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Custom Middleware"))
	})
	http.ListenAndServe(":8080", r)
}

var xApiKey = "apikey"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request", r.Method, r.URL.Path)
		apikey := r.Header.Get("X-API-KEY")
		if apikey != xApiKey {
			fmt.Println("failed x api key")
		}
		next.ServeHTTP(w, r)
	})
}
