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
	EventsCollection       = config.Collections.EventsCollection
	ArtistsCollection      = config.Collections.ArtistsCollection
	ArtistEventsCollection = config.Collections.ArtistEventsCollection
	ArtistAlbumsCollection = config.Collections.ArtistAlbumsCollection
	SubscribersCollection  = config.Collections.SubscribersCollection
)

func SQLInsertArtist(ctx context.Context, app *infra.Deps, artist *Artist) error {
	return app.SQLDB.Insert(ctx, ArtistsCollection, artist)
}

func SQLFindArtistByID(ctx context.Context, app *infra.Deps, artistID string, artist *Artist) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindOne(ctx, ArtistsCollection, query, args, artist)
}

func SQLUpdateArtistByID(ctx context.Context, app *infra.Deps, artistID string, update map[string]any) (int64, error) {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.Update(ctx, ArtistsCollection, query, args, update)
}

func SQLFindArtistEvents(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistEvent) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindMany(ctx, ArtistEventsCollection, query, args, result)
}

func SQLFindArtistAlbumsByArtistID(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistAlbum) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.FindMany(ctx, ArtistAlbumsCollection, query, args, result)
}

func SQLDeleteArtistRecordByID(ctx context.Context, app *infra.Deps, artistID string) (int64, error) {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.DeleteOne(ctx, ArtistsCollection, query, args)
}

// FindSubscribersForArtist checks if a specific user is subscribed to an artist.
func SQLFindSubscribersForArtist(ctx context.Context, app *infra.Deps, userID, artistID string) (bool, error) {
	var results []map[string]any
	query := "userid = $1 AND $2 = ANY(subscribed)"
	args := []any{userID, artistID}

	err := app.SQLDB.FindMany(ctx, SubscribersCollection, query, args, &results)
	if err != nil {
		return false, err
	}

	return len(results) > 0, nil
}

func SQLFindArtistsByEventID(ctx context.Context, app *infra.Deps, eventID string, result *[]Artist) error {
	query := "$1 = ANY(events)"
	args := []any{eventID}
	return app.SQLDB.FindMany(ctx, ArtistsCollection, query, args, result)
}

func SQLFindAllArtists(ctx context.Context, app *infra.Deps, result *[]Artist) error {
	return app.SQLDB.FindMany(ctx, ArtistsCollection, "", nil, result)
}

func SQLAddArtistMemberDB(ctx context.Context, app *infra.Deps, artistID string, member BandMember) error {
	query := "artistid = $1"
	args := []any{artistID}
	return app.SQLDB.AddToSet(ctx, ArtistsCollection, query, args, "members", member)
}

func SQLUpdateArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string, update map[string]any) (int64, error) {
	query := "artistid = $1 AND memberid = $2"
	args := []any{artistID, memberID}
	return app.SQLDB.Update(ctx, ArtistsCollection, query, args, update)
}

func SQLDeleteArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string) (int64, error) {
	query := "artistid = $1 AND memberid = $2"
	args := []any{artistID, memberID}
	return app.SQLDB.Delete(ctx, ArtistsCollection, query, args)
}

func SQLInsertArtistEvent(ctx context.Context, app *infra.Deps, artistevent *ArtistEvent) error {
	return app.SQLDB.Insert(ctx, ArtistEventsCollection, artistevent)
}

func SQLUpdateArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string, update map[string]any) (int64, error) {
	query := "eventid = $1"
	args := []any{artisteventID}
	return app.SQLDB.Update(ctx, ArtistEventsCollection, query, args, update)
}

// DeleteArtistEventByID deletes an artist event entry by its ID.
func SQLDeleteArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string) error {
	query := "eventid = $1"
	args := []any{artisteventID}
	_, err := app.SQLDB.DeleteOne(ctx, ArtistEventsCollection, query, args)
	return err
}

func SQLFindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *events.Event) error {
	query := "eventid = $1"
	args := []any{eventID}
	return app.SQLDB.FindOne(ctx, EventsCollection, query, args, event)
}

func SQLFindArtistEventsByEventAndArtist(ctx context.Context, app *infra.Deps, eventID, artistID string, result *[]ArtistEvent) error {
	query := "eventid = $1 AND artistid = $2"
	args := []any{eventID, artistID}
	return app.SQLDB.FindMany(ctx, ArtistEventsCollection, query, args, result)
}

func SQLAddArtistToEventDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) (int64, error) {
	if err := app.SQLDB.Insert(ctx, ArtistEventsCollection, artistEvent); err != nil {
		return 0, err
	}

	query := "eventid = $1"
	args := []any{artistEvent.EventID}
	err := app.SQLDB.AddToSet(ctx, EventsCollection, query, args, "artists", artistEvent.ArtistID)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func SQLAddEventToDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) error {
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

	if err := app.SQLDB.Insert(ctx, EventsCollection, event); err != nil {
		return err
	}

	userdata.SetUserData("event", event.EventID, artistEvent.ArtistID, "", "", app)
	return nil
}
