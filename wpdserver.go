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
	fmt.Fprintln(http_out, "{article}")

	vars := mux.Vars(r)
	article := vars["article"]

	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(
		`select 
			a.page_id, a.page_name, b.year, b.page_views 
			from 
			    (select * from dict_en where page_name = $1) as a 
                            left join imports_en as b on a.page_id = b.page_id 
			;`, article)
	checkErr(err)

	for rows.Next() {
		var page_id string
		var page_name string
		var year string
		var page_views string
		err = rows.Scan(&page_id, &page_name, &year, &page_views)
		checkErr(err)
		fmt.Fprintf(http_out, "{")
		fmt.Fprintf(http_out, "\n  \"search\": \"%s\",", article)
		fmt.Fprintf(http_out, "\n  \"page_id\": %s,", page_id)
		fmt.Fprintf(http_out, "\n  \"page_name\": \"%s\",", page_name)
		fmt.Fprintf(http_out, "\n  \"year\": %s,", year)
		fmt.Fprintf(http_out, "\n  \"page_views\": [%s]", page_views)
		fmt.Fprintf(http_out, "\n}")
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
