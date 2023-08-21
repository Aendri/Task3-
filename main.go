package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Countries struct {
	Name     string `json:"name"`
	Capital  string `json:"capital"`
	Currency string `json:"currency"`
	Region   string `json:"region"`
}

var countries []Countries

func Loadcountriesname(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &countries)
	if err != nil {
		return err
	}
	return nil

}

func Getcountries(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	filteredCountries := countries
	if region != "" {
		filteredCountries = filteredCountriesByRegion(region)

	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(filteredCountries)

}

func filteredCountriesByRegion(region string) []Countries {
	var fcountry []Countries
	for _, country := range countries {
		if strings.EqualFold(country.Region, region) {
			fcountry = append(fcountry, country)
		}
	}
	return fcountry
}

func main() {
	router := mux.NewRouter()
	filepath := "countries.json"

	err := Loadcountriesname(filepath)
	if err != nil {
		fmt.Println("Error in loading the countries data")
		return
	}
	router.HandleFunc("/countries", Getcountries).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Server is running on port 8080")

	http.ListenAndServe(":8080", nil)

}
