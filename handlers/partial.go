package handlers

import (
	"html/template"
	"myapp/auth"
	"myapp/db"
	"net/http"
	"path/filepath"
	"slices"
)

func PartialHandler(w http.ResponseWriter, r *http.Request) {
    claims, err := auth.VerifyToken(r)
    if err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    partial := r.URL.Path[len("/load-partial/"):]
    allowed := slices.Contains(claims.ValidPartials, partial)
    if !allowed {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Rotate token
    auth.SetTokenCookieWithClaims(w, claims)

    tmplPath := filepath.Join("templates", "partials", partial+".html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Not Found", http.StatusNotFound)
        return
    }

    // Fetch user info from DB
    user, err := db.GetUserByID(claims.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    tmpl.Execute(w, user)
}
