package main

import (
	"fmt"
	"myapp/db"
	"myapp/handlers"
	"net/http"
)

func main() {
    err := db.InitDB("db/myapp.db")
    if err != nil {
        panic(fmt.Sprintf("Database init failed: %v", err))
    }
    db.CreateTables()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        handlers.RootHandler(w, r)
    })
    http.HandleFunc("/load-partial/", handlers.PartialHandler)
    http.HandleFunc("/dashboard", handlers.DashboardHandler)

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))


    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
