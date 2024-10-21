package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func handleIndex(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func handleCategories(w http.ResponseWriter, r * http.Request, db * gorm.DB) {
	var categories []Category
	result := db.Find(&categories)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
	// for _, category := range categories {
	// 	fmt.Println(category.CategoryName, category.Slug)
	// }
}

type Category struct {
	CategoryName string
	Slug string
}

func main() {
	var port = "8000"

	db, err := gorm.Open(mysql.Open("root@tcp(localhost)/db_ujikom"), &gorm.Config{}) 
	if err != nil {
		panic("Failed Connect to Database!")
	}
	fmt.Println("Success Connect to Database!")

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/categories", func(w http.ResponseWriter, r * http.Request) {
		handleCategories(w, r, db)
	})

	fmt.Printf("Server Started at Port %s\n", port)
	http.ListenAndServe(":" + port, nil)
}