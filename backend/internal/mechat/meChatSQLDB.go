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
	MessagesTable = config.Tables.MessagesTable
	MereChatTable = config.Tables.MerechatTable
)

// ================= REPOSITORY (MONGODB LOGIC) =================
func SQLdbEnsureChatAccess(ctx context.Context, app *infra.Deps, chatID, user string) error {
	return app.DB.FindOne(ctx, MereChatTable, map[string]any{
		"chatid":       chatID,
		"participants": user,
	}, &struct{}{})
}

func SQLdbFindChat(ctx context.Context, app *infra.Deps, filter map[string]any, out *Chat) error {
	return app.DB.FindOne(ctx, MereChatTable, filter, out)
}

func SQLdbInsertChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	return app.DB.InsertOne(ctx, MereChatTable, chat)
}

func SQLdbFindMessagesForChat(ctx context.Context, app *infra.Deps, chatID string, user string, limit, skip int) ([]Message, error) {
	if err := dbEnsureChatAccess(ctx, app, chatID, user); err != nil {
		return nil, err
	}

	filter := map[string]any{
		"chatid":     chatID,
		"deleted_ne": true,
	}

	opts := db.FindManyOptions{
		Limit: limit,
		Skip:  skip,
		Sort:  []bson.E{{Key: "createdAt", Value: -1}},
	}

	var msgs []Message
	if err := app.DB.FindManyWithOptions(ctx, MessagesTable, filter, opts, &msgs); err != nil {
		return nil, err
	}
	if msgs == nil {
		msgs = make([]Message, 0)
	}
	return msgs, nil
}

func SQLdbFindChatByUser(ctx context.Context, app *infra.Deps, chatID, user string) (Chat, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, MereChatTable, map[string]any{
		"chatid":       chatID,
		"participants": user,
	}, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func SQLdbFindUserChats(ctx context.Context, app *infra.Deps, user string, skip, limit int) ([]Chat, error) {
	opts := db.FindManyOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  []bson.E{{Key: "updatedAt", Value: -1}},
	}

	var chats []Chat
	if err := app.DB.FindManyWithOptions(
		ctx,
		MereChatTable,
		map[string]any{"participants": user},
		opts,
		&chats,
	); err != nil {
		return nil, err
	}
	if chats == nil {
		chats = make([]Chat, 0)
	}
	return chats, nil
}

func SQLdbPersistAttachmentMessage(ctx context.Context, app *infra.Deps, chatID, user string, msg *Message) error {
	if err := app.DB.InsertOne(ctx, MessagesTable, msg); err != nil {
		return err
	}

	_, _ = app.DB.UpdateOne(
		ctx,
		MereChatTable,
		map[string]any{"chatid": chatID},
		map[string]any{
			"updatedAt": nowUTC(),
			"lastMessage": map[string]any{
				"text":      "[attachment]",
				"senderId":  user,
				"timestamp": time.Now(),
			},
		},
	)
	return nil
}

func SQLnowUTC() time.Time { return time.Now() }

func SQLdbUpdateLastMessage(ctx context.Context, app *infra.Deps, chatID string, msg *Message) {
	if msg == nil {
		return
	}

	preview := MessagePreview{
		Text:      msg.Content,
		UserID:    msg.UserID,
		Timestamp: msg.CreatedAt,
	}

	_, _ = app.DB.UpdateOne(ctx,
		MereChatTable,
		map[string]any{"chatid": chatID},
		map[string]any{
			"$set": map[string]any{
				"lastMessage": preview,
				"updatedAt":   time.Now(),
			},
		},
	)
}

func SQLdbInsertMessage(ctx context.Context, app *infra.Deps, msg *Message) error {
	return app.DB.InsertOne(ctx, MessagesTable, msg)
}

func SQLdbEditMessage(ctx context.Context, app *infra.Deps, msgID, userID, newContent string) (*Message, error) {
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
	if err := app.DB.FindOneAndUpdate(ctx, MessagesTable, filter, update, &msg); err != nil {
		return nil, err
	}

	dbUpdateLastMessage(ctx, app, msg.ChatID, &msg)
	return &msg, nil
}

func SQLdbDeleteMessage(ctx context.Context, app *infra.Deps, msgID, userID string) (*Message, error) {
	filter := map[string]any{
		"messageid": msgID,
		"userid":    userID,
	}
	update := map[string]any{
		"$set": map[string]any{"deleted": true},
	}

	var msg Message
	if err := app.DB.FindOneAndUpdate(ctx, MessagesTable, filter, update, &msg); err != nil {
		return nil, err
	}

	// Clear last message in chat preview if this was the last message
	_, _ = app.DB.UpdateOne(
		ctx,
		MereChatTable,
		map[string]any{
			"chatid":               msg.ChatID,
			"lastMessage.senderId": msg.UserID,
		},
		map[string]any{"$set": map[string]any{"lastMessage": nil}},
	)

	return &msg, nil
}

func SQLdbMarkAsRead(ctx context.Context, app *infra.Deps, msgID, userID string) error {
	return app.DB.AddToSet(
		ctx,
		MessagesTable,
		map[string]any{"messageid": msgID},
		"readBy",
		userID,
	)
}

func SQLdbUpdateReaction(ctx context.Context, app *infra.Deps, msgID, userID string, add bool) error {
	if add {
		return app.DB.AddToSet(
			ctx,
			MessagesTable,
			map[string]any{"messageid": msgID},
			"reactions",
			userID,
		)
	}

	_, err := app.DB.UpdateOne(
		ctx,
		MessagesTable,
		map[string]any{"messageid": msgID},
		map[string]any{
			"$pull": map[string]any{"reactions": userID},
		},
	)
	return err
}

func SQLdbGetChatParticipants(ctx context.Context, app *infra.Deps, chatID string) ([]string, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, MereChatTable, map[string]any{"chatid": chatID}, &chat); err != nil {
		return nil, err
	}
	return chat.Participants, nil
}

func SQLdbGetUnreadCountsPerChat(ctx context.Context, app *infra.Deps, user string) ([]Chat, map[string]int64, error) {
	var chats []Chat
	if err := app.DB.FindMany(ctx, MereChatTable, map[string]any{
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
	_ = app.DB.Aggregate(ctx, MessagesTable, pipeline, &results)

	countsMap := make(map[string]int64)
	for _, res := range results {
		countsMap[res.ChatID] = res.Count
	}

	return chats, countsMap, nil
}

func SQLdbSearchMessages(ctx context.Context, app *infra.Deps, chatID, term string, limit, skip int) ([]Message, error) {
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
	if err := app.DB.FindManyWithOptions(ctx, MessagesTable, filter, opts, &msgs); err != nil {
		return nil, err
	}

	if msgs == nil {
		msgs = make([]Message, 0)
	}

	return msgs, nil
}
