package handlers

import (
	"html/template"
	"myapp/auth"
	"net/http"
	"path/filepath"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    _, err := auth.VerifyToken(r)
    if err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    if err := auth.SetTokenCookie(w); err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    tmplPath := filepath.Join("templates", "dashboard.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}
