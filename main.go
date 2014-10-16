package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// Handle the routes
	r := mux.NewRouter()

	// Student routes
	r.HandleFunc("/auth/login", addDefaultHeaders(StudentLogin))
	r.HandleFunc("/auth/register", addDefaultHeaders(StudentRegister))

	// Server this on /api and use nginx to proxy
	http.Handle("/api/", r)
	http.ListenAndServe(":3000", nil)
}
