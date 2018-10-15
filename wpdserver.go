package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "****"
	DB_PASSWORD = "****"
	DB_NAME     = "****"
	DB_HOST     = "****"
)

func main() {

	// initialize router
	router := mux.NewRouter().StrictSlash(true)

	// Indes/
	router.HandleFunc("/", api_index)

	// Article/
	router.HandleFunc("/article/{article}", api_article)

	// Search/
	router.HandleFunc("/search/{search}", api_search)

	// initialize server with router and routes
	log.Fatal(http.ListenAndServe(":8880", router))
}

func api_index(http_out http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(http_out, "- /article/{artiklename}\n- /search/{regex}\n")
}

func api_article(http_out http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(http_out, "{article}")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT distinct page_view_date  FROM import_jobs;")
	checkErr(err)
	fmt.Fprintln(http_out, "page_view_date")

	for rows.Next() {
		var page_view_date string
		err = rows.Scan(&page_view_date)
		checkErr(err)
		fmt.Fprintf(http_out, "%s\n", page_view_date)
	}
}

func api_search(http_out http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(http_out, "{search}")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
