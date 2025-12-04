package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSecret = []byte("supersecretkey")

// TokenClaims holds all user info stored in the token.
type TokenClaims struct {
    UserID        string   `json:"user_id"`
    Name          string   `json:"name"`
    Email         string   `json:"email"`
    ValidPartials []string `json:"valid_partials"`
    ValidPages    []string `json:"valid_pages"`
    jwt.RegisteredClaims
}

// GenerateTokenForUser creates a fresh token for a user (used for login).
func GenerateTokenForUser(userID, name, email string) (string, error) {
    claims := TokenClaims{
        UserID:        userID,
        Name:          name,
        Email:         email,
        ValidPartials: []string{"user-profile", "settings-form"},
        ValidPages:    []string{"dashboard", "reports"},
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(hmacSecret)
}

// VerifyToken parses and validates the token from the "app-token" cookie.
func VerifyToken(r *http.Request) (*TokenClaims, error) {
    cookie, err := r.Cookie("app-token")
    if err != nil {
        return nil, err
    }
    token, err := jwt.ParseWithClaims(cookie.Value, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
        return hmacSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    claims, ok := token.Claims.(*TokenClaims)
    if !ok {
        return nil, err
    }
    return claims, nil
}

// SetTokenCookie issues a brand new token cookie (no existing claims).
func SetTokenCookie(w http.ResponseWriter) error {
    tokenStr, err := GenerateTokenForUser("anon", "Anonymous", "anon@example.com")
    if err != nil {
        return err
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "app-token",
        Value:    tokenStr,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set to true in production with HTTPS
    })
    return nil
}

// SetTokenCookieWithClaims rotates/refreshes a token using existing claims.
// Expects a pointer to TokenClaims, refreshes expiry and sets cookie.
func SetTokenCookieWithClaims(w http.ResponseWriter, claims *TokenClaims) error {
    // Refresh expiry
    claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(5 * time.Minute))

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString(hmacSecret)
    if err != nil {
        return err
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "app-token",
        Value:    tokenStr,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set true in production with HTTPS
    })
    return nil
}
