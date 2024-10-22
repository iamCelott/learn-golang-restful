package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func handleIndex(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

type Category struct {
	CategoryName string `json:"category_name"`
	Slug string `json:"slug"`
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
}

func handleCreateCategory(w http.ResponseWriter, r * http.Request, db * gorm.DB) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var category Category	
	if err := json.Unmarshal(body, &category); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	
	if result := db.Create(&category); result.Error != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
        return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Category created successfully"})
}

func handleEditCategory(w http.ResponseWriter, r * http.Request, db * gorm.DB) {
    id := r.PathValue("id")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var category Category	
	if err := json.Unmarshal(body, &category); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	if result := db.Where("id = ?", id).Updates(category); result.Error != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
    	return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id":id, "message": "Category updated successfully"})
}

func handleDeleteCategory(w http.ResponseWriter, r * http.Request, db * gorm.DB) {
	id := r.PathValue("id")

	if result := db.Where("id = ?", id).Delete(&Category{}); result.Error != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
    	return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id":id,"message": "Category deleted successfully"})
}

func main() {
	mux := http.NewServeMux()
	var port = "8000"
	db, err := gorm.Open(mysql.Open("root@tcp(localhost)/db_ujikom"), &gorm.Config{}) 
	if err != nil {
		panic("Failed Connect to Database!")
	}
	fmt.Println("Success Connect to Database!") 
 
	mux.HandleFunc("/", handleIndex)

	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		handleCategories(w, r, db)
	})

	mux.HandleFunc("/category/store", func(w http.ResponseWriter, r *http.Request) {
		handleCreateCategory(w, r, db)
	})

	mux.HandleFunc("/category/{id}", func(w http.ResponseWriter, r * http.Request) {
		if r.Method == http.MethodPut {
			handleEditCategory(w, r, db)
		} else if r.Method == http.MethodDelete {
			handleDeleteCategory(w, r, db)
		}
	})

	fmt.Printf("Server Started at Port %s\n", port)
	http.ListenAndServe(":" + port, mux)
}