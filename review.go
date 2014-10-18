package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Review struct to hold data in
type Review struct {
	Id            uint
	DisplayName   string `db:"display_name"`
	StudentEmail  string `db:"student_email"`
	ClassName     string `db:"class_name"`
	Review        string
	InstructorId  uint `db:"instructor_id"`
	Grading       uint
	Helpfulness   uint
	Availability  uint
	Homework      uint
	Novateachers  sql.NullString
	DateSubmitted time.Time `db:"date_submitted"`
}

func ReviewGetReviewsById(res http.ResponseWriter, req *http.Request) {

	// Get the instructor id
	vars := mux.Vars(req)
	id := vars["id"]

	// Convert to int (if we can)
	instructorId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("invalid instructor id given")
		log.Println(err.Error())
		fmt.Fprint(res, Response{"success": false, "message": "Not a valid instructor id!"})
		return
	}

	// Make sure we have connection first
	// Defer close to functon end
	db, err := DBConnect()
	defer db.Close()
	if err != nil {
		log.Println("failed to connect to database")
		log.Println(err.Error())
		fmt.Fprint(res, Response{"success": false, "message": "Could not connect to the database!"})
		return
	}

	// Go validate our user (using the old password first)
	reviews := []Review{}
	err = db.Select(&reviews, "SELECT * FROM reviews WHERE instructor_id=$1", instructorId)
	if err != nil {
		log.Println("error looking up reviews for instructor")
		log.Println(err)
		fmt.Fprint(res, Response{"success": false, "message": "Something went wrong when getting the reviews!"})
		return
	}

	// Go through the rows and read in the data
	for i := range reviews {
		log.Println(reviews[i].ClassName)
	}
}
