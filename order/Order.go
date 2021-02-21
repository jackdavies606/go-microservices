package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"strconv"
	"time"
)

var db *gorm.DB
var err error

// DB model
type OrderEntry struct {
	gorm.Model
	ItemId uint `json:"itemId"`
	CustomerId uint  `json:"customerId"`
	OrderId uint `json:"orderId"`
}

// DB model
type Order struct {
	gorm.Model
	CustomerId uint  `json:"customerId"`
	IsComplete bool `json:"isComplete"`
}

// Response model
type OrderResponse struct {
	OrderId uint
	IsComplete bool
	CustomerId uint
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

// todo : this method WORKS :)
// get open order by customer
func GetCustomersOpenOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	parsedCustomerId, customerIdParseErr := strconv.ParseUint(vars["customerId"], 10, 64)

	if customerIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id provided")
		return
	}

	customerId := uint(parsedCustomerId)

	// get Order
	var order Order
	order, err := findCustomerOrder(customerId, false)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Customer does not have an open Order")
		return
	}

	// get OrderEntry
	var entries []OrderEntry
	db.Where("order_id = ?", order.ID).Find(&entries)

	// todo - for each OrderEntry get the item and create an OrderResponse

	json.NewEncoder(w).Encode(order)
}

// todo : MAKE THIS WORK AS IT SHOULD
// gets open and closed orders for a customer
func GetAllCustomerOrders(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	parsedCustomerId, customerIdParseErr := strconv.ParseUint(vars["customerId"], 10, 64)

	if customerIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id provided")
		return
	}

	customerId := uint(parsedCustomerId)

	var orders []OrderEntry
	err = db.Where("customer_id = ?", customerId).First(&orders).Error

	json.NewEncoder(w).Encode(orders)
}

// todo : this method WORKS :)
// create an OrderEntry - creates Order if an open Order does not exist
func AddToOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	parsedCustomerId, customerIdParseErr := strconv.ParseUint(vars["customerId"], 10, 64)
	parsedItemId, itemIdParseErr := strconv.ParseUint(vars["itemId"], 10, 64)
	if customerIdParseErr != nil || itemIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id or item id provided")
		return
	}

	customerId := uint(parsedCustomerId)
	itemId := uint(parsedItemId)

	order, err := findCustomerOrder(customerId, false)

	// create an Order if one does not exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var newOrder = Order{
			CustomerId: customerId,
			IsComplete: false,
		}

		db.Create(&newOrder)
		time.Sleep(1000)
		order, err = findCustomerOrder(newOrder.CustomerId, false)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed create customer Order")
			return
		}
	}

	// todo : add a call to the item service to validate ItemId is valid
	// todo : add a call to the customer service to validate customerId is valid

	// create order entry
	var orderEntry = OrderEntry{
		CustomerId: customerId,
		ItemId: itemId,
		OrderId: order.ID,
	}
	db.Table("order_entries").Create(&orderEntry)

	fmt.Fprint(w, "New item added")
}

func findCustomerOrder(customerId uint, isOpen bool) (Order, error) {
	var order Order
	err := db.Where("customer_id = ? AND is_complete = ?", customerId, isOpen).First(&order).Error
	return order, err
}

// todo : this method WORKS :)
// deletes Order and related OrderEntry
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	customerId := vars["customerId"]

	var order Order
	err := db.Where("customer_id = ?", customerId).First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Customer does not have Order to delete")
		return
	}

	err = db.Where("customer_id = ? AND is_complete = ?", customerId, false).Delete(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to delete Order")
		return
	}

	// todo : cleanup related OrderEntries

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