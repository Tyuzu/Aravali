// File: internal/posts/relatedposts.go

package posts

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

type Post struct {
	PostID      string    `db:"postid" json:"postid"`
	Title       string    `db:"title" json:"title"`
	Category    string    `db:"category" json:"category"`
	Subcategory string    `db:"subcategory" json:"subcategory"`
	Tags        []string  `db:"tags" json:"tags"`
	CreatedAt   time.Time `db:"createdAt" json:"createdAt"`
	CreatedBy   string    `db:"createdBy" json:"createdBy"`
}

func GetRelatedPosts(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID := r.URL.Query().Get("postid")
		category := r.URL.Query().Get("category")
		subcategory := r.URL.Query().Get("subcategory")
		tags := r.URL.Query()["tags"]

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		related, err := GetRelatedPostsPage(ctx, app, postID, category, subcategory, tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"related": related,
		})
	}
}
