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

// DB model
type OrderEntry struct {
	gorm.Model
	ItemId int `json:"itemId"`
	CustomerId int  `json:"customerId"`
	OrderId int `json:"orderId"`
}

// DB model
type Order struct {
	gorm.Model
	CustomerId int  `json:"customerId"`
	IsComplete bool `json:"isComplete"`
}

// Response model
type OrderResponse struct {
	OrderId uint
	IsComplete bool
	CustomerId int
	Items []Item
}

// Response Model
type Item struct {
	Name string
	Price int
}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to Order database")
	}
	defer db.Close()

	db.AutoMigrate(&Order{})
	db.AutoMigrate(&OrderEntry{})
}

// get open order by customer
func GetCustomersOpenOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var order OrderEntry
	db.Where("name = ?", name).Find(&order)
	json.NewEncoder(w).Encode(order)
}

// gets open and closed orders for a customer
func GetAllCustomerOrders(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	var orders []OrderEntry
	db.Find(&orders)

	json.NewEncoder(w).Encode(orders)
}

// create an OrderEntry - creates Order if an open Order does not exist
func AddToOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	var item Order
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Create(&item)

	fmt.Fprint(w, "New item added")
}

// deletes Order and related OrderEntry
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var item Order
	db.Where("name = ?", name).Find(&item)
	db.Delete(&item)

	fmt.Fprint(w, "Item deleted")
}

func createOrderResponse(order Order, items []Item) OrderResponse {
	return OrderResponse{
		OrderId: order.ID,
		IsComplete: order.IsComplete,
		CustomerId: order.CustomerId,
		Items: items,
	}
}