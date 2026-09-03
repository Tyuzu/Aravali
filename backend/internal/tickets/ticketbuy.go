package tickets

import (
	"encoding/json"
	"fmt"
	"net/http"
	log "scav/utils/logger"
	"sync"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/internal/pay/stripe"
	"scav/utils"
)

// ------------------------------------------------------------------
// Realtime Event Channels with Cleanup
// ------------------------------------------------------------------

var eventUpdateChannels = struct {
	sync.RWMutex
	channels    map[string]chan map[string]any
	subscribers map[string]int // Track number of active subscribers per channel
}{
	channels:    make(map[string]chan map[string]any),
	subscribers: make(map[string]int),
}

func GetUpdatesChannel(eventId string) chan map[string]any {
	eventUpdateChannels.Lock()
	defer eventUpdateChannels.Unlock()

	if ch, ok := eventUpdateChannels.channels[eventId]; ok {
		eventUpdateChannels.subscribers[eventId]++
		return ch
	}

	ch := make(chan map[string]any, 10)
	eventUpdateChannels.channels[eventId] = ch
	eventUpdateChannels.subscribers[eventId] = 1
	return ch
}

// CloseUpdatesChannel marks a subscriber as finished. Channel is cleaned up when all subscribers are gone.
func CloseUpdatesChannel(eventId string) {
	eventUpdateChannels.Lock()
	defer eventUpdateChannels.Unlock()

	if count, ok := eventUpdateChannels.subscribers[eventId]; ok {
		count--
		if count <= 0 {
			if ch, ok := eventUpdateChannels.channels[eventId]; ok {
				close(ch)
			}
			delete(eventUpdateChannels.channels, eventId)
			delete(eventUpdateChannels.subscribers, eventId)
		} else {
			eventUpdateChannels.subscribers[eventId] = count
		}
	}
}

// ------------------------------------------------------------------
// Stripe Session
// ------------------------------------------------------------------

func CreateTicketPaymentSession(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ticketId := utils.GetParam(r, "ticketid")
		eventId := utils.GetParam(r, "eventid")

		var body struct {
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Quantity < 1 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		session, err := stripe.CreateTicketSession(ticketId, eventId, body.Quantity)
		if err != nil {
			http.Error(w, "Failed to create payment session", http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.TicketPaymentSessionCreatedEvent, mqevent.TicketPaymentSessionCreatedPayload{}); err != nil {
			log.Printf("failed to publish ticket payment session created event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"paymentUrl": session.URL,
				"eventId":    session.EventID,
				"ticketId":   session.TicketID,
				"stock":      session.Quantity,
			},
		})
	}
}

// ------------------------------------------------------------------
// SSE Updates
// ------------------------------------------------------------------

func EventUpdates(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventId := utils.GetParam(r, "eventId")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := GetUpdatesChannel(eventId)
		defer CloseUpdatesChannel(eventId)

		for {
			select {
			case update := <-ch:
				data, _ := json.Marshal(update)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}

func BroadcastTicketUpdate(eventId, ticketId string, remaining int) {
	ch := GetUpdatesChannel(eventId)
	select {
	case ch <- map[string]any{
		"type":             "ticket_update",
		"ticketId":         ticketId,
		"remainingTickets": remaining,
	}:
	default:
		log.Printf("event %s update channel full", eventId)
	}
}
