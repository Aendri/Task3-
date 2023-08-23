package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Countries struct {
	Name     string `json:"name"`
	Capital  string `json:"capital"`
	Currency string `json:"currency"`
	Region   string `json:"region"`
}

func CreateTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS countries (name TEXT NOT NULL, capital TEXT NOT NULL, currency TEXT NOT NULL, region TEXT NOT NULL)`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("table exists")
		log.Fatal("error creating table", err)
	}
}

func Loadcountriesname(db *sql.DB, filepath string) error {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var countries []Countries

	err = json.Unmarshal(data, &countries)
	if err != nil {
		return err
	}

	for _, country := range countries {
		_, err := db.Exec("INSERT INTO countries(name, capital,currency,region)VALUES(? ,?, ?, ?)", country.Name, country.Capital, country.Currency, country.Region)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Data entered succesfully")
	return nil

}

func Getcountries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var countries []Countries
		region := r.URL.Query().Get("region")
		query := "SELECT name,capital,currency,region FROM countries"
		if region != "" {
			query += " WHERE region =?"

		}
		rows, err := db.Query(query, region)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var country Countries
			err := rows.Scan(&country.Name, &country.Capital, &country.Currency, &country.Region)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			countries = append(countries, country)

		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(countries)

	}

}

func main() {

	db, err := sql.Open("sqlite3", "countries.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	CreateTable(db)
	router := mux.NewRouter()

	filepath := "./countries.json"

	err = Loadcountriesname(db, filepath)
	if err != nil {
		fmt.Println("Error in loading the countries data")
		return
	}
	router.HandleFunc("/countries", Getcountries(db)).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Server is running on port 8080")

	http.ListenAndServe(":8080", nil)

}
