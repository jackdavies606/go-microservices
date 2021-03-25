package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB
var err error

// DB model
type OrderEntry struct {
	ID int `json:"id" gorm:"primaryKey"`
	ItemId int `json:"itemId"`
	OrderId int `json:"orderId"`
}

// DB model
type Order struct {
	ID int `json:"id" gorm:"primaryKey"`
	CustomerId int  `json:"customerId"`
	IsComplete bool `json:"isComplete"`
}

// Response model
type OrderResponse struct {
	OrderId int
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

	orderLines, err := ReadCsv("orders.csv")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	populateOrdersDatabase(orderLines)

	orderEntryLines, err := ReadCsv("order_entries.csv")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	populateOrderEntriesDatabase(orderEntryLines)
}

func populateOrdersDatabase(lines [][]string) {
	for _, line := range lines {
		var id int
		if id, err = strconv.Atoi(line[0]); err != nil {
			panic(err)
		}

		var customerId int
		if customerId, err = strconv.Atoi(line[1]); err != nil {
			panic(err)
		}

		var isComplete bool
		if isComplete, err = strconv.ParseBool(line[2]); err != nil {
			panic(err)
		}

		order := Order{
			ID: id,
			CustomerId: customerId,
			IsComplete: isComplete,
		}

		fmt.Printf("Read: CustomerId %s, IsComplete %s", strconv.Itoa(order.CustomerId),
			strconv.FormatBool(order.IsComplete))

		db.Create(&order)
	}
}

func populateOrderEntriesDatabase(lines [][]string) {
	for _, line := range lines {
		var id int
		if id, err = strconv.Atoi(line[0]); err != nil {
			panic(err)
		}

		var itemId int
		if itemId, err = strconv.Atoi(line[1]); err != nil {
			panic(err)
		}

		var orderId int
		if orderId, err = strconv.Atoi(line[2]); err != nil {
			panic(err)
		}

		orderEntry := OrderEntry {
			ID: id,
			ItemId: itemId,
			OrderId: orderId,
		}
		db.Create(&orderEntry)
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

// get open order by customer
func GetCustomersOpenOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	parsedCustomerId, customerIdParseErr := strconv.Atoi(vars["customerId"])

	if customerIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id provided")
		return
	}

	// get Order
	var order Order
	order, err := findCustomerOrder(parsedCustomerId, false)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Customer does not have an open Order")
		return
	}

	// get OrderEntry
	var entries []OrderEntry
	db.Where("order_id = ? ", order.ID).Find(&entries)

	var items []Item
	for _, entry := range entries {
		itemIdString := strconv.Itoa(entry.ItemId)
		item, err := getItem(w, itemIdString)

		if err != nil {
			fmt.Fprintf(w, "Failed to retreive item for order.")
			return
		}

		fmt.Printf("Adding Item '%s' with price '%s' to Items", item.Name, strconv.Itoa(item.Price))
		items = append(items, item)
	}

	var response = OrderResponse {
		OrderId:    order.ID,
		IsComplete: order.IsComplete,
		CustomerId: order.CustomerId,
		Items: items,
	}

	json.NewEncoder(w).Encode(response)
}

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

// create an OrderEntry - creates Order if an open Order does not exist
func AddToOrder(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Order database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	parsedCustomerId, customerIdParseErr := strconv.Atoi(vars["customerId"])
	parsedItemId, itemIdParseErr := strconv.Atoi(vars["itemId"] )
	if customerIdParseErr != nil || itemIdParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid customer id or item id provided")
		return
	}

	customerId := parsedCustomerId
	itemId := parsedItemId

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

	itemIdString := strconv.Itoa(itemId)
	_, err = getItem(w, itemIdString)
	if err != nil {
		fmt.Fprintf(w, "Could not retrieve item to add to order with itemId '%s'", itemIdString)
		return
	}

	// create order entry
	var orderEntry = OrderEntry{
		ItemId: itemId,
		OrderId: order.ID,
	}
	db.Table("order_entries").Create(&orderEntry)

	fmt.Fprint(w, "New item added")
}

func findCustomerOrder(customerId int, isOpen bool) (Order, error) {
	var order Order
	err := db.Where("customer_id = ? AND is_complete = ?", customerId, isOpen).First(&order).Error
	return order, err
}

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
	err := db.Where("customer_id = ? AND is_complete = ?", customerId, false).First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Customer does not have an incomplete Order to delete")
		return
	}

	err = db.Where("customer_id = ? AND is_complete = ?", customerId, false).Delete(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to delete Order")
		return
	}

	// deleting all orderEntries related to order
	db.Where("order_id = ?", order.ID).Delete(OrderEntry{})

	fmt.Fprint(w, "Order deleted")
}

func getItem(w http.ResponseWriter, itemId string) (Item, error) {
	itemUrl := os.Getenv("ITEM_SERVICE_URL")
	requestUrl := itemUrl + "/item/id/" + itemId

	fmt.Printf("Making a request to: %s", requestUrl)
	resp, err := http.Get(requestUrl)
	var item Item

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Request to item service failed for itemId '%s': ", itemId)
		return item, err
	} else if resp.Status != "200 OK" {
		w.WriteHeader(http.StatusNotFound)
		err = errors.New("item could not be found")
		fmt.Printf("The response status on the request for itemId '%s' was not 200: %s. %s", itemId, resp.Status, err)
		return item, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not read the response body for itemId '%s'. %s", itemId, err)
		return item, err
	}

	err = json.Unmarshal(body, &item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not create Item from response body for itemId %s. %s", itemId, err)
		return item, err
	}

	fmt.Printf("Successfully retrieved and created item with itemId %s", itemId)

	return item, nil
}

func createOrderResponse(order Order, items []Item) OrderResponse {
	return OrderResponse{
		OrderId: order.ID,
		IsComplete: order.IsComplete,
		CustomerId: order.CustomerId,
		Items: items,
	}
}