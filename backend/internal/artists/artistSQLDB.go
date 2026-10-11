// File: internal/artists/artistSQLDB.go

package artists

import (
	"context"
	"database/sql"
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

func parseArtistMetadata(metadata any) Artist {
	artist := Artist{}
	if metadata == nil {
		return artist
	}

	var payload map[string]any
	switch v := metadata.(type) {
	case []byte:
		if len(v) == 0 {
			return artist
		}
		_ = json.Unmarshal(v, &payload)
	case string:
		if v == "" {
			return artist
		}
		_ = json.Unmarshal([]byte(v), &payload)
	case map[string]any:
		payload = v
	case map[string]string:
		payload = make(map[string]any, len(v))
		for k, val := range v {
			payload[k] = val
		}
	}
	if len(payload) == 0 {
		return artist
	}

	artist.Category = stringFromAny(payload["category"])
	artist.Name = stringFromAny(payload["name"])
	artist.Place = stringFromAny(payload["place"])
	artist.Country = stringFromAny(payload["country"])
	artist.Bio = stringFromAny(payload["bio"])
	artist.DOB = stringFromAny(payload["dob"])
	artist.Photo = stringFromAny(payload["photo"])
	artist.Banner = stringFromAny(payload["banner"])
	artist.CreatorID = stringFromAny(payload["creatorid"])
	artist.Genres = stringSliceFromAny(payload["genres"])
	artist.Socials = stringMapFromAny(payload["socials"])

	if createdAt := stringFromAny(payload["createdAt"]); createdAt != "" {
		if ts, err := time.Parse(time.RFC3339, createdAt); err == nil {
			artist.CreatedAt = ts
		}
	}
	if artist.CreatedAt.IsZero() {
		if createdAt := stringFromAny(payload["created_at"]); createdAt != "" {
			if ts, err := time.Parse(time.RFC3339, createdAt); err == nil {
				artist.CreatedAt = ts
			}
		}
	}
	return artist
}

func stringFromAny(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case []byte:
		return string(x)
	default:
		return ""
	}
}

func stringSliceFromAny(v any) []string {
	switch x := v.(type) {
	case []string:
		return append([]string(nil), x...)
	case []any:
		out := make([]string, 0, len(x))
		for _, item := range x {
			if s := stringFromAny(item); s != "" {
				out = append(out, s)
			}
		}
		return out
	case string:
		if x == "" {
			return nil
		}
		return []string{x}
	default:
		return nil
	}
}

func stringMapFromAny(v any) map[string]string {
	out := make(map[string]string)
	switch x := v.(type) {
	case map[string]string:
		for k, val := range x {
			out[k] = val
		}
	case map[string]any:
		for k, val := range x {
			if s := stringFromAny(val); s != "" {
				out[k] = s
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func FindArtistByID(ctx context.Context, app *infra.Deps, artistID string, artist *Artist) error {
	if app == nil || app.SQLDB == nil || artist == nil {
		return nil
	}
	if artistID == "" {
		*artist = Artist{}
		return sql.ErrNoRows
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT artistid, name, bio, userids, status, created_at, updated_at, metadata FROM %s WHERE artistid = $1 LIMIT 1", ArtistsTable), artistID)
	var (
		artistIDDB, name, bio, status sql.NullString
		userIDs                       []string
		createdAt, updatedAt          time.Time
		metadata                      []byte
	)
	if err := row.Scan(&artistIDDB, &name, &bio, &userIDs, &status, &createdAt, &updatedAt, &metadata); err != nil {
		*artist = Artist{}
		return err
	}
	got := parseArtistMetadata(metadata)
	got.ArtistID = artistIDDB.String
	if got.Name == "" {
		got.Name = name.String
	}
	if got.Bio == "" {
		got.Bio = bio.String
	}
	if got.CreatedAt.IsZero() {
		got.CreatedAt = createdAt
	}
	if got.CreatedAt.IsZero() {
		got.CreatedAt = updatedAt
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
	if _, err := app.SQLDB.Exec(ctx, fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", ArtistEventsTable)); err != nil {
		*result = []Artist{}
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
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT artistid, name, bio, userids, status, created_at, updated_at, metadata FROM %s ORDER BY created_at DESC", ArtistsTable))
	if err != nil {
		return err
	}
	defer rows.Close()

	artists := make([]Artist, 0)
	for rows.Next() {
		var (
			artistID, name, bio, status sql.NullString
			userIDs                     []string
			createdAt, updatedAt        time.Time
			metadata                    []byte
		)
		if err := rows.Scan(&artistID, &name, &bio, &userIDs, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		artist := parseArtistMetadata(metadata)
		artist.ArtistID = artistID.String
		if artist.Name == "" {
			artist.Name = name.String
		}
		if artist.Bio == "" {
			artist.Bio = bio.String
		}
		if artist.CreatedAt.IsZero() {
			artist.CreatedAt = createdAt
		}
		if artist.CreatedAt.IsZero() {
			artist.CreatedAt = updatedAt
		}
		artists = append(artists, artist)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*result = artists
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
	if app == nil {
		return nil
	}
	if artistEvent.EventID == "" || artistEvent.ArtistID == "" {
		return nil
	}

	userdata.SetUserData("event", artistEvent.EventID, artistEvent.ArtistID, "", "", app)
	return nil
}
