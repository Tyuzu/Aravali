package newchat

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var (
	chatsTable    = config.Tables.ChatsTable
	messagesTable = config.Tables.MessagesTable
)

func SQLgetChatByID(ctx context.Context, app *infra.Deps, chatID string) (Chat, error) {
	var chat Chat
	query := "chatid = $1"
	args := []any{chatID}

	if err := app.SQLDB.FindOne(ctx, chatsTable, query, args, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func SQLgetChatForUser(ctx context.Context, app *infra.Deps, chatID, userID string) (Chat, error) {
	var chat Chat
	query := "chatid = $1 AND $2 = ANY(users)"
	args := []any{chatID, userID}

	if err := app.SQLDB.FindOne(ctx, chatsTable, query, args, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func SQLgetChatMessages(ctx context.Context, app *infra.Deps, chatID string) ([]Message, error) {
	var messages []Message
	query := "chatid = $1"
	args := []any{chatID}

	opts := sqldb.FindManyOptions{
		OrderBy: "created_at ASC",
	}
	if err := app.SQLDB.FindManyWithOptions(ctx, messagesTable, query, args, opts, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func SQLgetRoomMessages(ctx context.Context, app *infra.Deps, room string) ([]Message, error) {
	var messages []Message
	query := "room = $1"
	args := []any{room}

	opts := sqldb.FindManyOptions{
		OrderBy: "timestamp DESC",
		Limit:   20,
	}
	if err := app.SQLDB.FindManyWithOptions(ctx, messagesTable, query, args, opts, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func SQLinsertMessage(ctx context.Context, app *infra.Deps, msg Message) error {
	return app.SQLDB.InsertOne(ctx, messagesTable, msg)
}

func SQLupdateChatLastMessage(ctx context.Context, app *infra.Deps, chatID, userID string, timestamp time.Time, previewText string) error {
	if previewText == "" {
		return nil
	}

	query := "chatid = $1"
	args := []any{chatID}

	update := map[string]any{
		"last_message": MessagePreview{
			Text:      previewText,
			UserID:    userID,
			Timestamp: timestamp,
		},
		"updated_at": timestamp,
	}

	_, err := app.SQLDB.UpdateOne(ctx, chatsTable, query, args, update)
	return err
}

func SQLUpdatexMessage(userID string, id string, newContent string, app *infra.Deps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "messageid = $1 AND senderid = $2"
	args := []any{id, userID}
	update := map[string]any{"content": newContent}

	rowsAffected, err := app.SQLDB.UpdateOne(ctx, messagesTable, query, args, update)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized")
	}
	return nil
}

func SQLDeletexMessage(userID string, id string, app *infra.Deps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "messageid = $1 AND senderid = $2"
	args := []any{id, userID}

	rowsAffected, err := app.SQLDB.DeleteOne(ctx, messagesTable, query, args)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized")
	}
	return nil
}

func SQLfindMessageRoom(id string, app *infra.Deps) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var msg Message
	query := "messageid = $1"
	args := []any{id}

	err := app.SQLDB.FindOne(ctx, messagesTable, query, args, &msg)
	if err != nil {
		return "", err
	}
	return msg.Room, nil
}

func SQLfindMessageByID(ctx context.Context, app *infra.Deps, msgID string) (Message, error) {
	var msg Message
	query := "messageid = $1"
	args := []any{msgID}

	if err := app.SQLDB.FindOne(ctx, messagesTable, query, args, &msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func SQLupdateMessageText(ctx context.Context, app *infra.Deps, msgID, text string) error {
	query := "messageid = $1"
	args := []any{msgID}
	update := map[string]any{"text": text}

	_, err := app.SQLDB.UpdateOne(ctx, messagesTable, query, args, update)
	return err
}

func SQLdeleteMessageByID(ctx context.Context, app *infra.Deps, msgID string) error {
	query := "messageid = $1"
	args := []any{msgID}

	_, err := app.SQLDB.DeleteOne(ctx, messagesTable, query, args)
	return err
}

func SQLtouchChatUpdatedAt(ctx context.Context, app *infra.Deps, chatID string) error {
	query := "chatid = $1"
	args := []any{chatID}
	update := map[string]any{"updated_at": time.Now()}

	_, err := app.SQLDB.UpdateOne(ctx, chatsTable, query, args, update)
	return err
}

func SQLfindChatByUsers(ctx context.Context, app *infra.Deps, users []string) (Chat, error) {
	var chat Chat
	query := "users = $1"
	args := []any{users}

	if err := app.SQLDB.FindOne(ctx, chatsTable, query, args, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func SQLcreateChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	return app.SQLDB.InsertOne(ctx, chatsTable, chat)
}

func SQLgetUserChats(ctx context.Context, app *infra.Deps, userID string) ([]Chat, error) {
	var chats []Chat
	query := "$1 = ANY(users)"
	args := []any{userID}

	opts := sqldb.FindManyOptions{
		OrderBy: "updated_at DESC",
		Limit:   15,
	}
	if err := app.SQLDB.FindManyWithOptions(ctx, chatsTable, query, args, opts, &chats); err != nil {
		return nil, err
	}
	if chats == nil {
		return []Chat{}, nil
	}
	return chats, nil
}
