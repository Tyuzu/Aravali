package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

// // SearchResult represents a single search result item
// type SearchResult struct {
// 	ID          string    `json:"id,omitempty"`
// 	EntityID    string    `json:"entityid,omitempty"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Image       string    `json:"image,omitempty"`
// 	CreatedAt   time.Time `json:"createdAt"`
// }

// // AllSearchResults represents results grouped by entity type
// type AllSearchResults struct {
// 	Events       []SearchResult `json:"events,omitempty"`
// 	Places       []SearchResult `json:"places,omitempty"`
// 	Feedposts    []SearchResult `json:"feedposts,omitempty"`
// 	Merch        []SearchResult `json:"merch,omitempty"`
// 	Blogposts    []SearchResult `json:"blogposts,omitempty"`
// 	Farms        []SearchResult `json:"farms,omitempty"`
// 	Songs        []SearchResult `json:"songs,omitempty"`
// 	Users        []SearchResult `json:"users,omitempty"`
// 	Recipes      []SearchResult `json:"recipes,omitempty"`
// 	Products     []SearchResult `json:"products,omitempty"`
// 	Menu         []SearchResult `json:"menu,omitempty"`
// 	Media        []SearchResult `json:"media,omitempty"`
// 	Crops        []SearchResult `json:"crops,omitempty"`
// 	Baitoworkers []SearchResult `json:"baitoworkers,omitempty"`
// 	Baitos       []SearchResult `json:"baitos,omitempty"`
// 	Artists      []SearchResult `json:"artists,omitempty"`
// }

// GetAutocompleteSuggestions returns autocomplete suggestions for a prefix using SQL
func SQLGetAutocompleteSuggestions(ctx context.Context, app *infra.Deps, prefix string) ([]string, error) {
	prefix = strings.ToLower(prefix)
	pattern := prefix + "%"

	suggestions := make(map[string]bool)

	type tableInfo struct {
		table string
		field string
	}

	tables := []tableInfo{
		{config.Tables.EventsTable, "title"},
		{config.Tables.PlacesTable, "name"},
		{config.Tables.FeedPostsTable, "title"},
		{config.Tables.MerchTable, "name"},
		{config.Tables.BlogPostsTable, "title"},
		{config.Tables.FarmsTable, "name"},
		{config.Tables.SongsTable, "title"},
		{config.Tables.UserTable, "username"},
		{config.Tables.RecipeTable, "title"},
		{config.Tables.ProductTable, "name"},
		{config.Tables.MenuTable, "name"},
		{config.Tables.MediaTable, "title"},
		{config.Tables.CropsTable, "name"},
		{config.Tables.BaitoWorkerTable, "name"},
		{config.Tables.BaitoTable, "title"},
		{config.Tables.ArtistsTable, "name"},
	}

	for _, tbl := range tables {
		var results []map[string]any
		where := fmt.Sprintf("LOWER(%s) LIKE $1", tbl.field)
		args := []any{pattern}

		err := app.SQLDB.FindMany(ctx, tbl.table, where, args, &results)
		if err != nil {
			continue // Skip individual table errors
		}

		for _, result := range results {
			if val, ok := result[tbl.field]; ok {
				if str, ok := val.(string); ok && str != "" {
					suggestions[str] = true
				}
			}
		}
	}

	output := make([]string, 0, len(suggestions))
	for k := range suggestions {
		output = append(output, k)
	}

	return output, nil
}

// SearchByEntity searches in a specific entity type using SQL
func SQLSearchByEntity(ctx context.Context, app *infra.Deps, entityType, query string) ([]SearchResult, error) {
	var results []SearchResult
	pattern := "%" + strings.ToLower(query) + "%"

	entityInfo := map[string]struct {
		table      string
		fields     []string
		titleField string
	}{
		"events": {
			table:      config.Tables.EventsTable,
			fields:     []string{"title", "description", "location"},
			titleField: "title",
		},
		"places": {
			table:      config.Tables.PlacesTable,
			fields:     []string{"name", "description", "address"},
			titleField: "name",
		},
		"feedposts": {
			table:      config.Tables.FeedPostsTable,
			fields:     []string{"title", "description", "caption"},
			titleField: "title",
		},
		"merch": {
			table:      config.Tables.MerchTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"blogposts": {
			table:      config.Tables.BlogPostsTable,
			fields:     []string{"title", "description", "content"},
			titleField: "title",
		},
		"farms": {
			table:      config.Tables.FarmsTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"songs": {
			table:      config.Tables.SongsTable,
			fields:     []string{"title", "description", "artist"},
			titleField: "title",
		},
		"users": {
			table:      config.Tables.UserTable,
			fields:     []string{"username", "displayname", "bio"},
			titleField: "username",
		},
		"recipes": {
			table:      config.Tables.RecipeTable,
			fields:     []string{"title", "description", "ingredients"},
			titleField: "title",
		},
		"products": {
			table:      config.Tables.ProductTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"menu": {
			table:      config.Tables.MenuTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"media": {
			table:      config.Tables.MediaTable,
			fields:     []string{"title", "description"},
			titleField: "title",
		},
		"crops": {
			table:      config.Tables.CropsTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"baitoworkers": {
			table:      config.Tables.BaitoWorkerTable,
			fields:     []string{"name", "description"},
			titleField: "name",
		},
		"baitos": {
			table:      config.Tables.BaitoTable,
			fields:     []string{"title", "description"},
			titleField: "title",
		},
		"artists": {
			table:      config.Tables.ArtistsTable,
			fields:     []string{"name", "description", "bio"},
			titleField: "name",
		},
	}

	info, exists := entityInfo[entityType]
	if !exists {
		return results, nil
	}

	// Construct SQL WHERE clause for multi-field case-insensitive search
	var whereConditions []string
	args := []any{pattern}
	for _, field := range info.fields {
		whereConditions = append(whereConditions, fmt.Sprintf("LOWER(%s) LIKE $1", field))
	}
	where := "(" + strings.Join(whereConditions, " OR ") + ")"

	var docs []map[string]any
	opts := sqldb.FindManyOptions{
		Limit: 50,
	}

	err := app.SQLDB.FindManyWithOptions(ctx, info.table, where, args, opts, &docs)
	if err != nil {
		return results, err
	}

	for _, doc := range docs {
		result := SearchResult{
			Title: getStringField(doc, info.titleField),
		}

		if desc, ok := doc["description"].(string); ok {
			result.Description = desc
		} else if desc, ok := doc["bio"].(string); ok {
			result.Description = desc
		} else if desc, ok := doc["caption"].(string); ok {
			result.Description = desc
		} else if desc, ok := doc["content"].(string); ok && len(desc) > 200 {
			result.Description = desc[:200] + "..."
		}

		if id, ok := doc["id"].(string); ok {
			result.ID = id
		}
		if id, ok := doc["entityid"].(string); ok {
			result.EntityID = id
		}
		if id, ok := doc["eventid"].(string); ok {
			result.EntityID = id
		}
		if id, ok := doc["placeid"].(string); ok {
			result.EntityID = id
		}
		if id, ok := doc["userid"].(string); ok {
			result.EntityID = id
		}

		if img, ok := doc["image"].(string); ok {
			result.Image = img
		} else if img, ok := doc["banner"].(string); ok {
			result.Image = img
		} else if img, ok := doc["banner_image"].(string); ok {
			result.Image = img
		}

		if t, ok := doc["created_at"].(time.Time); ok {
			result.CreatedAt = t
		} else if t, ok := doc["createdAt"].(time.Time); ok {
			result.CreatedAt = t
		}

		results = append(results, result)
	}

	return results, nil
}

// SearchAll searches across all entity types and returns grouped results
func SQL(ctx context.Context, app *infra.Deps, query string) (AllSearchResults, error) {
	allResults := AllSearchResults{}

	entityTypes := []string{
		"events", "places", "feedposts", "merch", "blogposts",
		"farms", "songs", "users", "recipes", "products",
		"menu", "media", "crops", "baitoworkers", "baitos", "artists",
	}

	resultsChan := make(chan struct {
		entityType string
		results    []SearchResult
	}, len(entityTypes))

	for _, entityType := range entityTypes {
		go func(et string) {
			results, _ := SearchByEntity(ctx, app, et, query)
			if len(results) > 10 {
				results = results[:10]
			}
			resultsChan <- struct {
				entityType string
				results    []SearchResult
			}{et, results}
		}(entityType)
	}

	for range entityTypes {
		res := <-resultsChan
		switch res.entityType {
		case "events":
			allResults.Events = res.results
		case "places":
			allResults.Places = res.results
		case "feedposts":
			allResults.Feedposts = res.results
		case "merch":
			allResults.Merch = res.results
		case "blogposts":
			allResults.Blogposts = res.results
		case "farms":
			allResults.Farms = res.results
		case "songs":
			allResults.Songs = res.results
		case "users":
			allResults.Users = res.results
		case "recipes":
			allResults.Recipes = res.results
		case "products":
			allResults.Products = res.results
		case "menu":
			allResults.Menu = res.results
		case "media":
			allResults.Media = res.results
		case "crops":
			allResults.Crops = res.results
		case "baitoworkers":
			allResults.Baitoworkers = res.results
		case "baitos":
			allResults.Baitos = res.results
		case "artists":
			allResults.Artists = res.results
		}
	}

	return allResults, nil
}

// Helper function to safely get string field from document
func SQLgetStringField(doc map[string]any, fieldName string) string {
	if val, ok := doc[fieldName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
