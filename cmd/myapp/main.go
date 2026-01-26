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
	"github.com/jonahhess/ds/internal/db"
	"github.com/jonahhess/ds/internal/router"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authKey := envs["AUTH_KEY"]
	encKey := envs["ENC_KEY"]
	path := envs["MYAPP_DB_PATH"]

	if err := db.InitDB(path); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	if err := db.CreateTables(); err != nil {
		log.Printf("table creation failed: %v", err)
	}

	defer db.CloseDB()

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

	r := router.SetupRoutes(store, db.DB)

	// --- server ---
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server started on :8080")
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
