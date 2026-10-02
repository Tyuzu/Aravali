// File: internal/artists/artistSQLDB.go

package artists

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/beats/userdata"
	"scav/internal/events"
)

var (
	EventsTable       = config.Tables.EventsTable
	ArtistsTable      = config.Tables.ArtistsTable
	ArtistEventsTable = config.Tables.ArtistEventsTable
	ArtistAlbumsTable = config.Tables.ArtistAlbumsTable
	SubscribersTable  = config.Tables.SubscribersTable
)

func InsertArtist(ctx context.Context, app *infra.Deps, artist *Artist) error {
	return app.SQLDB.Insert(ctx, ArtistsTable, artist)
}

func FindArtistByID(ctx context.Context, app *infra.Deps, artistID string, artist *Artist) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindOne(ctx, ArtistsTable, query, args, artist)
}

func UpdateArtistByID(ctx context.Context, app *infra.Deps, artistID string, update map[string]any) (int64, error) {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.Update(ctx, ArtistsTable, query, args, update)
}

func FindArtistEvents(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistEvent) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindMany(ctx, ArtistEventsTable, query, args, result)
}

func FindArtistAlbumsByArtistID(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistAlbum) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindMany(ctx, ArtistAlbumsTable, query, args, result)
}

func DeleteArtistRecordByID(ctx context.Context, app *infra.Deps, artistID string) (int64, error) {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.DeleteOne(ctx, ArtistsTable, query, args)
}

// FindSubscribersForArtist checks if a specific user is subscribed to an artist.
func FindSubscribersForArtist(ctx context.Context, app *infra.Deps, userID, artistID string) (bool, error) {
	var results []map[string]any
	query := "userid = $1 AND $2 = ANY(subscribed)"
	args := []any{userID, artistID}

	err := app.SQLDB.FindMany(ctx, SubscribersTable, query, args, &results)
	if err != nil {
		return false, err
	}

	return len(results) > 0, nil
}

func FindArtistsByEventID(ctx context.Context, app *infra.Deps, eventID string, result *[]Artist) error {
	query := "$1 = ANY(events)"
	args := []any{eventID}
	return app.SQLDB.FindMany(ctx, ArtistsTable, query, args, result)
}

func FindAllArtists(ctx context.Context, app *infra.Deps, result *[]Artist) error {
	return app.SQLDB.FindMany(ctx, ArtistsTable, "", nil, result)
}

func AddArtistMemberDB(ctx context.Context, app *infra.Deps, artistID string, member BandMember) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.AddToSet(ctx, ArtistsTable, query, args, "members", member)
}

func UpdateArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string, update map[string]any) (int64, error) {
	query := "artistid = $1 AND memberid = $2"
	args := []any{artistID, memberID}
	return app.SQLDB.Update(ctx, ArtistsTable, query, args, update)
}

func DeleteArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string) (int64, error) {
	query := "artistid = $1 AND memberid = $2"
	args := []any{artistID, memberID}
	return app.SQLDB.Delete(ctx, ArtistsTable, query, args)
}

func InsertArtistEvent(ctx context.Context, app *infra.Deps, artistevent *ArtistEvent) error {
	return app.SQLDB.Insert(ctx, ArtistEventsTable, artistevent)
}

func UpdateArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string, update map[string]any) (int64, error) {
	query := "eventid = $1"
	args := []any{artisteventID}
	return app.SQLDB.Update(ctx, ArtistEventsTable, query, args, update)
}

// DeleteArtistEventByID deletes an artist event entry by its ID.
func DeleteArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string) error {
	query := "eventid = $1"
	args := []any{artisteventID}
	_, err := app.SQLDB.DeleteOne(ctx, ArtistEventsTable, query, args)
	return err
}

func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *events.Event) error {
	query := "eventid = $1"
	args := []any{eventID}
	return app.SQLDB.FindOne(ctx, EventsTable, query, args, event)
}

func FindArtistEventsByEventAndArtist(ctx context.Context, app *infra.Deps, eventID, artistID string, result *[]ArtistEvent) error {
	query := "eventid = $1 AND artistid = $2"
	args := []any{eventID, artistID}
	return app.SQLDB.FindMany(ctx, ArtistEventsTable, query, args, result)
}

func AddArtistToEventDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) (int64, error) {
	if err := app.SQLDB.Insert(ctx, ArtistEventsTable, artistEvent); err != nil {
		return 0, err
	}

	query := "eventid = $1"
	args := []any{artistEvent.EventID}
	err := app.SQLDB.AddToSet(ctx, EventsTable, query, args, "artists", artistEvent.ArtistID)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func AddEventToDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) error {
	var event events.Event
	dateString := artistEvent.Date
	layout := "2006-01-02"
	dateToSave, _ := time.Parse(layout, dateString)

	event.CreatorID = artistEvent.CreatorID
	event.CreatedAt = time.Now().UTC()
	event.Date = dateToSave.UTC()
	event.Status = "active"
	event.EventID = artistEvent.EventID
	event.Artists = []string{artistEvent.ArtistID}
	event.Title = artistEvent.Title
	event.Location = artistEvent.Venue
	event.Published = "draft"
	event.Category = "concert"

	if err := app.SQLDB.Insert(ctx, EventsTable, event); err != nil {
		return err
	}

	userdata.SetUserData("event", event.EventID, artistEvent.ArtistID, "", "", app)
	return nil
}
