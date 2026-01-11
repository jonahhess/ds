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
	"github.com/michaeljs1990/sqlitestore"

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
	db.CreateUsersTable()
	db.CreateSessionsTable()

	// --- sessions (SQLite-backed) ---
	sessionStore, err := sqlitestore.NewSqliteStore(
		"./db/myapp.db",  // SQLite file
		"sessions",      // table name
		"/",             // path
		86400*30,        // max age
		[]byte("dev-secret-change-me"),
	)
	if err != nil {
		log.Fatalf("session store init failed: %v", err)
	}

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true,
	}

	// --- router ---
	r := router.SetupRoutes(ctx, sessionStore)

	// --- server ---
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
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
