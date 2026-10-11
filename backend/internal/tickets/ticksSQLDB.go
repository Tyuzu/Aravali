// File: internal/tickets/ticksSQLDB.go

package tickets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"

	"github.com/jackc/pgx/v5"
)

var ticketsTable = config.Tables.TicketsTable
var bookingsTable = config.Tables.BookingsTable
var purchasedTicketsTable = config.Tables.PurchasedTicketsTable
var refundsTable = config.Tables.RefundsTable

func ensureTicketDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func normalizeTicketWhere(query string) string {
	if strings.TrimSpace(query) == "" {
		return "1 = 1"
	}
	q := query
	q = strings.ReplaceAll(q, "ticketid", "tickid")
	q = strings.ReplaceAll(q, "purchasedticketid", "purchasedticketid")
	q = strings.ReplaceAll(q, "uniquecode", "metadata->>'uniqueCode'")
	q = strings.ReplaceAll(q, "buyername", "metadata->>'buyerName'")
	q = strings.ReplaceAll(q, "cancelledreason", "metadata->>'canceledReason'")
	q = strings.ReplaceAll(q, "canceledreason", "metadata->>'canceledReason'")
	return q
}

func normalizeRefundWhere(query string) string {
	if strings.TrimSpace(query) == "" {
		return "1 = 1"
	}
	q := query
	q = strings.ReplaceAll(q, "uniquecode", "metadata->>'uniqueCode'")
	q = strings.ReplaceAll(q, "eventid", "metadata->>'eventID'")
	q = strings.ReplaceAll(q, "ticketid", "metadata->>'ticketID'")
	q = strings.ReplaceAll(q, "userid", "userid")
	return q
}

func parseJSONMap(data []byte) map[string]any {
	if len(data) == 0 {
		return map[string]any{}
	}
	payload := map[string]any{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return map[string]any{}
	}
	if payload == nil {
		payload = map[string]any{}
	}
	return payload
}

func coerceInt64(v any) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case string:
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return int64(parsed)
		}
	}
	return 0
}

func coerceInt(v any) int {
	return int(coerceInt64(v))
}

func parseNumericString(raw string) int64 {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return int64(value * 100)
}

func ticketFromRow(ticketID, eventID, title string, priceValue any, quantityValue any, createdAt, updatedAt time.Time, metadata []byte) (*Ticket, error) {
	t := &Ticket{
		TicketID:  ticketID,
		EventID:   eventID,
		Name:      title,
		Quantity:  coerceInt(quantityValue),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if value, ok := priceValue.(string); ok {
		t.Price = parseNumericString(value)
	} else {
		t.Price = coerceInt64(priceValue)
	}
	payload := parseJSONMap(metadata)
	if payload != nil {
		if v, ok := payload["ticketid"]; ok {
			if s, ok := v.(string); ok && s != "" {
				t.TicketID = s
			}
		}
		if v, ok := payload["eventid"]; ok {
			if s, ok := v.(string); ok && s != "" {
				t.EventID = s
			}
		}
		if v, ok := payload["name"]; ok {
			if s, ok := v.(string); ok && s != "" {
				t.Name = s
			}
		}
		if v, ok := payload["title"]; ok {
			if s, ok := v.(string); ok && s != "" {
				t.Name = s
			}
		}
		if v, ok := payload["currency"]; ok {
			if s, ok := v.(string); ok {
				t.Currency = s
			}
		}
		if v, ok := payload["color"]; ok {
			if s, ok := v.(string); ok {
				t.Color = s
			}
		}
		if v, ok := payload["description"]; ok {
			if s, ok := v.(string); ok {
				t.Description = s
			}
		}
		if v, ok := payload["entity_id"]; ok {
			if s, ok := v.(string); ok {
				t.EntityID = s
			}
		}
		if v, ok := payload["entity_type"]; ok {
			if s, ok := v.(string); ok {
				t.EntityType = s
			}
		}
		if v, ok := payload["quantity"]; ok {
			t.Quantity = coerceInt(v)
		}
		if v, ok := payload["available"]; ok {
			t.Available = coerceInt(v)
		} else {
			t.Available = t.Quantity
		}
		if v, ok := payload["total"]; ok {
			t.Total = coerceInt(v)
		} else {
			t.Total = t.Quantity
		}
		if v, ok := payload["sold"]; ok {
			t.Sold = coerceInt(v)
		}
		if v, ok := payload["seatstart"]; ok {
			t.SeatStart = coerceInt(v)
		}
		if v, ok := payload["seatend"]; ok {
			t.SeatEnd = coerceInt(v)
		}
		if v, ok := payload["price"]; ok {
			t.Price = coerceInt64(v)
		}
		if t.Price == 0 && title != "" {
			if pval, ok := payload["price"]; ok {
				t.Price = coerceInt64(pval)
			}
		}
		if v, ok := payload["seats"]; ok {
			if arr, ok := v.([]any); ok {
				seats := make([]Seat, 0, len(arr))
				for _, item := range arr {
					if m, ok := item.(map[string]any); ok {
						seat := Seat{}
						if s, ok := m["id"]; ok {
							seat.SeatID, _ = s.(string)
						}
						if s, ok := m["entity_id"]; ok {
							seat.EntityID, _ = s.(string)
						}
						if s, ok := m["entity_type"]; ok {
							seat.EntityType, _ = s.(string)
						}
						if s, ok := m["seat_number"]; ok {
							seat.SeatNumber, _ = s.(string)
						}
						if s, ok := m["userid"]; ok {
							seat.UserID, _ = s.(string)
						}
						if s, ok := m["status"]; ok {
							seat.Status, _ = s.(string)
						}
						seats = append(seats, seat)
					}
				}
				t.Seats = seats
			}
		}
	}
	if t.Name == "" {
		t.Name = title
	}
	if t.EventID == "" {
		t.EventID = eventID
	}
	if t.TicketID == "" {
		t.TicketID = ticketID
	}
	if t.Available == 0 {
		t.Available = t.Quantity
	}
	if t.Total == 0 {
		t.Total = t.Quantity
	}
	return t, nil
}

func purchasedTicketFromRow(purchasedTicketID, userID, eventID, ticketID, status string, amount any, createdAt, updatedAt time.Time, metadata []byte) (*PurchasedTicket, error) {
	p := &PurchasedTicket{
		EventID:      eventID,
		TicketID:     ticketID,
		UserID:       userID,
		Status:       status,
		PurchaseDate: createdAt,
		Price:        int(coerceInt64(amount)),
		Canceled:     strings.EqualFold(status, "canceled") || strings.EqualFold(status, "cancelled"),
		Transferred:  false,
	}
	payload := parseJSONMap(metadata)
	if len(payload) > 0 {
		if v, ok := payload["eventid"]; ok {
			if s, ok := v.(string); ok {
				p.EventID = s
			}
		}
		if v, ok := payload["userid"]; ok {
			if s, ok := v.(string); ok {
				p.UserID = s
			}
		}
		if v, ok := payload["ticketid"]; ok {
			if s, ok := v.(string); ok {
				p.TicketID = s
			}
		}
		if v, ok := payload["buyerName"]; ok {
			if s, ok := v.(string); ok {
				p.BuyerName = s
			}
		}
		if v, ok := payload["uniqueCode"]; ok {
			if s, ok := v.(string); ok {
				p.UniqueCode = s
			}
		}
		if v, ok := payload["status"]; ok {
			if s, ok := v.(string); ok {
				p.Status = s
			}
		}
		if v, ok := payload["price"]; ok {
			p.Price = coerceInt(v)
		}
		if v, ok := payload["canceled"]; ok {
			if b, ok := v.(bool); ok {
				p.Canceled = b
			}
		}
		if v, ok := payload["transferred"]; ok {
			if b, ok := v.(bool); ok {
				p.Transferred = b
			}
		}
		if v, ok := payload["transferredTo"]; ok {
			if s, ok := v.(string); ok {
				p.TransferredTo = s
			}
		}
		if v, ok := payload["canceledReason"]; ok {
			if s, ok := v.(string); ok {
				p.CanceledReason = s
			}
		}
		if v, ok := payload["cancelledreason"]; ok {
			if s, ok := v.(string); ok {
				p.CanceledReason = s
			}
		}
		if v, ok := payload["purchaseDate"]; ok {
			switch val := v.(type) {
			case string:
				if t, err := time.Parse(time.RFC3339, val); err == nil {
					p.PurchaseDate = t
				}
			case time.Time:
				p.PurchaseDate = val
			}
		}
		if v, ok := payload["canceledAt"]; ok {
			switch val := v.(type) {
			case string:
				if t, err := time.Parse(time.RFC3339, val); err == nil {
					p.CanceledAt = &t
				}
			case time.Time:
				t := val
				p.CanceledAt = &t
			}
		}
	}
	if p.UniqueCode == "" && strings.TrimSpace(purchasedTicketID) != "" {
		p.UniqueCode = purchasedTicketID
	}
	return p, nil
}

func refundFromRow(refundID, userID, orderID, status string, amount any, createdAt, updatedAt time.Time, metadata []byte) (*RefundRequest, error) {
	r := &RefundRequest{
		EventID:     "",
		TicketID:    "",
		UserID:      userID,
		UniqueCode:  "",
		RequestDate: createdAt,
		Status:      status,
		Amount:      int(coerceInt64(amount)),
	}
	payload := parseJSONMap(metadata)
	if len(payload) > 0 {
		if v, ok := payload["eventID"]; ok {
			if s, ok := v.(string); ok {
				r.EventID = s
			}
		}
		if v, ok := payload["eventid"]; ok {
			if s, ok := v.(string); ok {
				r.EventID = s
			}
		}
		if v, ok := payload["ticketID"]; ok {
			if s, ok := v.(string); ok {
				r.TicketID = s
			}
		}
		if v, ok := payload["ticketid"]; ok {
			if s, ok := v.(string); ok {
				r.TicketID = s
			}
		}
		if v, ok := payload["userID"]; ok {
			if s, ok := v.(string); ok {
				r.UserID = s
			}
		}
		if v, ok := payload["userid"]; ok {
			if s, ok := v.(string); ok {
				r.UserID = s
			}
		}
		if v, ok := payload["uniqueCode"]; ok {
			if s, ok := v.(string); ok {
				r.UniqueCode = s
			}
		}
		if v, ok := payload["uniquecode"]; ok {
			if s, ok := v.(string); ok {
				r.UniqueCode = s
			}
		}
		if v, ok := payload["amount"]; ok {
			r.Amount = int(coerceInt64(v))
		}
		if v, ok := payload["status"]; ok {
			if s, ok := v.(string); ok {
				r.Status = s
			}
		}
		if v, ok := payload["requestDate"]; ok {
			switch val := v.(type) {
			case string:
				if t, err := time.Parse(time.RFC3339, val); err == nil {
					r.RequestDate = t
				}
			case time.Time:
				r.RequestDate = val
			}
		}
	}
	if r.EventID == "" {
		r.EventID = orderID
	}
	if r.UniqueCode == "" {
		r.UniqueCode = orderID
	}
	return r, nil
}

func buildTicketMetadata(ticket Ticket) map[string]any {
	meta := map[string]any{
		"ticketid":    ticket.TicketID,
		"eventid":     ticket.EventID,
		"name":        ticket.Name,
		"title":       ticket.Name,
		"price":       ticket.Price,
		"currency":    ticket.Currency,
		"color":       ticket.Color,
		"quantity":    ticket.Quantity,
		"available":   ticket.Available,
		"total":       ticket.Total,
		"description": ticket.Description,
		"sold":        ticket.Sold,
		"seatstart":   ticket.SeatStart,
		"seatend":     ticket.SeatEnd,
		"entity_id":   ticket.EntityID,
		"entity_type": ticket.EntityType,
		"seats":       ticket.Seats,
	}
	return meta
}

func buildPurchasedTicketMetadata(p PurchasedTicket) map[string]any {
	meta := map[string]any{
		"eventid":         p.EventID,
		"ticketid":        p.TicketID,
		"userid":          p.UserID,
		"buyerName":       p.BuyerName,
		"uniqueCode":      p.UniqueCode,
		"status":          p.Status,
		"purchaseDate":    p.PurchaseDate.Format(time.RFC3339),
		"price":           p.Price,
		"canceled":        p.Canceled,
		"canceledReason":  p.CanceledReason,
		"cancelledreason": p.CanceledReason,
		"transferred":     p.Transferred,
		"transferredTo":   p.TransferredTo,
	}
	if p.CanceledAt != nil {
		meta["canceledAt"] = p.CanceledAt.Format(time.RFC3339)
	}
	return meta
}

func buildRefundMetadata(refund RefundRequest) map[string]any {
	meta := map[string]any{
		"eventID":     refund.EventID,
		"eventid":     refund.EventID,
		"ticketID":    refund.TicketID,
		"ticketid":    refund.TicketID,
		"userID":      refund.UserID,
		"userid":      refund.UserID,
		"uniqueCode":  refund.UniqueCode,
		"uniquecode":  refund.UniqueCode,
		"status":      refund.Status,
		"amount":      refund.Amount,
		"requestDate": refund.RequestDate.Format(time.RFC3339),
	}
	if refund.ProcessedAt != nil {
		meta["processedAt"] = refund.ProcessedAt.Format(time.RFC3339)
	}
	if refund.RefundedAt != nil {
		meta["refundedAt"] = refund.RefundedAt.Format(time.RFC3339)
	}
	return meta
}

// Tickets DB helpers
func FindTicketsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]Ticket) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT tickid, eventid, title, price, quantity, created_at, updated_at, metadata FROM ticks WHERE eventid = $1 ORDER BY created_at DESC`, eventID)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]Ticket, 0)
	for rows.Next() {
		var ticketID, rowEventID, title string
		var priceValue any
		var quantityValue any
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&ticketID, &rowEventID, &title, &priceValue, &quantityValue, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		t, err := ticketFromRow(ticketID, rowEventID, title, priceValue, quantityValue, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		items = append(items, *t)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*out = items
	return nil
}

func FindTicketByID(ctx context.Context, app *infra.Deps, eventID, ticketID string) (*Ticket, error) {
	if err := ensureTicketDB(app); err != nil {
		return nil, err
	}
	var ticketIDRow, rowEventID, title string
	var priceValue any
	var quantityValue any
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT tickid, eventid, title, price, quantity, created_at, updated_at, metadata FROM ticks WHERE eventid = $1 AND (tickid = $2 OR metadata->>'ticketid' = $2) LIMIT 1`, eventID, ticketID).Scan(&ticketIDRow, &rowEventID, &title, &priceValue, &quantityValue, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, err
	}
	return ticketFromRow(ticketIDRow, rowEventID, title, priceValue, quantityValue, createdAt, updatedAt, metadata)
}

func InsertTicketDB(ctx context.Context, app *infra.Deps, t Ticket) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if t.TicketID == "" {
		t.TicketID = utils.GenerateRandomString(12)
	}
	if t.EventID == "" {
		return errors.New("event id required")
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = t.CreatedAt
	}
	if t.Available == 0 {
		t.Available = t.Quantity
	}
	if t.Total == 0 {
		t.Total = t.Quantity
	}
	meta, err := json.Marshal(buildTicketMetadata(t))
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO ticks (tickid, eventid, title, price, quantity, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tickid) DO UPDATE SET
			eventid = EXCLUDED.eventid,
			title = EXCLUDED.title,
			price = EXCLUDED.price,
			quantity = EXCLUDED.quantity,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		t.TicketID,
		t.EventID,
		t.Name,
		float64(t.Price)/100.0,
		t.Quantity,
		t.CreatedAt,
		t.UpdatedAt,
		meta,
	)
	return err
}

func UpdateTicketDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if err := ensureTicketDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}
	var metadata []byte
	if q := normalizeTicketWhere(query); q != "" {
		fetchQuery := fmt.Sprintf(`SELECT metadata FROM ticks WHERE %s LIMIT 1`, q)
		if err := app.SQLDB.QueryRow(ctx, fetchQuery, args...).Scan(&metadata); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	base := parseJSONMap(metadata)
	for key, val := range update {
		base[key] = val
	}
	if v, ok := update["name"]; ok {
		base["name"] = v
		base["title"] = v
	}
	if v, ok := update["quantity"]; ok {
		base["quantity"] = v
	}
	if v, ok := update["available"]; ok {
		base["available"] = v
	}
	if v, ok := update["total"]; ok {
		base["total"] = v
	}
	if v, ok := update["price"]; ok {
		base["price"] = v
	}
	payload, err := json.Marshal(base)
	if err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx,
		`UPDATE ticks SET metadata = $1::jsonb, updated_at = NOW() WHERE `+normalizeTicketWhere(query),
		append([]any{payload}, args...)...,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func DeleteTicketDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if err := ensureTicketDB(app); err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM ticks WHERE `+normalizeTicketWhere(query), args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Purchased tickets helpers
func FindPurchasedTickets(ctx context.Context, app *infra.Deps, query string, args []any, out *[]PurchasedTicket) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT purchasedticketid, userid, eventid, ticketid, status, amount, created_at, updated_at, metadata FROM purticks WHERE `+normalizeTicketWhere(query)+` ORDER BY created_at DESC`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]PurchasedTicket, 0)
	for rows.Next() {
		var purchasedTicketID, userID, eventID, ticketID, status string
		var amount any
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&purchasedTicketID, &userID, &eventID, &ticketID, &status, &amount, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		p, err := purchasedTicketFromRow(purchasedTicketID, userID, eventID, ticketID, status, amount, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		items = append(items, *p)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*out = items
	return nil
}

func FindPurchasedTicketByUnique(ctx context.Context, app *infra.Deps, eventID, uniqueCode string) (*PurchasedTicket, error) {
	if err := ensureTicketDB(app); err != nil {
		return nil, err
	}
	var purchasedTicketID, userID, rowEventID, ticketID, status string
	var amount any
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT purchasedticketid, userid, eventid, ticketid, status, amount, created_at, updated_at, metadata FROM purticks WHERE eventid = $1 AND metadata->>'uniqueCode' = $2 LIMIT 1`, eventID, uniqueCode).Scan(&purchasedTicketID, &userID, &rowEventID, &ticketID, &status, &amount, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("purchased ticket not found")
		}
		return nil, err
	}
	return purchasedTicketFromRow(purchasedTicketID, userID, rowEventID, ticketID, status, amount, createdAt, updatedAt, metadata)
}

func InsertPurchasedTicket(ctx context.Context, app *infra.Deps, p PurchasedTicket) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if p.UniqueCode == "" {
		p.UniqueCode = utils.GenerateRandomString(12)
	}
	if p.Status == "" {
		p.Status = "paid"
	}
	if p.PurchaseDate.IsZero() {
		p.PurchaseDate = time.Now().UTC()
	}
	payload, err := json.Marshal(buildPurchasedTicketMetadata(p))
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO purticks (purchasedticketid, userid, eventid, ticketid, status, amount, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (purchasedticketid) DO UPDATE SET
			userid = EXCLUDED.userid,
			eventid = EXCLUDED.eventid,
			ticketid = EXCLUDED.ticketid,
			status = EXCLUDED.status,
			amount = EXCLUDED.amount,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		utils.GenerateRandomString(12),
		p.UserID,
		p.EventID,
		p.TicketID,
		p.Status,
		p.Price,
		p.PurchaseDate,
		p.PurchaseDate,
		payload,
	)
	return err
}

func UpdatePurchasedTicket(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if err := ensureTicketDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}
	var metadata []byte
	q := normalizeTicketWhere(query)
	if q != "" {
		if err := app.SQLDB.QueryRow(ctx, `SELECT metadata FROM purticks WHERE `+q+` LIMIT 1`, args...).Scan(&metadata); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	base := parseJSONMap(metadata)
	for key, val := range update {
		base[key] = val
	}
	payload, err := json.Marshal(base)
	if err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx,
		`UPDATE purticks SET metadata = $1::jsonb, updated_at = NOW() WHERE `+q,
		append([]any{payload}, args...)...,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// Refunds
func FindRefunds(ctx context.Context, app *infra.Deps, query string, args []any, out *[]RefundRequest) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT refundid, userid, orderid, amount, status, created_at, updated_at, metadata FROM refunds WHERE `+normalizeRefundWhere(query)+` ORDER BY created_at DESC`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]RefundRequest, 0)
	for rows.Next() {
		var refundID, userID, orderID, status string
		var amount any
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&refundID, &userID, &orderID, &amount, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		r, err := refundFromRow(refundID, userID, orderID, status, amount, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		items = append(items, *r)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*out = items
	return nil
}

func InsertRefund(ctx context.Context, app *infra.Deps, refund RefundRequest) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if refund.RequestDate.IsZero() {
		refund.RequestDate = time.Now().UTC()
	}
	if refund.UniqueCode == "" {
		refund.UniqueCode = fmt.Sprintf("refund-%d", time.Now().UnixNano())
	}
	refundID := utils.GenerateRandomString(12)
	payload, err := json.Marshal(buildRefundMetadata(refund))
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO refunds (refundid, userid, orderid, amount, status, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (refundid) DO UPDATE SET
			userid = EXCLUDED.userid,
			orderid = EXCLUDED.orderid,
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		refundID,
		refund.UserID,
		refund.EventID,
		refund.Amount,
		refund.Status,
		refund.RequestDate,
		refund.RequestDate,
		payload,
	)
	return err
}

// Events helper (events table is named literal "events")
func FindEventByID(ctx context.Context, app *infra.Deps, eventID string, out interface{}) error {
	if err := ensureTicketDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	var rowEventID, creatorID string
	if err := app.SQLDB.QueryRow(ctx, `SELECT eventid, userid FROM events WHERE eventid = $1 LIMIT 1`, eventID).Scan(&rowEventID, &creatorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("event not found")
		}
		return err
	}
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("out must be pointer to struct")
	}
	field := v.Elem().FieldByName("CreatorID")
	if field.IsValid() && field.CanSet() {
		field.SetString(creatorID)
	}
	if field.Kind() == 0 {
		// Fallback for anonymous struct fields
		for i := 0; i < v.Elem().NumField(); i++ {
			name := v.Elem().Type().Field(i).Name
			if strings.EqualFold(name, "CreatorID") || strings.EqualFold(name, "creatorid") {
				v.Elem().Field(i).SetString(creatorID)
				break
			}
		}
	}
	_ = rowEventID
	return nil
}
