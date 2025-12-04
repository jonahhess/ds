package main

import (
	"fmt"
	"myapp/db"
	"myapp/handlers"
	"net/http"
)

func main() {
    // Initialize SQLite database
    err := db.InitDB("db/myapp.db")
    if err != nil {
        panic(fmt.Sprintf("Database init failed: %v", err))
    }
    db.CreateTables()

    http.HandleFunc("/", handlers.RootHandler)
    http.HandleFunc("/load-partial/", handlers.PartialHandler)
    http.HandleFunc("/dashboard", handlers.DashboardHandler)

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
