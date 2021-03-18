package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/orders/customer/{customerId}", GetAllCustomerOrders).Methods("GET")
	router.HandleFunc("/order/customer/{customerId}", GetCustomersOpenOrder).Methods("GET")
	router.HandleFunc("/order/customer/{customerId}/item/{itemId}", AddToOrder).Methods("POST")
	router.HandleFunc("/order/customer/{customerId}", CancelOrder).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	InitialMigration()

	handleRequests()
}
