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
	return QQQcreateCropAbout(ctx, app, crop)
}

func GetCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) (*CropAbout, error) {
	return QQQgetCropAboutByID(ctx, app, cropID)
}

func GetAllCropAbouts(
	ctx context.Context,
	app *infra.Deps,
) ([]CropAbout, error) {
	return QQQgetAllCropAbouts(ctx, app)
}

func UpdateCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
	crop *CropAbout,
) (any, error) {
	return QQQupdateCropAbout(ctx, app, cropID, crop)
}

func DeleteCropAbout(
	ctx context.Context,
	app *infra.Deps,
	cropID string,
) error {
	return QQQdeleteCropAbout(ctx, app, cropID)
}
