package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
)

var db *gorm.DB
var err error

type Customer struct {
	gorm.Model
	Name string `json:"name"`
}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Customer{})
}

func GetCustomerByName(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Cusomer database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var customer Customer
	db.Where("name = ?", name).Find(&customer)
	json.NewEncoder(w).Encode(customer)
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Cusomer database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]


	var customer Customer
	db.Where("ID = ?", id).Find(&customer)
	json.NewEncoder(w).Encode(customer)
}

func AddCustomer(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Cusomer database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	db.Create(&Customer{Name: name})

	fmt.Fprint(w, "New customer added")
}

func RemoveCustomer(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Cusomer database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var customer Customer
	db.Where("name = ?", name).Find(&customer)
	db.Delete(&customer)

	fmt.Fprint(w, "Customer deleted")
}