// File: internal/crops/cropaboutrepo.go

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
	return createCropAbout(ctx, app.SQLDB, crop)
}

func GetCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) (*CropAbout, error) {
	return getCropAboutByID(ctx, app.SQLDB, cropID)
}

func GetAllCropAbouts(
	ctx context.Context,
	app *infra.Deps,
) ([]CropAbout, error) {
	return getAllCropAbouts(ctx, app.SQLDB)
}

func UpdateCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
	crop *CropAbout,
) (any, error) {
	return updateCropAbout(ctx, app.SQLDB, cropID, crop)
}

func DeleteCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) error {
	return deleteCropAbout(ctx, app.SQLDB, cropID)
}
