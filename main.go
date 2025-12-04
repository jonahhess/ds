package main

import (
	"fmt"
	"net/http"

	"myapp/handlers"
)

func main() {
    http.HandleFunc("/", handlers.RootHandler)
    http.HandleFunc("/load-partial/", handlers.PartialHandler)
    http.HandleFunc("/dashboard", handlers.DashboardHandler)

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
