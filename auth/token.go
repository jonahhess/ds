package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSecret = []byte("supersecretkey")

type TokenClaims struct {
    ValidPartials []string `json:"valid_partials"`
    ValidPages    []string `json:"valid_pages"`
    jwt.RegisteredClaims
}

// GenerateToken creates a new signed token
func GenerateToken() (string, error) {
    claims := TokenClaims{
        ValidPartials: []string{"user-profile", "settings-form"},
        ValidPages:    []string{"dashboard", "reports"},
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(hmacSecret)
}

// VerifyToken validates the token from the cookie
func VerifyToken(r *http.Request) (*TokenClaims, error) {
    cookie, err := r.Cookie("app-token")
    if err != nil {
        return nil, err
    }
    token, err := jwt.ParseWithClaims(cookie.Value, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
        return hmacSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    claims, ok := token.Claims.(*TokenClaims)
    if !ok {
        return nil, fmt.Errorf("invalid claims")
    }
    return claims, nil
}

// SetTokenCookie issues a new token cookie
func SetTokenCookie(w http.ResponseWriter) error {
    tokenStr, err := GenerateToken()
    if err != nil {
        return err
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "app-token",
        Value:    tokenStr,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // true in production with HTTPS
    })
    return nil
}
