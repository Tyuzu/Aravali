// File: internal/mechat/meChatSQLDB.go

package mechat

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
)

var (
	MessagesTable = config.Tables.MessagesTable
	MereChatTable = config.Tables.MerechatTable
)

// ================= REPOSITORY (POSTGRESQL LOGIC) =================

type UnreadCountResult struct {
	ChatID string `json:"chatid"`
	Count  int64  `json:"count"`
}

func SQLnowUTC() time.Time { return nowUTC() }

func dbEnsureChatAccess(ctx context.Context, app *infra.Deps, chatID, user string) error {
	return nil
}

func dbFindChat(ctx context.Context, app *infra.Deps, query string, args []any, out *Chat) error {
	return nil
}

func dbInsertChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	return nil
}

func dbFindMessagesForChat(ctx context.Context, app *infra.Deps, chatID string, user string, limit, offset int) ([]Message, error) {
	return nil, nil
}

func dbFindChatByUser(ctx context.Context, app *infra.Deps, chatID, user string) (Chat, error) {
	var chat Chat
	return chat, nil
}

func dbFindUserChats(ctx context.Context, app *infra.Deps, user string, offset, limit int) ([]Chat, error) {
	return nil, nil
}

func dbPersistAttachmentMessage(ctx context.Context, app *infra.Deps, chatID, user string, msg *Message) error {
	return nil
}

func nowUTC() time.Time { return time.Now().UTC() }

func dbUpdateLastMessage(ctx context.Context, app *infra.Deps, chatID string, msg *Message) {
}

func dbInsertMessage(ctx context.Context, app *infra.Deps, msg *Message) error {
	return nil
}

func dbEditMessage(ctx context.Context, app *infra.Deps, msgID, userID, newContent string) (*Message, error) {
	var msg Message
	return &msg, nil
}

func dbDeleteMessage(ctx context.Context, app *infra.Deps, msgID, userID string) (*Message, error) {
	var msg Message
	return &msg, nil
}

func dbMarkAsRead(ctx context.Context, app *infra.Deps, msgID, userID string) error {
	return nil
}

func dbUpdateReaction(ctx context.Context, app *infra.Deps, msgID, userID string, add bool) error {
	return nil
}

func dbGetChatParticipants(ctx context.Context, app *infra.Deps, chatID string) ([]string, error) {
	var chat Chat
	return chat.Participants, nil
}

func dbGetUnreadCountsPerChat(ctx context.Context, app *infra.Deps, user string) ([]Chat, map[string]int64, error) {
	var chats []Chat

	countsMap := make(map[string]int64)
	if len(chats) == 0 {
		return chats, countsMap, nil
	}

	return chats, countsMap, nil
}

func dbSearchMessages(ctx context.Context, app *infra.Deps, chatID, term string, limit, offset int) ([]Message, error) {
	return nil, nil
}
