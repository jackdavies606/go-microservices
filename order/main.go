package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items", GetItems).Methods("GET")
	router.HandleFunc("/item/{name}", GetItem).Methods("GET")
	router.HandleFunc("/item", AddItem).Methods("POST")
	router.HandleFunc("/item/{name}", RemoveItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	InitialMigration()

	handleRequests()
}
