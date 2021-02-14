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

type Order struct {
	gorm.Model
	Name string `json:"name"`
	ItemId int `json:"itemId"`
	CustomerId int  `json:"customerId"`

}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to Item database")
	}
	defer db.Close()

	db.AutoMigrate(&Item{})
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var item Item
	db.Where("name = ?", name).Find(&item)
	json.NewEncoder(w).Encode(item)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	var items []Item
	db.Find(&items)

	json.NewEncoder(w).Encode(items)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Create(&item)

	fmt.Fprint(w, "New item added")
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var item Item
	db.Where("name = ?", name).Find(&item)
	db.Delete(&item)

	fmt.Fprint(w, "Item deleted")
}