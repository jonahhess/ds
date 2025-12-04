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

  // Parse base layout + page + JS component
    tmpl, err := template.ParseFiles(
        filepath.Join("templates", "layouts", "base.html"),
        filepath.Join("templates", "dashboard.html"),
    )
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    if err := tmpl.ExecuteTemplate(w, "base.html", map[string]string{
        "Title": "Dashboard",
    }); err != nil {
        http.Error(w, "Template execution error", http.StatusInternalServerError)
        return
    }
}
