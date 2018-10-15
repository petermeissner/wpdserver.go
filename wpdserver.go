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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/todos", TodoIndex)
	router.HandleFunc("/todos/{todoId}", TodoShow)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT distinct page_view_date  FROM import_jobs;")
	checkErr(err)
	fmt.Println("page_view_date")

	for rows.Next() {
		var page_view_date string
		err = rows.Scan(&page_view_date)
		checkErr(err)
		fmt.Printf("%s\n", page_view_date)
	}

	log.Fatal(http.ListenAndServe(":8880", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo Index!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
