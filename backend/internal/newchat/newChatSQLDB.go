// File: internal/newchat/newChatSQLDB.go

package newchat

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
)

var (
	chatsTable    = config.Tables.ChatsTable
	messagesTable = config.Tables.MessagesTable
)

func getChatByID(ctx context.Context, app *infra.Deps, chatID string) (Chat, error) {
	var chat Chat
	return chat, nil
}

func getChatForUser(ctx context.Context, app *infra.Deps, chatID, userID string) (Chat, error) {
	var chat Chat
	return chat, nil
}

func getChatMessages(ctx context.Context, app *infra.Deps, chatID string) ([]Message, error) {
	var messages []Message
	return messages, nil
}

func getRoomMessages(ctx context.Context, app *infra.Deps, room string) ([]Message, error) {
	var messages []Message
	return messages, nil
}

func insertMessage(ctx context.Context, app *infra.Deps, msg Message) error {

}

func updateChatLastMessage(ctx context.Context, app *infra.Deps, chatID, userID string, timestamp time.Time, previewText string) error {

	return err
}

func UpdatexMessage(userID string, id string, newContent string, app *infra.Deps) error {

	return nil
}

func DeletexMessage(userID string, id string, app *infra.Deps) error {

	return nil
}

func findMessageRoom(id string, app *infra.Deps) (string, error) {

	var msg Message
	return msg.Room, nil
}

func findMessageByID(ctx context.Context, app *infra.Deps, msgID string) (Message, error) {
	var msg Message
	return msg, nil
}

func updateMessageText(ctx context.Context, app *infra.Deps, msgID, text string) error {
	return err
}

func deleteMessageByID(ctx context.Context, app *infra.Deps, msgID string) error {
	return err
}

func touchChatUpdatedAt(ctx context.Context, app *infra.Deps, chatID string) error {
	return err
}

func findChatByUsers(ctx context.Context, app *infra.Deps, users []string) (Chat, error) {
	var chat Chat
	return chat, nil
}

func createChat(ctx context.Context, app *infra.Deps, chat Chat) error {

}

func getUserChats(ctx context.Context, app *infra.Deps, userID string) ([]Chat, error) {
	var chats []Chat
	return chats, nil
}
