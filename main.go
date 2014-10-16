package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// Handle the routes
	r := mux.NewRouter()

	// Student routes
	r.HandleFunc("/api/auth/login", addDefaultHeaders(StudentLogin))
	r.HandleFunc("/api/auth/register", addDefaultHeaders(StudentRegister))

	// Server this on /api and use nginx to proxy
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
