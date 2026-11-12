package crops

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"scav/infra"
	"scav/internal/farms"
	"scav/utils"
)

/* ---------------------------------------------------- */
/* Filtered Crops                                       */
/* ---------------------------------------------------- */

func GetFilteredCrops(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		params := r.URL.Query()
		filter := map[string]any{}

		if c := params.Get("category"); c != "" {
			filter["category"] = c
		}
		if region := params.Get("region"); region != "" {
			filter["farmLocation"] = region
		}
		if params.Get("inStock") == "true" {
			filter["quantity"] = map[string]any{"$gt": 0}
		}

		price := map[string]any{}
		if min := utils.ParseFloat(params.Get("minPrice")); min > 0 {
			price["$gte"] = min
		}
		if max := utils.ParseFloat(params.Get("maxPrice")); max > 0 {
			price["$lte"] = max
		}
		if len(price) > 0 {
			filter["price"] = price
		}

		crops, err := findFilteredCrops(ctx, app.DB, filter)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch crops")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"crops":   crops,
		})
	}
}

/* ---------------------------------------------------- */
/* Pre Crop Catalogue (Redis → DB → CSV)                */
/* ---------------------------------------------------- */

func GetPreCropCatalogue(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		const cacheKey = "crop_catalogue"
		var crops []farms.CropCatalogueItem

		/* ------------------------------------------------ */
		/* 1. Cache                                         */
		/* ------------------------------------------------ */

		if data, err := app.Cache.Get(ctx, cacheKey); err == nil && len(data) > 0 {
			if json.Unmarshal(data, &crops) == nil {
				utils.RespondWithJSON(w, http.StatusOK, utils.M{
					"success": true,
					"crops":   crops,
				})
				return
			}
		}

		/* ------------------------------------------------ */
		/* 2. Database                                      */
		/* ------------------------------------------------ */

		if items, err := findCatalogueItems(ctx, app.DB, map[string]any{}); err == nil && len(items) > 0 {
			crops = items
			if jsonBytes, err := json.Marshal(crops); err == nil {
				_ = app.Cache.Set(ctx, cacheKey, jsonBytes, 2*time.Hour)
			}

			utils.RespondWithJSON(w, http.StatusOK, utils.M{
				"success": true,
				"crops":   crops,
			})
			return
		}

		/* ------------------------------------------------ */
		/* 3. CSV fallback                                  */
		/* ------------------------------------------------ */

		file, err := os.Open("data/pre_crop_catalogue.csv")
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to retrieve catalogue",
			})
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Invalid CSV",
			})
			return
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil || len(record) != len(headers) {
				continue
			}

			item := farms.CropCatalogueItem{}
			for i, h := range headers {
				switch strings.ToLower(h) {
				case "name":
					item.Name = record[i]
				case "category":
					item.Category = record[i]
				case "banner":
					item.Banner = record[i]
				case "stock":
					item.Stock, _ = strconv.Atoi(record[i])
				case "unit":
					item.Unit = record[i]
				case "featured":
					item.Featured = strings.ToLower(record[i]) == "true"
				case "pricerange":
					parts := strings.Split(record[i], "-")
					if len(parts) == 2 {
						min, _ := strconv.Atoi(parts[0])
						max, _ := strconv.Atoi(parts[1])
						item.PriceRange = []int{min, max}
					}
				}
			}

			crops = append(crops, item)
		}

		if jsonBytes, err := json.Marshal(crops); err == nil {
			_ = app.Cache.Set(ctx, cacheKey, jsonBytes, 2*time.Hour)
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"crops":   crops,
		})
	}
}

/* ---------------------------------------------------- */
/* Crop Catalogue (unique by name + catalogueId)        */
/* ---------------------------------------------------- */

func GetCropCatalogue(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		allCrops, err := getAllCrops(ctx, app.DB)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to fetch crop catalogue",
			})
			return
		}

		seen := make(map[string]bool)
		unique := make([]farms.Crop, 0, len(allCrops))

		for _, c := range allCrops {
			key := strings.ToLower(c.Name + c.CatalogueId)
			if !seen[key] {
				seen[key] = true
				unique = append(unique, c)
			}
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"crops":   unique,
		})
	}
}

type MarketHighlightItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func GetMarketHighlights(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fallback := map[string][]MarketHighlightItem{
			"activeDeals": {
				{ID: "deal-1", Text: "Fresh harvest offers are arriving this week from nearby farms."},
				{ID: "deal-2", Text: "Weekend market bundles on leafy greens and tomatoes are in stock."},
			},
			"seasonalPicks": {
				{ID: "season-1", Text: "Seasonal favourites are at their peak quality and freshness."},
				{ID: "season-2", Text: "Cool-weather staples are selling fast across local mandis."},
			},
			"cropTrends": {
				{ID: "trend-1", Text: "Demand is trending up for high-value, low-water crops."},
				{ID: "trend-2", Text: "Farmers are pacing harvests around buyer demand and weather."},
			},
			"announcements": {
				{ID: "announce-1", Text: "New supply loads are hitting the market every morning."},
				{ID: "announce-2", Text: "Farm-to-door delivery slots are open for premium produce."},
			},
			"showcase": {
				{ID: "showcase-1", Text: "Farmer spotlight: fresh, traceable produce from trusted growers."},
				{ID: "showcase-2", Text: "Community harvest stories are now live for the latest market arrivals."},
			},
		}

		if app == nil || app.DB == nil {
			utils.RespondWithJSON(w, http.StatusOK, fallback)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		crops, err := getAllCrops(ctx, app.DB)
		if err != nil || len(crops) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, fallback)
			return
		}

		seen := map[string]bool{}
		cropNames := make([]string, 0, len(crops))
		for _, crop := range crops {
			name := strings.TrimSpace(crop.Name)
			if name == "" {
				continue
			}
			k := strings.ToLower(name)
			if !seen[k] {
				seen[k] = true
				cropNames = append(cropNames, name)
			}
		}
		if len(cropNames) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, fallback)
			return
		}

		sequence := append([]string{}, cropNames...)
		if len(sequence) > 5 {
			sequence = sequence[:5]
		}

		payload := map[string][]MarketHighlightItem{
			"activeDeals": {
				{ID: "deal-1", Text: "Fresh deals on " + sequence[0] + " are flowing in this week."},
				{ID: "deal-2", Text: "Local buyers are picking up " + sequence[min(1, len(sequence)-1)] + " at competitive prices."},
			},
			"seasonalPicks": {
				{ID: "season-1", Text: "Seasonal pick: " + sequence[0] + " is at peak quality."},
				{ID: "season-2", Text: "Best of the season: " + sequence[min(2, len(sequence)-1)] + " is stocked and fresh."},
			},
			"cropTrends": {
				{ID: "trend-1", Text: "Supply is rising for " + sequence[0] + " across nearby farms."},
				{ID: "trend-2", Text: "Market demand remains strong for " + sequence[min(3, len(sequence)-1)] + "."},
			},
			"announcements": {
				{ID: "announce-1", Text: "New harvest arrivals are expected for " + sequence[0] + "."},
				{ID: "announce-2", Text: "Fresh collections of " + sequence[min(4, len(sequence)-1)] + " are now available."},
			},
			"showcase": {
				{ID: "showcase-1", Text: "Farmer showcase: " + sequence[0] + " from a trusted local grower."},
				{ID: "showcase-2", Text: "Spotlight crop: " + sequence[min(1, len(sequence)-1)] + " bringing in premium quality."},
			},
		}

		utils.RespondWithJSON(w, http.StatusOK, payload)
	}
}

/* ---------------------------------------------------- */
/* Crop Types (grouped in application layer)             */
/* ---------------------------------------------------- */

func GetCropTypes(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		crops, err := getAllCrops(ctx, app.DB)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to fetch crops",
			})
			return
		}

		m := map[string]*Agg{}

		for _, c := range crops {
			a, ok := m[c.Name]
			if !ok {
				m[c.Name] = &Agg{
					Name:     c.Name,
					MinPrice: c.Price,
					MaxPrice: c.Price,
					Unit:     c.Unit,
					Banner:   getCropBanner(strings.ToLower(c.Name)),
				}
				a = m[c.Name]
			}

			if c.Price < a.MinPrice {
				a.MinPrice = c.Price
			}
			if c.Price > a.MaxPrice {
				a.MaxPrice = c.Price
			}
			if c.Quantity > 0 {
				a.AvailableCount++
			}
		}

		result := make([]Agg, 0, len(m))
		for _, v := range m {
			result = append(result, *v)
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success":   true,
			"cropTypes": result,
		})
	}
}
