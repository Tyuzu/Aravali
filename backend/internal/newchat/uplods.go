package newchat

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/infra/mq"
	log "scav/utils/logger"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/utils"
)

// ------------------------- HTTP handlers -------------------------

type editPayload struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func EditMessageHandler(hub *Hub, app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload editPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if err := UpdatexMessage(userID, payload.ID, payload.Content, app); err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "forbidden", http.StatusForbidden)
			}
			return
		}

		room, err := findMessageRoom(payload.ID, app)
		if err != nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		out := outboundPayload{
			Action:    "edit",
			ID:        payload.ID,
			Content:   payload.Content,
			Timestamp: time.Now().Unix(),
		}
		if data, err := json.Marshal(out); err == nil {
			hub.broadcast <- broadcastMsg{Room: room, Data: data}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteMessageHandler(hub *Hub, app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		room, err := findMessageRoom(payload.ID, app)
		if err != nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		if err := DeletexMessage(userID, payload.ID, app); err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "forbidden", http.StatusForbidden)
			}
			return
		}

		out := outboundPayload{
			Action: "delete",
			ID:     payload.ID,
		}
		if data, err := json.Marshal(out); err == nil {
			hub.broadcast <- broadcastMsg{Room: room, Data: data}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UploadHandler(hub *Hub, app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var payload struct {
			Chat  string       `json:"chat"`
			Files []Attachment `json:"files"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if payload.Chat == "" {
			http.Error(w, "chat missing", http.StatusBadRequest)
			return
		}

		ts := time.Now().Unix()
		msg := Message{
			MessageID: utils.GenerateRandomDigitString(16),
			Room:      payload.Chat,
			SenderID:  userID,
			Content:   "",
			Files:     payload.Files,
			Timestamp: ts,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := insertMessage(ctx, app, msg); err != nil {
			log.Println("db error:", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		previewText := buildLastMessagePreview(msg.Content, "", nil, len(msg.Files))
		if err := updateChatLastMessage(ctx, app, payload.Chat, userID, time.Unix(msg.Timestamp, 0), previewText); err != nil {
			log.Println("db error updating chat preview:", err)
		}

		out := outboundPayload{
			Action:    "chat",
			ID:        msg.MessageID,
			Room:      msg.Room,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			Files:     msg.Files,
			Timestamp: msg.Timestamp,
		}

		data, _ := json.Marshal(out)
		hub.broadcast <- broadcastMsg{Room: msg.Room, Data: data}

		mqpayload, _ := json.Marshal(mqevent.FileAddedToChatPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.FileAddedToChatEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, data)
	}
}
