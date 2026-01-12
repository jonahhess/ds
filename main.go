package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	"myapp/db"
	"myapp/router"
)

func main() {
	// --- root context (ctrl+c / docker stop) ---
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// --- database ---
	if err := db.InitDB("db/myapp.db"); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	if err := db.CreateTables(); err != nil {
		log.Fatalf("table creation failed: %v", err)
	}

	defer db.CloseDB()

	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authKey := envs["AUTH_KEY"]
	encKey := envs["ENC_KEY"]

	store := sessions.NewCookieStore(
		[]byte(authKey),
		[]byte(encKey),
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false, // false only for localhost
		SameSite: http.SameSiteLaxMode,
	}

	// --- router ---
	r := router.SetupRoutes(store)

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

	// --- graceful shutdown ---
	<-ctx.Done()
	fmt.Println("\nShutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	fmt.Println("Server stopped cleanly")
}
