// File: internal/artists/artistSQLDB.go

package artists

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

func buildSetClause(update map[string]any) (string, []any) {
	if len(update) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	idx := 1
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, idx))
		args = append(args, value)
		idx++
	}
	return strings.Join(parts, ", "), args
}

func InsertArtist(ctx context.Context, app *infra.Deps, artist *Artist) error {
	if app == nil || app.SQLDB == nil || artist == nil {
		return nil
	}
	if artist.CreatedAt.IsZero() {
		artist.CreatedAt = time.Now().UTC()
	}
	genresJSON, _ := json.Marshal(artist.Genres)
	socialsJSON, _ := json.Marshal(artist.Socials)
	_, err := app.SQLDB.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (artistid, category, name, place, country, bio, dob, photo, banner, genres, socials, createdat, creatorid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", ArtistsTable),
		artist.ArtistID, artist.Category, artist.Name, artist.Place, artist.Country, artist.Bio, artist.DOB, artist.Photo, artist.Banner, string(genresJSON), string(socialsJSON), artist.CreatedAt, artist.CreatorID,
	)
	return err
}

func FindArtistByID(ctx context.Context, app *infra.Deps, artistID string, artist *Artist) error {
	if app == nil || app.SQLDB == nil || artist == nil {
		return nil
	}
	if artistID == "" {
		*artist = Artist{}
		return nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT artistid, category, name, place, country, bio, dob, photo, banner, genres, socials, createdat, creatorid FROM %s WHERE artistid = $1 LIMIT 1", ArtistsTable), artistID)
	var got Artist
	var genresRaw, socialsRaw any
	if err := row.Scan(&got.ArtistID, &got.Category, &got.Name, &got.Place, &got.Country, &got.Bio, &got.DOB, &got.Photo, &got.Banner, &genresRaw, &socialsRaw, &got.CreatedAt, &got.CreatorID); err != nil {
		*artist = Artist{}
		return nil
	}
	if genresRaw != nil {
		if b, ok := genresRaw.([]byte); ok {
			_ = json.Unmarshal(b, &got.Genres)
		}
	}
	if socialsRaw != nil {
		if b, ok := socialsRaw.([]byte); ok {
			_ = json.Unmarshal(b, &got.Socials)
		}
	}
	*artist = got
	return nil
}

func UpdateArtistByID(ctx context.Context, app *infra.Deps, artistID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistID == "" || len(update) == 0 {
		return 0, nil
	}
	setClause, args := buildSetClause(update)
	if setClause == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE artistid = $%d", ArtistsTable, setClause, len(args)+1), append(args, artistID)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func FindArtistEvents(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistEvent) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if artistID == "" {
		*result = nil
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE artistid = $1", ArtistEventsTable), artistID)
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func FindArtistAlbumsByArtistID(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistAlbum) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if artistID == "" {
		*result = nil
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE artistid = $1", ArtistAlbumsTable), artistID)
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func DeleteArtistRecordByID(ctx context.Context, app *infra.Deps, artistID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistID == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE artistid = $1", ArtistsTable), artistID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func FindSubscribersForArtist(ctx context.Context, app *infra.Deps, userID, artistID string) (bool, error) {
	if app == nil || app.SQLDB == nil {
		return false, nil
	}
	if userID == "" || artistID == "" {
		return false, nil
	}
	var count int
	err := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT 1 FROM %s WHERE userid = $1 AND artistid = $2 LIMIT 1", SubscribersTable), userID, artistID).Scan(&count)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func FindArtistsByEventID(ctx context.Context, app *infra.Deps, eventID string, result *[]Artist) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if eventID == "" {
		*result = nil
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT a.* FROM %s a JOIN %s e ON e.artistid = a.artistid WHERE e.eventid = $1", ArtistsTable, ArtistEventsTable), eventID)
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func FindAllArtists(ctx context.Context, app *infra.Deps, result *[]Artist) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s ORDER BY createdat DESC", ArtistsTable))
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func AddArtistMemberDB(ctx context.Context, app *infra.Deps, artistID string, member BandMember) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if artistID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("INSERT INTO %s (artistid, memberid, ref_artistid, name, role, dob, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", ArtistsTable), artistID, member.MemberID, member.ReferenceArtist, member.Name, member.Role, member.DOB, member.Image)
	return err
}

func UpdateArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistID == "" || memberID == "" || len(update) == 0 {
		return 0, nil
	}
	setClause, args := buildSetClause(update)
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE artistid = $%d AND memberid = $%d", ArtistsTable, setClause, len(args)+1, len(args)+2), append(args, artistID, memberID)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func DeleteArtistMemberDB(ctx context.Context, app *infra.Deps, artistID, memberID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistID == "" || memberID == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE artistid = $1 AND memberid = $2", ArtistsTable), artistID, memberID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func InsertArtistEvent(ctx context.Context, app *infra.Deps, artistevent *ArtistEvent) error {
	if app == nil || app.SQLDB == nil || artistevent == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (eventid, artistid, title, date, venue, city, country, creatorid, ticket_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", ArtistEventsTable),
		artistevent.EventID, artistevent.ArtistID, artistevent.Title, artistevent.Date, artistevent.Venue, artistevent.City, artistevent.Country, artistevent.CreatorID, artistevent.TicketURL,
	)
	return err
}

func UpdateArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artisteventID == "" || len(update) == 0 {
		return 0, nil
	}
	setClause, args := buildSetClause(update)
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE eventid = $%d", ArtistEventsTable, setClause, len(args)+1), append(args, artisteventID)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func DeleteArtistEventByID(ctx context.Context, app *infra.Deps, artisteventID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if artisteventID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE eventid = $1", ArtistEventsTable), artisteventID)
	return err
}

func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *events.Event) error {
	if app == nil || app.SQLDB == nil || event == nil {
		return nil
	}
	if eventID == "" {
		*event = events.Event{}
		return nil
	}
	var out events.Event
	if err := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 LIMIT 1", EventsTable), eventID).Scan(&out.EventID, &out.Title, &out.Status, &out.Location, &out.CreatedAt, &out.Date, &out.CreatorID); err != nil {
		*event = events.Event{}
		return nil
	}
	*event = out
	return nil
}

func FindArtistEventsByEventAndArtist(ctx context.Context, app *infra.Deps, eventID, artistID string, result *[]ArtistEvent) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if eventID == "" || artistID == "" {
		*result = nil
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 AND artistid = $2", ArtistEventsTable), eventID, artistID)
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func AddArtistToEventDB(ctx context.Context, app *infra.Deps, artistEvent ArtistEvent) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistEvent.EventID == "" || artistEvent.ArtistID == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("INSERT INTO %s (eventid, artistid, title, date, venue, city, country, creatorid, ticket_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (eventid, artistid) DO NOTHING", ArtistEventsTable), artistEvent.EventID, artistEvent.ArtistID, artistEvent.Title, artistEvent.Date, artistEvent.Venue, artistEvent.City, artistEvent.Country, artistEvent.CreatorID, artistEvent.TicketURL)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
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
