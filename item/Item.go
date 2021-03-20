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
	"strconv"
)

var db *gorm.DB
var err error

type Item struct {
	gorm.Model
	Name string `json:"name"`
	Price int `json:"price"`
}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to Item database")
	}
	defer db.Close()

	db.AutoMigrate(&Item{})

	lines, err := ReadCsv("./items.csv")
	if err != nil {
		panic(err)
	}

	populateDatabase(lines)
}


func populateDatabase(lines [][]string) {
	for _, line := range lines {
		var price int64
		if price, err = strconv.ParseInt(line[1], 10, 64); err != nil {
			panic(err)
		}

		item := Item{
			Name: line[0],
			Price: int(price),
		}
		fmt.Printf("Read: Name %s, Price %s", item.Name, strconv.Itoa(item.Price))
		db.Create(&item)
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

func GetItemByName(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var item Item
	if findErr := db.Where("name = ?", name).First(&item).Error; errors.Is(findErr, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The requested item '%v' was not found", name)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func GetItemById(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var item Item

	if findErr := db.Where("ID = ?", id).First(&item).Error; errors.Is(findErr, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The requested item with id '%v' was not found", id)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to the Item database")
	}
	defer db.Close()

	var items []Item
	db.Find(&items)

	json.NewEncoder(w).Encode(items)
}

func AddItem(w http.ResponseWriter, r *http.Request) {
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

func RemoveItem(w http.ResponseWriter, r *http.Request) {
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