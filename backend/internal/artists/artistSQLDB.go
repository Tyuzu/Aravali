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
}

func FindArtistByID(ctx context.Context, app *infra.Deps, artistID string, artist *Artist) error {
}

func UpdateArtistByID(ctx context.Context, app *infra.Deps, artistID string, update map[string]any) (int64, error) {

}

func FindArtistEvents(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistEvent) error {

}

func FindArtistAlbumsByArtistID(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistAlbum) error {

}

func DeleteArtistRecordByID(ctx context.Context, app *infra.Deps, artistID string) (int64, error) {

}

// FindSubscribersForArtist checks if a specific user is subscribed to an artist.
func FindSubscribersForArtist(ctx context.Context, app *infra.Deps, userID, artistID string) (bool, error) {
	var results []map[string]any

	return len(results) > 0, nil
}

func FindArtistsByEventID(ctx context.Context, app *infra.Deps, eventID string, result *[]Artist) error {

}

func FindAllArtists(ctx context.Context, app *infra.Deps, result *[]Artist) error {

}

func AddArtistMemberDB(ctx context.Context, app *infra.Deps, artistID string, member BandMember) error {

}

func UpdateArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string, update map[string]any) (int64, error) {

}

func DeleteArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string) (int64, error) {

}

func InsertArtistEvent(ctx context.Context, app *infra.Deps, artistevent *ArtistEvent) error {

}

func UpdateArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string, update map[string]any) (int64, error) {

}

// DeleteArtistEventByID deletes an artist event entry by its ID.
func DeleteArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string) error {

}

func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *events.Event) error {

}

func FindArtistEventsByEventAndArtist(ctx context.Context, app *infra.Deps, eventID, artistID string, result *[]ArtistEvent) error {

}

func AddArtistToEventDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) (int64, error) {

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

	userdata.SetUserData("event", event.EventID, artistEvent.ArtistID, "", "", app)
	return nil
}
