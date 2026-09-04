package crops

import (
	"context"
	"scav/infra"
)

func CreateCropAbout(
	ctx context.Context,
	app *infra.Deps,
	crop *CropAbout,
) error {
	return app.DB.InsertOne(
		ctx,
		cropsAboutCollection,
		crop,
	)
}

func GetCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) (*CropAbout, error) {

	var crop CropAbout

	err := app.DB.FindOne(
		ctx,
		cropsAboutCollection,
		map[string]any{"id": cropID},
		&crop,
	)

	if err != nil {
		return nil, err
	}

	return &crop, nil
}

func GetAllCropAbouts(
	ctx context.Context,
	app *infra.Deps,
) ([]CropAbout, error) {

	var crops []CropAbout

	err := app.DB.FindMany(
		ctx,
		cropsAboutCollection,
		map[string]any{},
		&crops,
	)

	return crops, err
}

func UpdateCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
	crop *CropAbout,
) (any, error) {

	return app.DB.UpdateOne(
		ctx,
		cropsAboutCollection,
		map[string]any{"id": cropID},
		map[string]any{
			"$set": crop,
		},
	)
}

func DeleteCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) error {

	_, err := app.DB.DeleteOne(
		ctx,
		cropsAboutCollection,
		map[string]any{"id": cropID},
	)

	return err
}
