// File: internal/crops/cropAboutModels.go

package crops

type CropAbout struct {
	ID                 string             `db:"id" json:"id"`
	CommonName         string             `db:"commonName" json:"commonName"`
	ScientificName     string             `db:"scientificName" json:"scientificName"`
	Image              string             `db:"image" json:"image"`
	ImageAlt           string             `db:"imageAlt" json:"imageAlt"`
	Description        string             `db:"description" json:"description"`
	NutritionalValues  []NutritionalValue `db:"nutritionalValues" json:"nutritionalValues"`
	GrowingConditions  GrowingConditions  `db:"growingConditions" json:"growingConditions"`
	PlantingHarvesting string             `db:"plantingHarvesting" json:"plantingHarvesting"`
	CareTips           []string           `db:"careTips" json:"careTips"`
	Varieties          []string           `db:"varieties" json:"varieties"`
	Usage              string             `db:"usage" json:"usage"`
	FunFacts           []string           `db:"funFacts" json:"funFacts"`
}

type NutritionalValue struct {
	Label string `db:"label" json:"label"`
	Value string `db:"value" json:"value"`
}

type GrowingConditions struct {
	Soil        string `db:"soil" json:"soil"`
	Sunlight    string `db:"sunlight" json:"sunlight"`
	Water       string `db:"water" json:"water"`
	Temperature string `db:"temperature" json:"temperature"`
}
