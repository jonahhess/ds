package errors

import (
	"log"
	"net/http"
)

// HandleError writes an error response with appropriate status code
// It logs the error but doesn't expose internal details to the client
func HandleError(w http.ResponseWriter, r *http.Request, err error, statusCode int, userMessage string) {
	log.Printf("Error [%s %s]: %v", r.Method, r.URL.Path, err)
	http.Error(w, userMessage, statusCode)
}

// HandleInternalError logs the error and returns a generic 500 error
func HandleInternalError(w http.ResponseWriter, r *http.Request, err error) {
	HandleError(w, r, err, http.StatusInternalServerError, "Internal server error")
}

// HandleNotFound returns a 404 error
func HandleNotFound(w http.ResponseWriter, r *http.Request, resource string) {
	http.Error(w, resource+" not found", http.StatusNotFound)
}

// HandleBadRequest returns a 400 error
func HandleBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

// HandleUnauthorized returns a 401 error
func HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

// HandleForbidden returns a 403 error
func HandleForbidden(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "forbidden", http.StatusForbidden)
}
