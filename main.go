package main

import (
	"fmt"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	var port = "8000"
	http.HandleFunc("/", handleIndex)


	
	fmt.Printf("Server Started at Port %s\n", port)
	http.ListenAndServe(":" + port, nil)
}