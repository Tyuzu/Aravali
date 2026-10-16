// File: internal/beats/autocomplete/autocompleteSQLDB.go

package autocomplete

import (
	"context"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/places"
)

var (
	AutocompleteTable = config.Tables.AutocompleteTable
)

func findPlacesByQuery(ctx context.Context, app *infra.Deps, query string, results *[]places.Place) error {
	if app == nil || app.SQLDB == nil || results == nil {
		return nil
	}
	q := strings.TrimSpace(query)
	if q == "" {
		*results = []places.Place{}
		return nil
	}
	search := "%" + q + "%"
	rows, err := app.SQLDB.Query(ctx, `SELECT placeid, name, short_desc, description, place, capacity, date, address, createdby, organizer_name, organizer_contact, category, banner, website_url, status, accessibility_info, social_links, tags, custom_fields, created_at, updated_at, city, country, zipcode, jobs, location, phone, website, isopen, distance, views, reviewcount, sociallinks, updatedby, deletedat, amenities, events, operatinghours, keywords FROM places WHERE LOWER(name) LIKE LOWER($1) OR LOWER(category) LIKE LOWER($1) OR LOWER(address) LIKE LOWER($1) ORDER BY name LIMIT 10`, search)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]places.Place, 0, 10)
	for rows.Next() {
		var item places.Place
		if err := rows.Scan(&item.PlaceID, &item.Name, &item.ShortDesc, &item.Description, &item.Place, &item.Capacity, &item.Date, &item.Address, &item.CreatedBy, &item.OrganizerName, &item.OrganizerContact, &item.Category, &item.Banner, &item.WebsiteURL, &item.Status, &item.AccessibilityInfo, &item.SocialMediaLinks, &item.Tags, &item.CustomFields, &item.CreatedAt, &item.UpdatedAt, &item.City, &item.Country, &item.ZipCode, &item.Jobs, &item.Location, &item.Phone, &item.Website, &item.IsOpen, &item.Distance, &item.Views, &item.ReviewCount, &item.SocialLinks, &item.UpdatedBy, &item.DeletedAt, &item.Amenities, &item.Events, &item.OperatingHours, &item.Keywords); err != nil {
			return err
		}
		items = append(items, item)
	}
	*results = items
	return rows.Err()
}

func findUsersByQuery(ctx context.Context, app *infra.Deps, query string, users *[]auth.User) error {
	if app == nil || app.SQLDB == nil || users == nil {
		return nil
	}
	q := strings.TrimSpace(query)
	if q == "" {
		*users = []auth.User{}
		return nil
	}
	search := "%" + q + "%"
	rows, err := app.SQLDB.Query(ctx, `SELECT userid, username, email, name, bio, avatar, banner, created_at, updated_at, online, last_login, profile_views, phone_number, address, social_links, is_verified, email_verified, followerscount, followscount, wallet_balance FROM users WHERE LOWER(username) LIKE LOWER($1) OR LOWER(name) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1) ORDER BY username LIMIT 10`, search)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]auth.User, 0, 10)
	for rows.Next() {
		var item auth.User
		if err := rows.Scan(&item.UserID, &item.Username, &item.Email, &item.Name, &item.Bio, &item.Avatar, &item.Banner, &item.CreatedAt, &item.UpdatedAt, &item.Online, &item.LastLogin, &item.ProfileViews, &item.PhoneNumber, &item.Address, &item.SocialLinks, &item.IsVerified, &item.EmailVerified, &item.FollowersCount, &item.FollowingCount, &item.WalletBalance); err != nil {
			return err
		}
		items = append(items, item)
	}
	*users = items
	return rows.Err()
}
