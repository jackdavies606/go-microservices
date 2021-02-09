package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/customer/{name}", GetCustomer).Methods("GET")
	router.HandleFunc("/customer/{name}", AddCustomer).Methods("POST")
	router.HandleFunc("/customer/{name}", RemoveCustomer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	InitialMigration()

	handleRequests()
}
