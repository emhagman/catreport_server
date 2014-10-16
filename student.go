package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
)

func StudentLogin(res http.ResponseWriter, req *http.Request) {

	// Get our POST vars
	vars := mux.Vars(req)

	// Get username and password
	username := vars["username"]
	password := vars["password"]

	// Make sure we actually got a username and password
	// In Go, strings are never null or 'nil'
	if username == "" || password == "" {
		fmt.Fprint(res, Response{"success": false})
	}

	// Automatically add @villanova.edu so people can't mess with us
	// Villanova students only, remember?
	if !strings.Contains(username, "@villanova.edu") {
		username += "@villanova.edu"
	}

	// Back before I knew anything about programming...
	// I used MD5 to hash the passwords without any salt
	// Granted, there isn't really anything intense stored in the database
	// But this is terrible practice and I am ashamed of old me
	oldPassword := OldHashPassowrd(password)

	// Make sure we have connection first
	db, err := DBConnect()
	if err != nil {
		fmt.Fprint(res, Response{"success": false})
	}

	// Go validate our user (using the old password first)
	found := false
	rows, err := db.Query("SELECT * FROM students WHERE username = $1 AND password = $2", username, oldPassword)
	for rows.Next() {
		found = true
	}

	// Check to see if old hash works
	// If the old hash did work, update the database with new hash
	if found {
		StudentLogin(res, req)
		fmt.Fprint(res, Response{"success": true, "old": true})
		return
	}

	// New password is hashed using pbkdf2 and the salt is kept on server :P
	newPassword := HashPassword([]byte(password), []byte(os.Getenv("CATREPORT_SALT")))

	// Let's check to see if they have the new hashing algo
	// I really should just convert everyone eventually...
	rows, err = db.Query("SELECT * FROM students WHERE username = $1 AND password = $2", username, newPassword)
	for rows.Next() {
		SessionLogin(res, req)
		fmt.Fprint(res, Response{"success": true, "new": true})
		return
	}

	// Wrong password or username bro
	fmt.Fprint(res, Response{"success": false})
}

func StudentRegister(res http.ResponseWriter, req *http.Request) {

}
