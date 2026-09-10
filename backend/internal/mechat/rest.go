package mechat

import (
	"encoding/json"
	"net/http"
	"scav/internal/verticals/media"
	"sort"
	"strconv"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
)

//
// ================= CHATS =================
//

func StartNewChat(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := getUser(r)

		var body struct {
			Participants []string `json:"participants"`
			EntityType   string   `json:"entityType"`
			EntityId     string   `json:"entityId"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			writeErr(w, 400, "invalid body")
			return
		}

		set := map[string]struct{}{user: {}}
		for _, p := range body.Participants {
			if p != "" {
				set[p] = struct{}{}
			}
		}

		participants := make([]string, 0, len(set))
		for p := range set {
			participants = append(participants, p)
		}
		sort.Strings(participants)

		filter := map[string]any{
			"participants": participants,
			"entitytype":   body.EntityType,
			"entityid":     body.EntityId,
		}

		var existing Chat
		if err := dbFindChat(ctx, app, filter, &existing); err == nil {
			utils.RespondWithJSON(w, http.StatusOK, existing)
			return
		}

		now := time.Now()
		chat := Chat{
			ChatID:       utils.GenerateRandomString(16),
			Participants: participants,
			EntityType:   body.EntityType,
			EntityId:     body.EntityId,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := dbInsertChat(ctx, app, chat); err != nil {
			writeErr(w, 500, "failed to create chat")
			return
		}

		mqpayload, _ := json.Marshal(mqevent.MechatCreatedPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.MechatCreatedEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, chat)
	}
}

func GetChatMessages(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := utils.GetUserIDFromRequest(r)

		chatID := strings.TrimSpace(utils.GetParam(r, "chatid"))
		if chatID == "" {
			writeErr(w, 400, "missing chat id")
			return
		}

		limit := 50
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}

		skip := 0
		if s := r.URL.Query().Get("skip"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v >= 0 {
				skip = v
			}
		}

		msgs, err := dbFindMessagesForChat(ctx, app, chatID, user, limit, skip)
		if err != nil {
			writeErr(w, 404, "not found or access denied")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, msgs)
	}
}

func GetChatByID(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := utils.GetUserIDFromRequest(r)

		chatID := utils.GetParam(r, "chatid")
		chat, err := dbFindChatByUser(ctx, app, chatID, user)
		if err != nil {
			writeErr(w, 404, "not found or access denied")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, chat)
	}
}

func GetUserChats(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := utils.GetUserIDFromRequest(r)

		skip := 0
		limit := 20

		if v := r.URL.Query().Get("skip"); v != "" {
			if i, err := strconv.Atoi(v); err == nil && i >= 0 {
				skip = i
			}
		}
		if v := r.URL.Query().Get("limit"); v != "" {
			if i, err := strconv.Atoi(v); err == nil && i > 0 {
				limit = i
			}
		}

		chats, err := dbFindUserChats(ctx, app, user, skip, limit)
		if err != nil {
			writeErr(w, 500, "failed to load chats")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, chats)
	}
}

func UploadAttachment(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := utils.GetUserIDFromRequest(r)
		chatID := utils.GetParam(r, "chatid")

		mediaID := r.FormValue("mediaid")
		savedName := r.FormValue("savedname")
		extn := r.FormValue("extn")
		mimeType := r.FormValue("mimeType")
		fileSizeS := r.FormValue("fileSize")

		if mediaID == "" || savedName == "" || mimeType == "" {
			writeErr(w, 400, "invalid upload payload")
			return
		}

		fileSize, _ := strconv.ParseInt(fileSizeS, 10, 64)

		mediaType := "file"
		if strings.HasPrefix(mimeType, "image/") {
			mediaType = "image"
		} else if strings.HasPrefix(mimeType, "video/") {
			mediaType = "video"
		} else if strings.HasPrefix(mimeType, "audio/") {
			mediaType = "audio"
		}

		if err := dbEnsureChatAccess(ctx, app, chatID, user); err != nil {
			writeErr(w, 403, "chat not found or access denied")
			return
		}

		now := time.Now()

		// FIX: generate MessageID BEFORE insert
		msg := &Message{
			MessageID: utils.GenerateRandomDigitString(16),
			ChatID:    chatID,
			UserID:    user,
			Content:   "",
			Media: &media.Media{
				MediaID:   mediaID,
				Type:      mediaType,
				URL:       savedName,
				MimeType:  mimeType,
				FileSize:  fileSize,
				Extn:      extn,
				CreatedAt: now,
			},
			FileURL:   savedName,
			FileType:  mediaType,
			ReadBy:    []string{user},
			Status:    "sent",
			CreatedAt: now,
		}

		if err := dbPersistAttachmentMessage(ctx, app, chatID, user, msg); err != nil {
			writeErr(w, 500, "failed to persist message")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, msg)
	}
}
