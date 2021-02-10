package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"net/http"
)

var (
	Router           *mux.Router
)

func init()  {
	Router = mux.NewRouter()
	Router.Use(contextMiddleware)
	Router.HandleFunc("/", Hello)
}

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Router.ServeHTTP(w, r)
	return
}

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request-id", id)
		r = r.WithContext(ctx)
		w.Header().Set("request-id", id)
		next.ServeHTTP(w, r)
		return
	})
}