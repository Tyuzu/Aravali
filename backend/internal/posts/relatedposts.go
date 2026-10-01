// File: internal/posts/relatedposts.go

package posts

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/infra/sqldb"
	"scav/utils"
	"time"
)

type Post struct {
	PostID      string    `bson:"postid" json:"postid"`
	Title       string    `bson:"title" json:"title"`
	Category    string    `bson:"category" json:"category"`
	Subcategory string    `bson:"subcategory" json:"subcategory"`
	Tags        []string  `bson:"tags" json:"tags"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy   string    `bson:"createdBy" json:"createdBy"`
}

func GetRelatedPosts(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID := r.URL.Query().Get("postid")
		category := r.URL.Query().Get("category")
		subcategory := r.URL.Query().Get("subcategory")
		tags := r.URL.Query()["tags"]

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := "postid <> $1"
		args := []any{postID}
		if category != "" {
			query += " AND category = $2"
			args = append(args, category)
		}
		if subcategory != "" {
			query += " AND subcategory = $3"
			args = append(args, subcategory)
		}
		if len(tags) > 0 {
			query += " AND tags && $4"
			args = append(args, tags)
		}

		opts := sqldb.FindManyOptions{
			Limit:   10,
			OrderBy: "created_at DESC",
		}

		var related []Post
		if err := FindRelatedPostsWithOptions(ctx, app, query, args, opts, &related); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if related == nil {
			related = []Post{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"related": related,
		})
	}
}
