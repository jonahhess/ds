package main

import (
	"fmt"
	"myapp/db"
	"myapp/router"
	"net/http"
)

func main() {
	err := db.InitDB("db/myapp.db")
	if err != nil {
		panic(fmt.Sprintf("Database init failed: %v", err))
	}
	db.CreateTables()

	r := router.SetupRoutes()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
