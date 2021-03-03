package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items", GetItems).Methods("GET")
	router.HandleFunc("/item/name/{name}", GetItemByName).Methods("GET")
	router.HandleFunc("/item/id/{id}", GetItemById).Methods("GET")
	router.HandleFunc("/item/name/{name}", RemoveItem).Methods("DELETE")
	router.HandleFunc("/item", AddItem).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	InitialMigration()

	handleRequests()
}
