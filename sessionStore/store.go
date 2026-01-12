package sessions

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

func InitStore() *sessions.CookieStore {
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

	return store
}

func GetStore() *sessions.CookieStore {
	return store
}
