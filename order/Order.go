package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"strconv"
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

// get open order by customer
func GetCustomersOpenOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	customerId := vars["customerId"]

	// get Order
	var order Order
	db.Where("customer_id = ? AND is_complete = ?", customerId, false).Find(&order)

	// get OrderEntry
	var entries []OrderEntry
	db.Where("order_id = ?", order.ID).Find(&entries)

	// todo - for each OrderEntry get the item and create an OrderResponse

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
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	customerId := vars["customerId"]
	itemId := vars["itemId"]

	parsedCustomerId, customerIdParseErr := strconv.ParseUint(customerId, 10, 64)
	parsedItemId, itemIdParseErr := strconv.ParseUint(itemId, 10, 64)
	if customerIdParseErr != nil || itemIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id or item id provided")
		return
	}

	order := findCustomerOrder(customerId, false)

	// create an Order if one does not exist
	if &order == nil {
		var newOrder = Order{
			CustomerId: uint(parsedCustomerId),
			IsComplete: false,
		}

		db.Create(&newOrder)
		order = findCustomerOrder(customerId, false)
	}

	// todo : add a call to the item service to validate ItemId is valid
	// todo : add a call to the customer service to validate customerId is valid

	// create order entry
	var orderEntry = OrderEntry{
		CustomerId: uint(parsedCustomerId),
		ItemId: uint(parsedItemId),
		OrderId: order.ID,
	}
	db.Table("order_entries").Create(orderEntry)

	fmt.Fprint(w, "New item added")
}

func findCustomerOrder(customerId string, isOpen bool) Order {
	var order Order
	db.Where("customer_id = ? AND is_complete", customerId, isOpen).Find(&order)
	return order
}

// deletes Order and related OrderEntry
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["customerId"]

	var order Order
	db.Where("name = ?", name).Find(&order)
	db.Delete(&order)

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