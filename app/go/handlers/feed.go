package handlers

import (
	"nae/models"
	"nae/utils"
	"net/http"
)

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	utils.SendJSON(w, http.StatusOK, models.FeedData{
		Title: "Golang Performance Feed",
		Body:  "Serving high-throughput data to your Solar2D application with sub-millisecond response times.",
		Tags:  []string{"Go", "Solar2D", "REST"},
	})
}
