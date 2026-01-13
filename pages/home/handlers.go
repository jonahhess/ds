package home

import (
	"fmt"
	"myapp/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Home Page Cookies:", r.Cookies())

	err := layouts.
		Base("Home", Home()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
