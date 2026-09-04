package mechat

import (
	"context"
	"regexp"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/db"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	MessagesCollection = config.Collections.MessagesCollection
	MereChatCollection = config.Collections.MerechatCollection
)

//
// ================= REPOSITORY (MONGODB LOGIC) =================
//

type UnreadCountResult struct {
	ChatID string `bson:"_id" json:"chatid"`
	Count  int64  `bson:"count" json:"count"`
}

func dbEnsureChatAccess(ctx context.Context, app *infra.Deps, chatID, user string) error {
	return app.DB.FindOne(ctx, MereChatCollection, map[string]any{
		"chatid":       chatID,
		"participants": user,
	}, &struct{}{})
}

func dbUpdateLastMessage(ctx context.Context, app *infra.Deps, chatID string, msg *Message) {
	if msg == nil {
		return
	}

	preview := MessagePreview{
		Text:      msg.Content,
		UserID:    msg.UserID,
		Timestamp: msg.CreatedAt,
	}

	_, _ = app.DB.UpdateOne(ctx,
		MereChatCollection,
		map[string]any{"chatid": chatID},
		map[string]any{
			"$set": map[string]any{
				"lastMessage": preview,
				"updatedAt":   time.Now(),
			},
		},
	)
}

func dbInsertMessage(ctx context.Context, app *infra.Deps, msg *Message) error {
	return app.DB.InsertOne(ctx, MessagesCollection, msg)
}

func dbEditMessage(ctx context.Context, app *infra.Deps, msgID, userID, newContent string) (*Message, error) {
	now := time.Now()
	filter := map[string]any{
		"messageid": msgID,
		"userid":    userID,
		"deleted":   map[string]any{"$ne": true},
	}
	update := map[string]any{
		"$set": map[string]any{
			"content":  newContent,
			"editedAt": now,
		},
	}

	var msg Message
	if err := app.DB.FindOneAndUpdate(ctx, MessagesCollection, filter, update, &msg); err != nil {
		return nil, err
	}

	dbUpdateLastMessage(ctx, app, msg.ChatID, &msg)
	return &msg, nil
}

func dbDeleteMessage(ctx context.Context, app *infra.Deps, msgID, userID string) (*Message, error) {
	filter := map[string]any{
		"messageid": msgID,
		"userid":    userID,
	}
	update := map[string]any{
		"$set": map[string]any{"deleted": true},
	}

	var msg Message
	if err := app.DB.FindOneAndUpdate(ctx, MessagesCollection, filter, update, &msg); err != nil {
		return nil, err
	}

	// Clear last message in chat preview if this was the last message
	_, _ = app.DB.UpdateOne(
		ctx,
		MereChatCollection,
		map[string]any{
			"chatid":               msg.ChatID,
			"lastMessage.senderId": msg.UserID,
		},
		map[string]any{"$set": map[string]any{"lastMessage": nil}},
	)

	return &msg, nil
}

func dbMarkAsRead(ctx context.Context, app *infra.Deps, msgID, userID string) error {
	return app.DB.AddToSet(
		ctx,
		MessagesCollection,
		map[string]any{"messageid": msgID},
		"readBy",
		userID,
	)
}

func dbUpdateReaction(ctx context.Context, app *infra.Deps, msgID, userID string, add bool) error {
	if add {
		return app.DB.AddToSet(
			ctx,
			MessagesCollection,
			map[string]any{"messageid": msgID},
			"reactions",
			userID,
		)
	}

	_, err := app.DB.UpdateOne(
		ctx,
		MessagesCollection,
		map[string]any{"messageid": msgID},
		map[string]any{
			"$pull": map[string]any{"reactions": userID},
		},
	)
	return err
}

func dbGetChatParticipants(ctx context.Context, app *infra.Deps, chatID string) ([]string, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, MereChatCollection, map[string]any{"chatid": chatID}, &chat); err != nil {
		return nil, err
	}
	return chat.Participants, nil
}

func dbGetUnreadCountsPerChat(ctx context.Context, app *infra.Deps, user string) ([]Chat, map[string]int64, error) {
	var chats []Chat
	if err := app.DB.FindMany(ctx, MereChatCollection, map[string]any{
		"participants": user,
	}, &chats); err != nil {
		return nil, nil, err
	}

	chatIDs := make([]string, 0, len(chats))
	for _, c := range chats {
		chatIDs = append(chatIDs, c.ChatID)
	}

	pipeline := bson.A{
		map[string]any{
			"$match": map[string]any{
				"chatid":  map[string]any{"$in": chatIDs},
				"userid":  map[string]any{"$ne": user},
				"deleted": map[string]any{"$ne": true},
				"readBy":  map[string]any{"$ne": user},
			},
		},
		map[string]any{
			"$group": map[string]any{
				"_id":   "$chatid",
				"count": map[string]any{"$sum": 1},
			},
		},
	}

	var results []UnreadCountResult
	_ = app.DB.Aggregate(ctx, MessagesCollection, pipeline, &results)

	countsMap := make(map[string]int64)
	for _, res := range results {
		countsMap[res.ChatID] = res.Count
	}

	return chats, countsMap, nil
}

func dbSearchMessages(ctx context.Context, app *infra.Deps, chatID, term string, limit, skip int) ([]Message, error) {
	filter := map[string]any{
		"chatid":  chatID,
		"deleted": map[string]any{"$ne": true},
	}

	if term != "" {
		filter["content"] = map[string]any{
			"$regex":   regexp.QuoteMeta(term),
			"$options": "i",
		}
	}

	opts := db.FindManyOptions{
		Limit: limit,
		Skip:  skip,
		Sort:  []bson.E{{Key: "createdAt", Value: -1}},
	}

	var msgs []Message
	if err := app.DB.FindManyWithOptions(ctx, MessagesCollection, filter, opts, &msgs); err != nil {
		return nil, err
	}

	if msgs == nil {
		msgs = make([]Message, 0)
	}

	return msgs, nil
}
