package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gojp/goreportcard/database"
)

type HomeHandler struct {
	DB database.Database
}

// HomeHandler handles the homepage
func (h *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {

		recentRepos, err := h.DB.GetMostRecentlyViewed(5)
		if err != nil {
			log.Println("ERROR: while calling GetMostRecentlyViewed:", err)
			recentRepos = []string{}
		}
		t := template.Must(template.New("home.html").Delims("[[", "]]").ParseFiles("templates/home.html", "templates/footer.html"))
		t.Execute(w, map[string]interface{}{
			"Recent":               recentRepos,
			"google_analytics_key": googleAnalyticsKey,
		})

		return
	}

	errorHandler(w, r, http.StatusNotFound)
}
