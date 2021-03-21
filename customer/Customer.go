package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"os"
)

var db *gorm.DB
var err error

type Customer struct {
	ID int `json:"id" gorm:"primaryKey"` // todo: add to csv
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

	lines, err := ReadCsv("./customers.csv")
	if err != nil {
		fmt.Println("Failed to read ./customers.csv")
		panic(err)
	}

	populateDatabase(lines)
}

func populateDatabase(lines [][]string) {
	for _, line := range lines {
		customer := Customer{
			Name: line[0],
		}
		fmt.Printf("Read: %s", customer.Name)
		db.Create(&customer)
	}
}

func ReadCsv(csvPath string) ([][]string, error){
	f, err := os.Open(csvPath)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
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