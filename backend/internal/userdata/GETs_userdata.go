package userdata

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/config"
	"scav/infra"
	"scav/middleware"
	"scav/utils"
	log "scav/utils/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GetUserProfileData fetches user-specific entity data
func GetUserProfileData(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := utils.GetParam(r, "username")

		// Validate JWT
		tokenString := r.Header.Get("Authorization")
		claims, err := middleware.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if username != claims.UserID && username != claims.Username {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse query parameter
		entityType := r.URL.Query().Get("entity_type")
		if entityType == "" {
			http.Error(w, "Entity type is required", http.StatusBadRequest)
			return
		}
		if !IsValidEntityType(entityType) {
			http.Error(w, "Invalid entity type", http.StatusBadRequest)
			return
		}

		// Fetch data from Database interface
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var results []UserData
		filter := map[string]any{
			"entity_type": entityType,
			"userid":      username,
		}

		if err := app.DB.FindMany(ctx, userdataCollection, filter, &results); err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			log.Printf("Error fetching user data: %v", err)
			return
		}

		// Ensure empty slice instead of nil
		if results == nil {
			results = []UserData{}
		}

		// Respond with JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding response: %v", err)
			return
		}
	}
}

func GetOtherUserProfileData(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := utils.GetParam(r, "username")
		if username == "" {
			http.Error(w, "Missing user identifier", http.StatusBadRequest)
			return
		}

		entityType := r.URL.Query().Get("entity_type")
		if entityType == "" {
			http.Error(w, "Entity type is required", http.StatusBadRequest)
			return
		}
		if entityType != "feedpost" {
			http.Error(w, "Invalid entity type", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		type postDoc struct {
			PostID    string    `bson:"postid"`
			Title     string    `bson:"title"`
			Thumb     string    `bson:"thumb"`
			CreatedBy string    `bson:"createdBy"`
			Username  string    `bson:"username"`
			CreatedAt time.Time `bson:"createdAt"`
			Blocks    []struct {
				Type    string `bson:"type"`
				URL     string `bson:"url"`
				Caption string `bson:"caption"`
			} `bson:"blocks"`
		}

		var posts []postDoc
		filter := bson.M{
			"$or": []bson.M{
				{"createdBy": username},
				{"username": username},
			},
		}

		if err := app.DB.FindMany(ctx, config.Collections.FeedPostsCollection, filter, &posts); err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			log.Printf("Error fetching other user feed posts: %v", err)
			return
		}

		response := make([]map[string]any, 0, len(posts))
		for _, post := range posts {
			imageURL := post.Thumb
			caption := post.Title
			if imageURL == "" && len(post.Blocks) > 0 {
				for _, block := range post.Blocks {
					if block.URL != "" {
						imageURL = block.URL
						if block.Caption != "" {
							caption = block.Caption
						}
						break
					}
				}
			}
			if caption == "" {
				caption = "Post image"
			}

			response = append(response, map[string]any{
				"id":          post.PostID,
				"entity_id":   post.PostID,
				"postid":      post.PostID,
				"entity_type": "feedpost",
				"image_url":   imageURL,
				"caption":     caption,
				"created_at":  post.CreatedAt.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding other user feed posts: %v", err)
		}
	}
}
