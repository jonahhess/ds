package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const csrfTokenKey = "csrf_token"
const csrfTokenLength = 32

// CSRFMiddleware creates CSRF protection middleware
func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sess, ok := SessionFromContext(ctx)
			if !ok {
				http.Error(w, "Session not found", http.StatusInternalServerError)
				return
			}

			// Generate CSRF token if not present
			token, ok := sess.Values[csrfTokenKey].(string)
			if !ok || token == "" {
				token = generateCSRFToken()
				sess.Values[csrfTokenKey] = token
				if err := sess.Save(r, w); err != nil {
					http.Error(w, "Failed to save session", http.StatusInternalServerError)
					return
				}
			}

			// Store CSRF token in context for template access
			ctx = context.WithValue(ctx, "csrf_token", token)
			r = r.WithContext(ctx)

			// Skip CSRF check for safe methods
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate CSRF token for unsafe methods
			submittedToken := r.FormValue("csrf_token")
			if submittedToken == "" {
				// Try header as fallback
				submittedToken = r.Header.Get("X-CSRF-Token")
			}

			if submittedToken != token {
				http.Error(w, "CSRF token mismatch", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFToken retrieves the CSRF token from the session
func CSRFToken(r *http.Request) string {
	sess, ok := SessionFromContext(r.Context())
	if !ok {
		return ""
	}
	token, _ := sess.Values[csrfTokenKey].(string)
	return token
}

// CSRFTokenFromContext retrieves the CSRF token from context
func CSRFTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value("csrf_token").(string); ok {
		return token
	}
	return ""
}

func generateCSRFToken() string {
	b := make([]byte, csrfTokenLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
