package artists

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
	db "scav/infra/db"
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

func InsertArtist(ctx context.Context, db db.Database, artist *Artist) error {
	return db.Insert(ctx, ArtistsTable, artist)
}

func FindArtistByID(ctx context.Context, db db.Database, artistID string, artist *Artist) error {
	return db.FindOne(ctx, ArtistsTable, map[string]any{"artistid": artistID}, artist)
}

func UpdateArtistByID(ctx context.Context, db db.Database, artistID string, update map[string]any) (any, error) {
	return db.Update(ctx, ArtistsTable, map[string]any{"artistid": artistID}, map[string]any{"$set": update})
}

func FindArtistEvents(ctx context.Context, db db.Database, artistID string, result *[]ArtistEvent) error {
	return db.FindMany(ctx, ArtistEventsTable, map[string]any{"artistid": artistID}, result)
}

func FindArtistAlbumsByArtistID(ctx context.Context, db db.Database, artistID string, result *[]ArtistAlbum) error {
	return db.FindMany(ctx, ArtistAlbumsTable, map[string]any{"artistid": artistID}, result)
}

func DeleteArtistRecordByID(ctx context.Context, db db.Database, artistID string) (any, error) {
	return db.DeleteOne(ctx, ArtistsTable, map[string]any{"artistid": artistID})
}

// FindSubscribersForArtist checks if a specific user is subscribed to an artist.
func FindSubscribersForArtist(ctx context.Context, db db.Database, userID, artistID string) (bool, error) {
	var results []map[string]any
	err := db.FindMany(ctx, SubscribersTable, map[string]any{
		"userid": userID,
		"subscribed": map[string]any{
			"$in": []string{artistID},
		},
	}, &results)

	if err != nil {
		return false, err
	}

	return len(results) > 0, nil
}

func FindArtistsByEventID(ctx context.Context, db db.Database, eventID string, result *[]Artist) error {
	return db.FindMany(ctx, ArtistsTable, map[string]any{"events": eventID}, result)
}

func FindAllArtists(ctx context.Context, db db.Database, result *[]Artist) error {
	return db.FindMany(ctx, ArtistsTable, map[string]any{}, result)
}

func AddArtistMemberDB(ctx context.Context, db db.Database, artistID string, member BandMember) (any, error) {
	return db.Update(ctx, ArtistsTable, map[string]any{"artistid": artistID}, map[string]any{"$push": map[string]any{"members": member}})
}

func UpdateArtistMemberDB(ctx context.Context, db db.Database, artistID, memberID string, update map[string]any) (any, error) {
	return db.Update(ctx, ArtistsTable, map[string]any{"artistid": artistID, "members.memberid": memberID}, map[string]any{"$set": update})
}

func DeleteArtistMemberDB(ctx context.Context, db db.Database, artistID, memberID string) (any, error) {
	return db.Update(ctx, ArtistsTable, map[string]any{"artistid": artistID}, map[string]any{"$pull": map[string]any{"members": map[string]any{"memberid": memberID}}})
}

func InsertArtistEvent(ctx context.Context, db db.Database, artistevent *ArtistEvent) error {
	return db.Insert(ctx, ArtistEventsTable, artistevent)
}

func UpdateArtistEventByID(ctx context.Context, db db.Database, artisteventID string, update map[string]any) (any, error) {
	return db.Update(ctx, ArtistEventsTable, map[string]any{"eventid": artisteventID}, update)
}

// DeleteArtistEventByID deletes an artist event entry by its ID.
func DeleteArtistEventByID(ctx context.Context, db db.Database, artisteventID string) error {
	_, err := db.DeleteOne(ctx, ArtistEventsTable, map[string]any{"eventid": artisteventID})
	return err
}

func FindEventByID(ctx context.Context, db db.Database, eventID string, event *events.Event) error {
	return db.FindOne(ctx, EventsTable, map[string]any{"eventid": eventID}, event)
}

func FindArtistEventsByEventAndArtist(ctx context.Context, db db.Database, eventID, artistID string, result *[]ArtistEvent) error {
	return db.FindMany(ctx, ArtistEventsTable, map[string]any{"eventid": eventID, "artistid": artistID}, result)
}

func AddArtistToEventDB(ctx context.Context, db db.Database, artistEvent ArtistEvent) (any, error) {
	if err := db.Insert(ctx, ArtistEventsTable, artistEvent); err != nil {
		return nil, err
	}
	return db.Update(ctx, EventsTable, map[string]any{"eventid": artistEvent.EventID}, map[string]any{"$addToSet": map[string]any{"artists": artistEvent.ArtistID}})
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

	if err := app.DB.Insert(ctx, EventsTable, event); err != nil {
		return err
	}

	userdata.SetUserData("event", event.EventID, artistEvent.ArtistID, "", "", app)
	return nil
}
