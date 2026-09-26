package mechat

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var (
	MessagesTable = config.Tables.MessagesTable
	MereChatTable = config.Tables.MerechatTable
)

// ================= REPOSITORY (POSTGRESQL LOGIC) =================

func SQLdbEnsureChatAccess(ctx context.Context, app *infra.Deps, chatID, user string) error {
	query := "chatid = $1 AND $2 = ANY(participants)"
	args := []any{chatID, user}

	return app.SQLDB.FindOne(ctx, MereChatTable, query, args, &struct{}{})
}

func SQLdbFindChat(ctx context.Context, app *infra.Deps, query string, args []any, out *Chat) error {
	return app.SQLDB.FindOne(ctx, MereChatTable, query, args, out)
}

func SQLdbInsertChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	return app.SQLDB.InsertOne(ctx, MereChatTable, chat)
}

func SQLdbFindMessagesForChat(ctx context.Context, app *infra.Deps, chatID string, user string, limit, offset int) ([]Message, error) {
	if err := SQLdbEnsureChatAccess(ctx, app, chatID, user); err != nil {
		return nil, err
	}

	query := "chatid = $1 AND deleted IS NOT TRUE"
	args := []any{chatID}

	opts := sqldb.FindManyOptions{
		Limit:   limit,
		Offset:  offset,
		OrderBy: "created_at DESC",
	}

	var msgs []Message
	if err := app.SQLDB.FindManyWithOptions(ctx, MessagesTable, query, args, opts, &msgs); err != nil {
		return nil, err
	}
	if msgs == nil {
		msgs = make([]Message, 0)
	}
	return msgs, nil
}

func SQLdbFindChatByUser(ctx context.Context, app *infra.Deps, chatID, user string) (Chat, error) {
	var chat Chat
	query := "chatid = $1 AND $2 = ANY(participants)"
	args := []any{chatID, user}

	if err := app.SQLDB.FindOne(ctx, MereChatTable, query, args, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func SQLdbFindUserChats(ctx context.Context, app *infra.Deps, user string, offset, limit int) ([]Chat, error) {
	query := "$1 = ANY(participants)"
	args := []any{user}

	opts := sqldb.FindManyOptions{
		Offset:  offset,
		Limit:   limit,
		OrderBy: "updated_at DESC",
	}

	var chats []Chat
	if err := app.SQLDB.FindManyWithOptions(ctx, MereChatTable, query, args, opts, &chats); err != nil {
		return nil, err
	}
	if chats == nil {
		chats = make([]Chat, 0)
	}
	return chats, nil
}

func SQLdbPersistAttachmentMessage(ctx context.Context, app *infra.Deps, chatID, user string, msg *Message) error {
	if err := app.SQLDB.InsertOne(ctx, MessagesTable, msg); err != nil {
		return err
	}

	query := "chatid = $1"
	args := []any{chatID}
	update := map[string]any{
		"updated_at": SQLnowUTC(),
		"last_message": map[string]any{
			"text":      "[attachment]",
			"sender_id": user,
			"timestamp": time.Now(),
		},
	}

	_, _ = app.SQLDB.UpdateOne(ctx, MereChatTable, query, args, update)
	return nil
}

func SQLnowUTC() time.Time { return time.Now().UTC() }

func SQLdbUpdateLastMessage(ctx context.Context, app *infra.Deps, chatID string, msg *Message) {
	if msg == nil {
		return
	}

	preview := MessagePreview{
		Text:      msg.Content,
		UserID:    msg.UserID,
		Timestamp: msg.CreatedAt,
	}

	query := "chatid = $1"
	args := []any{chatID}
	update := map[string]any{
		"last_message": preview,
		"updated_at":   time.Now(),
	}

	_, _ = app.SQLDB.UpdateOne(ctx, MereChatTable, query, args, update)
}

func SQLdbInsertMessage(ctx context.Context, app *infra.Deps, msg *Message) error {
	return app.SQLDB.InsertOne(ctx, MessagesTable, msg)
}

func SQLdbEditMessage(ctx context.Context, app *infra.Deps, msgID, userID, newContent string) (*Message, error) {
	now := time.Now()
	query := "messageid = $1 AND userid = $2 AND deleted IS NOT TRUE"
	args := []any{msgID, userID}

	update := map[string]any{
		"content":   newContent,
		"edited_at": now,
	}

	var msg Message
	if err := app.SQLDB.FindOneAndUpdate(ctx, MessagesTable, query, args, update, &msg); err != nil {
		return nil, err
	}

	SQLdbUpdateLastMessage(ctx, app, msg.ChatID, &msg)
	return &msg, nil
}

func SQLdbDeleteMessage(ctx context.Context, app *infra.Deps, msgID, userID string) (*Message, error) {
	query := "messageid = $1 AND userid = $2"
	args := []any{msgID, userID}
	update := map[string]any{"deleted": true}

	var msg Message
	if err := app.SQLDB.FindOneAndUpdate(ctx, MessagesTable, query, args, update, &msg); err != nil {
		return nil, err
	}

	// Clear last message preview if this deleted message matches the current sender
	clearQuery := "chatid = $1 AND (last_message->>'sender_id') = $2"
	clearArgs := []any{msg.ChatID, msg.UserID}
	clearUpdate := map[string]any{"last_message": nil}

	_, _ = app.SQLDB.UpdateOne(ctx, MereChatTable, clearQuery, clearArgs, clearUpdate)

	return &msg, nil
}

func SQLdbMarkAsRead(ctx context.Context, app *infra.Deps, msgID, userID string) error {
	rawQuery := "UPDATE " + MessagesTable + " SET read_by = ARRAY_APPEND(read_by, $2) WHERE messageid = $1 AND NOT ($2 = ANY(read_by))"
	return app.SQLDB.QueryRaw(ctx, rawQuery, []any{msgID, userID}, nil)
}

func SQLdbUpdateReaction(ctx context.Context, app *infra.Deps, msgID, userID string, add bool) error {
	if add {
		rawQuery := "UPDATE " + MessagesTable + " SET reactions = ARRAY_APPEND(reactions, $2) WHERE messageid = $1 AND NOT ($2 = ANY(reactions))"
		return app.SQLDB.QueryRaw(ctx, rawQuery, []any{msgID, userID}, nil)
	}

	rawQuery := "UPDATE " + MessagesTable + " SET reactions = ARRAY_REMOVE(reactions, $2) WHERE messageid = $1"
	return app.SQLDB.QueryRaw(ctx, rawQuery, []any{msgID, userID}, nil)
}

func SQLdbGetChatParticipants(ctx context.Context, app *infra.Deps, chatID string) ([]string, error) {
	var chat Chat
	query := "chatid = $1"
	args := []any{chatID}

	if err := app.SQLDB.FindOne(ctx, MereChatTable, query, args, &chat); err != nil {
		return nil, err
	}
	return chat.Participants, nil
}

func SQLdbGetUnreadCountsPerChat(ctx context.Context, app *infra.Deps, user string) ([]Chat, map[string]int64, error) {
	var chats []Chat
	chatQuery := "$1 = ANY(participants)"
	chatArgs := []any{user}

	if err := app.SQLDB.FindMany(ctx, MereChatTable, chatQuery, chatArgs, &chats); err != nil {
		return nil, nil, err
	}

	countsMap := make(map[string]int64)
	if len(chats) == 0 {
		return chats, countsMap, nil
	}

	rawQuery := `
		SELECT chatid, COUNT(*) AS count
		FROM ` + MessagesTable + `
		WHERE $1 = ANY(chat_ids)
		  AND userid != $1
		  AND deleted IS NOT TRUE
		  AND NOT ($1 = ANY(read_by))
		GROUP BY chatid
	`

	var results []UnreadCountResult
	if err := app.SQLDB.QueryRaw(ctx, rawQuery, []any{user}, &results); err != nil {
		return chats, countsMap, nil
	}

	for _, res := range results {
		countsMap[res.ChatID] = res.Count
	}

	return chats, countsMap, nil
}

func SQLdbSearchMessages(ctx context.Context, app *infra.Deps, chatID, term string, limit, offset int) ([]Message, error) {
	query := "chatid = $1 AND deleted IS NOT TRUE"
	args := []any{chatID}

	if term != "" {
		query += " AND content ILIKE $2"
		args = append(args, "%"+term+"%")
	}

	opts := sqldb.FindManyOptions{
		Limit:   limit,
		Offset:  offset,
		OrderBy: "created_at DESC",
	}

	var msgs []Message
	if err := app.SQLDB.FindManyWithOptions(ctx, MessagesTable, query, args, opts, &msgs); err != nil {
		return nil, err
	}

	if msgs == nil {
		msgs = make([]Message, 0)
	}

	return msgs, nil
}
