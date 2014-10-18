package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Review struct to hold data in
type Review struct {
	Id            uint           `json:"id"`
	DisplayName   string         `db:"display_name" json:"display_name"`
	StudentEmail  string         `db:"student_email" json:"student_email"`
	ClassName     string         `db:"class_name" json:"class_name"`
	Review        string         `json:"review"`
	InstructorId  uint           `db:"instructor_id" json:"instructor_id"`
	Grading       uint           `json:"grading"`
	Helpfulness   uint           `json:"helpfulness"`
	Availability  uint           `json:"availability"`
	Homework      uint           `json:"homework"`
	Novateachers  sql.NullString `json:"novateachers"`
	DateSubmitted time.Time      `db:"date_submitted" json:"date_submitted"`
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

	// We go them all so return the data
	js, errjs := json.Marshal(reviews)
	if errjs != nil {
		log.Println("Error marshalling JSON data.")
		log.Println(errjs)
		fmt.Fprint(res, Response{"success": false, "message": "Could not format reviews!"})
		return
	}

	// Write it!
	fmt.Fprint(res, string(js))
}
