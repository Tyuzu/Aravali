// File: internal/posts/GETs_blogposts.go

package posts

import (
	"fmt"
	"net/http"
	"scav/infra"
	"scav/utils"
)

// --- Get single post ---
func GetPost(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		postID := utils.GetParam(r, "id")

		var post BlogPost
		if err := GetPostByID(ctx, app, postID, &post); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"post": post,
		})
	}
}

func GetAllPosts(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		limit := 20
		page := 1

		if l := query.Get("limit"); l != "" {
			_, _ = fmt.Sscanf(l, "%d", &limit)
		}
		if p := query.Get("page"); p != "" {
			_, _ = fmt.Sscanf(p, "%d", &page)
		}

		posts, err := GetPostsPage(ctx, app, limit, page)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch posts")
			return
		}

		// --- Collect unique user IDs ---
		userIDs := make(map[string]struct{})
		for _, p := range posts {
			if p.CreatedBy != "" {
				userIDs[p.CreatedBy] = struct{}{}
			}
		}

		ids := make([]string, 0, len(userIDs))
		for id := range userIDs {
			ids = append(ids, id)
		}

		// --- Fetch usernames ---
		usernames := make(map[string]string)
		if len(ids) > 0 {
			var users []struct {
				UserID   string `db:"userid"`
				Username string `db:"username"`
			}

			if err := FindUsersByFilter(ctx, app, "userid = ANY($1)", []any{ids}, &users); err == nil {
				for _, u := range users {
					usernames[u.UserID] = u.Username
				}
			}
		}

		// --- Build response ---
		resp := make([]BlogPostResponse, 0, len(posts))
		for _, p := range posts {
			resp = append(resp, BlogPostResponse{
				PostID:      p.PostID,
				Title:       p.Title,
				Category:    p.Category,
				Subcategory: p.Subcategory,
				ReferenceID: p.ReferenceID,
				Thumb:       pickThumb(p.Blocks),
				CreatedBy:   p.CreatedBy,
				Username:    usernames[p.CreatedBy],
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			})
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"posts": resp,
		})
	}
}
