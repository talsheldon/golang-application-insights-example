package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	Router           *mux.Router
)

func init()  {
	Router = mux.NewRouter()
	Router.HandleFunc("/", Hello)
}

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Router.ServeHTTP(w, r)
	return
}
