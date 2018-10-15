package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// read in database credentials
var db_credentials, err = ioutil.ReadFile(".db_credentials")

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

	// inform http_out that content is JSON
	http_out.Header().Add("Content-Type", "application/json")

	// get request varaibles
	vars := mux.Vars(r)
	article := vars["article"]

	// establish connection to databse
	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	// execute query
	rows, err := db.Query(
		`select array_to_json(array_agg(row_to_json(t))) as res_string from
			(select 
			a.page_id, a.page_name, b.year, b.page_views 
			from 
			(select * from dict_en where page_name = $1) as a 
                            left join imports_en as b on a.page_id = b.page_id 
			) as t 
		;`,
		article)

	checkErr(err)
	for rows.Next() {
		var res_string string
		err = rows.Scan(&res_string)
		checkErr(err)
		fmt.Fprintf(http_out, "%s", res_string)
	}

}

func api_search(http_out http.ResponseWriter, r *http.Request) {
	http_out.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	search := vars["search"]

	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(
		`select array_to_json(array_agg(row_to_json(t))) as res_string from
			(select a.page_id, a.page_name, b.year, b.page_views 
			from 
			    (select * from dict_en where page_name ~ $1) as a 
                            left join imports_en as b on a.page_id = b.page_id 
			limit 100) as t
			;`, search)
	checkErr(err)

	for rows.Next() {
		var res_string string
		err = rows.Scan(&res_string)
		checkErr(err)
		fmt.Fprintf(http_out, "%s", res_string)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
