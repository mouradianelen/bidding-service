package main

import (
	"fmt"
	"main/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.BidRequestHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
