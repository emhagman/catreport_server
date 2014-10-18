package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func StudentLogin(res http.ResponseWriter, req *http.Request) {

	// Get username and password
	username := req.FormValue("username")
	password := req.FormValue("password")

	// Make sure we actually got a username and password
	// In Go, strings are never null or 'nil'
	if username == "" || password == "" {
		log.Println("empty username or password empty")
		fmt.Fprint(res, Response{"success": false})
		return
	}

	// Automatically add @villanova.edu so people can't mess with us
	// Villanova students only, remember?
	if !strings.Contains(username, "@villanova.edu") {
		log.Println("appended villanova e-mail address")
		username += "@villanova.edu"
	}

	// Back before I knew anything about programming...
	// I used MD5 to hash the passwords without any salt
	// Granted, there isn't really anything intense stored in the database
	// But this is terrible practice and I am ashamed of old me
	oldPassword := OldHashPassowrd(password)

	// Make sure we have connection first
	// Defer close to function end
	db, err := DBConnect()
	defer db.Close()

	if err != nil {
		log.Println("failed to connect to database")
		log.Println(err.Error())
		fmt.Fprint(res, Response{"success": false})
		return
	}

	// Go validate our user (using the old password first)
	rows, err := db.Query("SELECT * FROM students WHERE email = $1 AND password = $2", username, oldPassword)
	if err != nil {
		log.Println("error looking up using old hash")
		log.Println(err)
		fmt.Fprint(res, Response{"success": false})
		return
	}

	// New password is hashed using pbkdf2 and the salt is kept on server :P
	newPassword := HashPassword([]byte(password), []byte(os.Getenv("CATREPORT_SALT")))

	// See if we have at least one row using old hash
	if rows.Next() {

		log.Println("login with old hash found")

		// update the with new hash
		// if this fails, we can keep using the old md5 but this isn't really safe...
		stmt, err := db.Prepare("UPDATE students SET password=$1 WHERE email=$2")
		if err != nil {
			log.Println("error creating stmt for updating new password hash")
			log.Println(err)
		} else {
			// update it
			// again, if this fails, we can still use md5 for now
			_, err2 := stmt.Exec(newPassword, username)
			if err2 != nil {
				log.Println("error updating new password hash")
				log.Println(err2)
			}
		}

		// login and return success
		SessionLogin(res, req)
		fmt.Fprint(res, Response{"success": true, "old": true})
		return
	}

	// Let's check to see if they have the new hashing algo
	// I really should just convert everyone eventually...
	rows, err = db.Query("SELECT * FROM students WHERE email = $1 AND password = $2", username, newPassword)
	if err != nil {
		log.Println("error looking up using new hash")
		log.Println(err)
		fmt.Fprint(res, Response{"success": false})
		return
	}

	// See if we have at least one row
	for rows.Next() {
		SessionLogin(res, req)
		fmt.Fprint(res, Response{"success": true, "new": true})
		return
	}

	// Wrong password or username bro
	log.Println("invalid username or password")
	fmt.Fprint(res, Response{"success": false})
}

func StudentRegister(res http.ResponseWriter, req *http.Request) {

	// Get register information
	email := req.FormValue("email")
	password := req.FormValue("password")
	first := req.FormValue("first")
	last := req.FormValue("last")

	// Validate that shit
	if email == "" || password == "" || first == "" || last == "" {
		fmt.Fprint(res, Response{"success": false, "message": "Missing something!"})
		return
	}

	// Make sure we have a villanova e-mail
	// Seriously, if you look at this, it's obvious you can just make up an e-mail
	// Most people won't read this code. E-mail verification would fix this problem...
	// Not really guarding NSA secrets so this doesn't matter
	if !strings.Contains(email, "@villanova.edu") {
		fmt.Fprint(res, Response{"success": false, "message": "Must be a valid Villanova e-mail."})
		return
	}

	// Make sure we have connection first
	// Defer close to function end
	db, err := DBConnect()
	defer db.Close()
	if err != nil {
		log.Println("failed to connect to database")
		log.Println(err.Error())
		fmt.Fprint(res, Response{"success": false, "message": "The database appears to be down..."})
		return
	}

	// Prepare our statement
	stmt, err := db.Prepare("INSERT INTO students (sfirst, last, email, password, verified) VALUES($1, $2, $3, $4, false)")
	if err != nil {
		log.Println("error creating stmt for registering user")
		log.Println(err.Error())
		fmt.Fprint(res, Response{"success": false, "message": "Something went horribly wrong!"})
		return
	}

	// Execute insert for registration
	_, err2 := stmt.Exec(first, last, email, password)
	if err2 != nil {
		log.Println("error registering new user")
		log.Println(err2)
		fmt.Fprint(res, Response{"success": false, "message": "Error registering new user!"})
	}

	// We made it this far which means they are a user. Give em a session
	SessionLogin(res, req)
	fmt.Fprint(res, Response{"success": true})
}
