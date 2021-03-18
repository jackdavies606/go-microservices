package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
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

	// TODO : use this to run the sql file
	query, err := ioutil.ReadFile("path/to/database.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}

// This method should only be callable by an admin user
func GetCustomers(w http.ResponseWriter, r *http.Request)  {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Cusomer database")
	}
	defer db.Close()

	var customers []Customer
	db.Find(&customers)

	json.NewEncoder(w).Encode(customers)
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
	if findErr := db.Where("name = ?", name).First(&customer).Error; errors.Is(findErr, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The requested customer with name '%v' was not found", name)
		return
	}

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
	if findErr := db.Where("ID = ?", id).First(&customer).Error; errors.Is(findErr, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The requested customer with id '%v' was not found", id)
		return
	}

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