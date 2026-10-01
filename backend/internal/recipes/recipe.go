// File: internal/recipes/recipe.go

package recipes

type IngredientAlternative struct {
	Name   string `json:"name" db:"name"`
	ItemID string `json:"itemId" db:"itemId"`
	Type   string `json:"type" db:"type"`
}

type Ingredient struct {
	Name         string                  `json:"name" db:"name"`
	ItemID       string                  `json:"itemId" db:"itemId"`
	Type         string                  `json:"type" db:"type"`
	Quantity     float64                 `json:"quantity" db:"quantity"`
	Unit         string                  `json:"unit" db:"unit"`
	Alternatives []IngredientAlternative `json:"alternatives" db:"alternatives"`
}

type Recipe struct {
	RecipeId    string       `db:"recipeid,omitempty" json:"recipeid"`
	UserID      string       `json:"userid" db:"userid"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	CookTime    string       `json:"cookTime" db:"cookTime"`       // replaced PrepTime
	Cuisine     string       `json:"cuisine" db:"cuisine"`         // new
	Dietary     []string     `json:"dietary" db:"dietary"`         // new
	PortionSize string       `json:"portionSize" db:"portionSize"` // new
	Season      string       `json:"season" db:"season"`           // new
	Tags        []string     `json:"tags" db:"tags"`
	Images      []string     `json:"images" db:"images"`
	Ingredients []Ingredient `json:"ingredients" db:"ingredients"`
	Steps       []string     `json:"steps" db:"steps"`
	Difficulty  string       `json:"difficulty" db:"difficulty"`
	Banner      string       `json:"banner" db:"banner"`
	Servings    int          `json:"servings" db:"servings"`
	VideoURL    string       `json:"videoUrl" db:"videoUrl"` // new
	Notes       string       `json:"notes" db:"notes"`       // new
	CreatedAt   int64        `json:"createdAt" db:"createdAt"`
	Views       int          `json:"views" db:"views"`
}
