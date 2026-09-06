package placedb

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var placesCollection = config.Collections.PlacesCollection
var eventsCollection = config.Collections.EventsCollection
var productsCollection = config.Collections.ProductCollection

// Places
func FindPlaces(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, placesCollection, filter, out)
}

func FindOnePlace(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindOne(ctx, placesCollection, filter, out)
}

func UpdatePlace(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.Update(ctx, placesCollection, filter, update)
}

func InsertPlace(ctx context.Context, app *infra.Deps, place any) error {
	return app.DB.Insert(ctx, placesCollection, place)
}

func DeletePlace(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, placesCollection, filter)
}

// Events
func CountEvents(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.Count(ctx, eventsCollection, filter)
}

func FindEventsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out any) error {
	return app.DB.FindManyWithOptions(ctx, eventsCollection, filter, opts, out)
}

// Place products (generic)
func FindPlaceProducts(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, productsCollection, filter, out)
}

func InsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
	return app.DB.InsertOne(ctx, productsCollection, product)
}

func UpdatePlaceProduct(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, productsCollection, filter, update)
}

func DeletePlaceProduct(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, productsCollection, filter)
}
