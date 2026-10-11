// File: internal/search/searchSQLDB.go

package search

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

type SearchResult struct {
	ID          string    `json:"id,omitempty"`
	EntityID    string    `json:"entityid,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AllSearchResults struct {
	Events       []SearchResult `json:"events,omitempty"`
	Places       []SearchResult `json:"places,omitempty"`
	Feedposts    []SearchResult `json:"feedposts,omitempty"`
	Merch        []SearchResult `json:"merch,omitempty"`
	Blogposts    []SearchResult `json:"blogposts,omitempty"`
	Farms        []SearchResult `json:"farms,omitempty"`
	Songs        []SearchResult `json:"songs,omitempty"`
	Users        []SearchResult `json:"users,omitempty"`
	Recipes      []SearchResult `json:"recipes,omitempty"`
	Products     []SearchResult `json:"products,omitempty"`
	Menu         []SearchResult `json:"menu,omitempty"`
	Media        []SearchResult `json:"media,omitempty"`
	Crops        []SearchResult `json:"crops,omitempty"`
	Baitoworkers []SearchResult `json:"baitoworkers,omitempty"`
	Baitos       []SearchResult `json:"baitos,omitempty"`
	Artists      []SearchResult `json:"artists,omitempty"`
}

func normalizeSearchQuery(q string) string {
	return strings.ToLower(strings.TrimSpace(q))
}

func readNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// GetAutocompleteSuggestions returns autocomplete suggestions for a prefix using SQL
func GetAutocompleteSuggestions(ctx context.Context, app *infra.Deps, prefix string) ([]string, error) {
	if app == nil || app.SQLDB == nil {
		return []string{}, nil
	}
	prefix = normalizeSearchQuery(prefix)
	if prefix == "" {
		return []string{}, nil
	}
	pattern := prefix + "%"
	suggestions := map[string]struct{}{}
	type tableInfo struct {
		table string
		field string
	}
	// A small, safe set of searchable columns across the app.
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
		query := fmt.Sprintf(`SELECT DISTINCT LOWER(%s) FROM %s WHERE LOWER(%s) LIKE $1 ORDER BY 1 LIMIT 20`, tbl.field, tbl.table, tbl.field)
		rows, err := app.SQLDB.Query(ctx, query, pattern)
		if err != nil {
			continue
		}
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err == nil && strings.TrimSpace(value) != "" {
				suggestions[value] = struct{}{}
			}
		}
		rows.Close()
	}
	output := make([]string, 0, len(suggestions))
	for k := range suggestions {
		output = append(output, k)
	}
	sort.Strings(output)
	return output, nil
}

// SearchByEntity searches in a specific entity type using SQL
func SearchByEntity(ctx context.Context, app *infra.Deps, entityType, query string) ([]SearchResult, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	results := make([]SearchResult, 0)
	pattern := "%" + normalizeSearchQuery(query) + "%"

	entityInfo := map[string]struct {
		table      string
		columns    []string
		titleField string
	}{
		"events":       {table: config.Tables.EventsTable, columns: []string{"id", "eventid", "title", "description", "image", "created_at"}, titleField: "title"},
		"places":       {table: config.Tables.PlacesTable, columns: []string{"id", "placeid", "name", "description", "image", "created_at"}, titleField: "name"},
		"feedposts":    {table: config.Tables.FeedPostsTable, columns: []string{"id", "feedpostid", "title", "description", "image", "created_at"}, titleField: "title"},
		"merch":        {table: config.Tables.MerchTable, columns: []string{"id", "merchid", "name", "description", "image", "created_at"}, titleField: "name"},
		"blogposts":    {table: config.Tables.BlogPostsTable, columns: []string{"id", "blogpostid", "title", "description", "image", "created_at"}, titleField: "title"},
		"farms":        {table: config.Tables.FarmsTable, columns: []string{"id", "farmid", "name", "description", "image", "created_at"}, titleField: "name"},
		"songs":        {table: config.Tables.SongsTable, columns: []string{"id", "songid", "title", "description", "image", "created_at"}, titleField: "title"},
		"users":        {table: config.Tables.UserTable, columns: []string{"id", "userid", "username", "bio", "image", "created_at"}, titleField: "username"},
		"recipes":      {table: config.Tables.RecipeTable, columns: []string{"recipeid", "recipeid", "title", "description", "banner", "created_at"}, titleField: "title"},
		"products":     {table: config.Tables.ProductTable, columns: []string{"id", "productid", "name", "description", "image", "created_at"}, titleField: "name"},
		"menu":         {table: config.Tables.MenuTable, columns: []string{"id", "menuid", "name", "description", "image", "created_at"}, titleField: "name"},
		"media":        {table: config.Tables.MediaTable, columns: []string{"id", "mediaid", "title", "description", "image", "created_at"}, titleField: "title"},
		"crops":        {table: config.Tables.CropsTable, columns: []string{"id", "cropid", "name", "description", "image", "created_at"}, titleField: "name"},
		"baitoworkers": {table: config.Tables.BaitoWorkerTable, columns: []string{"id", "baitoworkerid", "name", "description", "image", "created_at"}, titleField: "name"},
		"baitos":       {table: config.Tables.BaitoTable, columns: []string{"id", "baitoid", "title", "description", "image", "created_at"}, titleField: "title"},
		"artists":      {table: config.Tables.ArtistsTable, columns: []string{"id", "artistid", "name", "description", "image", "created_at"}, titleField: "name"},
	}
	info, exists := entityInfo[entityType]
	if !exists || info.table == "" {
		return results, nil
	}
	cond := make([]string, 0, len(info.columns)-2)
	for _, field := range info.columns[2 : len(info.columns)-1] {
		cond = append(cond, fmt.Sprintf("LOWER(COALESCE(%s, '')) LIKE $1", field))
	}
	querySQL := fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY created_at DESC LIMIT 20`,
		strings.Join(info.columns, ", "),
		info.table,
		strings.Join(cond, " OR "),
	)
	rows, err := app.SQLDB.Query(ctx, querySQL, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, entityID, title, description, image sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &entityID, &title, &description, &image, &createdAt); err != nil {
			continue
		}
		result := SearchResult{
			ID:          readNullString(id),
			EntityID:    readNullString(entityID),
			Title:       readNullString(title),
			Description: readNullString(description),
			Image:       readNullString(image),
			CreatedAt:   createdAt,
		}
		if result.ID == "" {
			result.ID = result.EntityID
		}
		if result.Title == "" && result.EntityID != "" {
			result.Title = result.EntityID
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

// SearchAll searches across all entity types and returns grouped results
func SearchAll(ctx context.Context, app *infra.Deps, query string) (AllSearchResults, error) {
	allResults := AllSearchResults{}
	entityTypes := []string{"events", "places", "feedposts", "merch", "blogposts", "farms", "songs", "users", "recipes", "products", "menu", "media", "crops", "baitoworkers", "baitos", "artists"}
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
			}{entityType: et, results: results}
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

func getStringField(doc map[string]any, fieldName string) string {
	if val, ok := doc[fieldName]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
