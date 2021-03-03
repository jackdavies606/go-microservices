package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/orders/{customerId}", GetAllCustomerOrders).Methods("GET")
	router.HandleFunc("/order/{customerId}", GetCustomersOpenOrder).Methods("GET")
	router.HandleFunc("/order/{customerId}/item/{itemId}", AddToOrder).Methods("POST")
	router.HandleFunc("/order/{customerId}", CancelOrder).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	InitialMigration()

	handleRequests()
}
