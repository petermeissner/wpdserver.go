package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// read in database credentials
var db_credentials, err = ioutil.ReadFile(".db_credentials")

// main function
func main() {

	// initialize router
	router := mux.NewRouter().StrictSlash(true)

	// Indes/
	router.HandleFunc("/", api_index)

	// Article/
	router.HandleFunc("/article/exact/{lang}/{article}", api_article_exact)
	router.HandleFunc("/article/search/{lang}/{article}", api_article_search)

	// Search/
	router.HandleFunc("/search/{lang}/{search}", api_search)

	// initialize server with router and routes
	log.Fatal(http.ListenAndServe(":8880", router))
}

// index route
func api_index(http_out http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(
		http_out,
		`{
"note": "Data has been gathered by Peter Meissner in a project comissioned by Hertie School of Governance (Simon Munzert) with support by Daimler and Benz Foundation.",
"api_1":
  [
    {
      "description": "article stats",
      "path": "article/exact/{lang}/{articlename}"
    },
    {
      "description": "article stats",
      "path": "article/search/{lang}/{articlename}"
    },
    {
      "description": "search article stats",
      "path": "search/{lang}/{regex}"
    }
  ]
}
`)
}

// article exact route
func api_article_exact(http_out http.ResponseWriter, r *http.Request) {

	// inform http_out that content is JSON
	http_out.Header().Add("Content-Type", "application/json")

	// get request varaibles
	vars := mux.Vars(r)
	article := strings.ToLower(vars["article"])

	var rx = regexp.MustCompile(`\W`)
	lang := vars["lang"]
	lang = strings.ToLower(lang)
	lang = rx.ReplaceAllString(lang, "")
	lang = lang[:2]

	// establish connection to databse
	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	// execute query
	rows, err := db.Query(
		`select array_to_json(array_agg(row_to_json(t))) as res_string from
		(select 
		a.page_id, 
		a.page_name, 
		b.year, 
		b.page_views as page_view_count
		-- generate_series((b.year || '-01-01')::date, (b.year || '-12-31')::date, '1 day'::interval)::date as page_view_date, 
		-- regexp_split_to_table(b.page_views, ',') as page_view_count 
		from 
			(select * from dict_`+lang+` where page_name = $1 limit 100) as a
			left join page_views_daily_`+lang+` as b on a.page_id = b.page_id 
	) as t;`, article)

	checkErr(err)
	for rows.Next() {

		var res_string string
		err = rows.Scan(&res_string)

		if err != nil {
			fmt.Fprintf(http_out, "%s", `{"status": "error"}`)
		} else {
			fmt.Fprintf(http_out, `{"status": "ok", "data": %s}`, res_string)
		}
	}

}

// articel search route
func api_article_search(http_out http.ResponseWriter, r *http.Request) {
	http_out.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	search := strings.ToLower(vars["article"])

	var rx = regexp.MustCompile(`\W`)
	lang := vars["lang"]
	lang = strings.ToLower(lang)
	lang = rx.ReplaceAllString(lang, "")
	lang = lang[:2]

	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(
		`select array_to_json(array_agg(row_to_json(t))) as res_string from
			(select 
				a.page_id, 
				a.page_name, 
				b.year,
				--generate_series((b.year || '-01-01')::date, (b.year || '-12-31')::date, '1 day'::interval)::date as page_view_date, 
				--regexp_split_to_table(b.page_views, ',') as page_view_count 
				b.page_views as page_view_count
				from 
					(select * from dict_`+lang+` where page_name ~ $1 limit 100) as a
					left join page_views_daily_`+lang+` as b on a.page_id = b.page_id 
		) as t;`, search)
	checkErr(err)

	for rows.Next() {
		var res_string string
		err = rows.Scan(&res_string)

		if err != nil {
			fmt.Fprintf(http_out, "%s", `{"status": "error"}`)
		} else {
			fmt.Fprintf(http_out, `{"status": "ok", "data": %s}`, res_string)
		}
	}

}

// search route
func api_search(http_out http.ResponseWriter, r *http.Request) {
	http_out.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	search := strings.ToLower(vars["search"])

	var rx = regexp.MustCompile(`\W`)
	lang := vars["lang"]
	lang = strings.ToLower(lang)
	lang = rx.ReplaceAllString(lang, "")
	lang = lang[:2]

	db, err := sql.Open("postgres", string(db_credentials))
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(
		`select array_to_json(array_agg(row_to_json(t))) as res_string from
			(select 
				a.page_id, 
				a.page_name, 
				'`+lang+`' as page_lang
				from 
			    (select * from dict_`+lang+` where page_name ~ $1) as a 
			) as t
			;`, search)
	checkErr(err)

	for rows.Next() {
		var res_string string
		err = rows.Scan(&res_string)

		if err != nil {
			fmt.Fprintf(http_out, "%s", `{"status": "error"}`)
		} else {
			fmt.Fprintf(http_out, `{"status": "ok", "data": %s}`, res_string)
		}
	}

}

// helper function: error checker
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
