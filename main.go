package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"myapp/db"
	"myapp/router"
	sessions "myapp/sessionStore"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := db.InitDB("db/myapp.db"); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	if err := db.CreateTables(); err != nil {
		log.Printf("table creation failed: %v", err)
	}

	defer db.CloseDB()

	sess := sessions.InitStore()
	r := router.SetupRoutes(sess)

	// --- server ---
	srv := &http.Server{
		Addr:         ":8088",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server started on :8088")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\nShutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	fmt.Println("Server stopped cleanly")
}
