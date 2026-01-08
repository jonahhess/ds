package main

import (
	"fmt"
	"myapp/db"
	"myapp/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
    err := db.InitDB("db/myapp.db")
    if err != nil {
        panic(fmt.Sprintf("Database init failed: %v", err))
    }
    db.CreateTables()

    r := chi.NewRouter()

    r.Get("/", handlers.RootHandler)
    r.Get("/dashboard", handlers.DashboardHandler)
    r.Get("/reports", handlers.ReportsHandler)

    fs := http.FileServer(http.Dir("./static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fs))

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", r)
}
