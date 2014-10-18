package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ReviewStruct
type ReviewStruct struct {
	display_name   string
	grading        uint
	class_name     string
	review         string
	date_submitted time.Time
	instructor_id  uint
	id             uint
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
	rows, err := db.Query("SELECT * FROM reviews WHERE instructor_id=$1", instructorId)
	if err != nil {
		log.Println("error looking up reviews for instructor")
		log.Println(err)
		fmt.Fprint(res, Response{"success": false, "message": "Something went wrong when getting the reviews!"})
		return
	}

	// Go through the rows and read in the data
	for rows.Next() {

		// load data into a struct
		data := ReviewStruct{}
		rows.Scan(&data)

		log.Println(data.class_name)
	}
}
