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
	if app == nil || app.SQLDB == nil || chatID == "" || user == "" {
		return nil
	}
	var exists bool
	err := app.SQLDB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM chats WHERE chatid = $1 AND $2 = ANY(participants))`, chatID, user).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return context.Canceled
	}
	return nil
}

func dbFindChat(ctx context.Context, app *infra.Deps, query string, args []any, out *Chat) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE `+query+` LIMIT 1`, args...)
	return row.Scan(&out.ChatID, &out.Users, &out.LastMessage, &out.ReadStatus, &out.Participants, &out.CreatedAt, &out.UpdatedAt, &out.EntityType, &out.EntityId, &out.LastSeq)
}

func dbInsertChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO chats (chatid, users, participants, createdat, updatedat, entitytype, entityid, lastseq) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (chatid) DO NOTHING`, chat.ChatID, chat.Users, chat.Participants, chat.CreatedAt, chat.UpdatedAt, chat.EntityType, chat.EntityId, chat.LastSeq)
	return err
}

func dbFindMessagesForChat(ctx context.Context, app *infra.Deps, chatID string, user string, limit, offset int) ([]Message, error) {
	if app == nil || app.SQLDB == nil || chatID == "" {
		return []Message{}, nil
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT messageid, chatid, sender, text, fileurl, filetype, createdat, replyto, sendername, avatarurl, content, media, editedat, deleted, readby, status, room, senderid, timestamp FROM messages WHERE chatid = $1 AND userid = $2 ORDER BY createdat DESC LIMIT $3 OFFSET $4`, chatID, user, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]Message, 0)
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.FileURL, &msg.FileType, &msg.CreatedAt, &msg.ReplyTo, &msg.SenderName, &msg.AvatarURL, &msg.Content, &msg.Media, &msg.EditedAt, &msg.Deleted, &msg.ReadBy, &msg.Status, &msg.RoomID, &msg.UserID, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func dbFindChatByUser(ctx context.Context, app *infra.Deps, chatID, user string) (Chat, error) {
	var chat Chat
	if app == nil || app.SQLDB == nil {
		return chat, nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE chatid = $1 AND $2 = ANY(participants) LIMIT 1`, chatID, user)
	if err := row.Scan(&chat.ChatID, &chat.Users, &chat.LastMessage, &chat.ReadStatus, &chat.Participants, &chat.CreatedAt, &chat.UpdatedAt, &chat.EntityType, &chat.EntityId, &chat.LastSeq); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func dbFindUserChats(ctx context.Context, app *infra.Deps, user string, offset, limit int) ([]Chat, error) {
	if app == nil || app.SQLDB == nil || user == "" {
		return []Chat{}, nil
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE $1 = ANY(participants) ORDER BY updatedat DESC LIMIT $2 OFFSET $3`, user, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Chat, 0)
	for rows.Next() {
		var item Chat
		if err := rows.Scan(&item.ChatID, &item.Users, &item.LastMessage, &item.ReadStatus, &item.Participants, &item.CreatedAt, &item.UpdatedAt, &item.EntityType, &item.EntityId, &item.LastSeq); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func dbPersistAttachmentMessage(ctx context.Context, app *infra.Deps, chatID, user string, msg *Message) error {
	if app == nil || app.SQLDB == nil || msg == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO messages (messageid, chatid, userid, content, fileurl, filetype, createdat, readby, status, roomid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, msg.MessageID, chatID, user, msg.Content, msg.FileURL, msg.FileType, msg.CreatedAt, msg.ReadBy, msg.Status, msg.RoomID)
	return err
}
func nowUTC() time.Time { return time.Now().UTC() }

func dbUpdateLastMessage(_ context.Context, _ *infra.Deps, _ string, _ *Message) {
}

func dbInsertMessage(_ context.Context, _ *infra.Deps, _ *Message) error {
	return nil
}

func dbEditMessage(_ context.Context, _ *infra.Deps, _, _, _ string) (*Message, error) {
	var msg Message
	return &msg, nil
}

func dbDeleteMessage(_ context.Context, _ *infra.Deps, _, _ string) (*Message, error) {
	var msg Message
	return &msg, nil
}

func dbMarkAsRead(_ context.Context, _ *infra.Deps, _, _ string) error {
	return nil
}

func dbUpdateReaction(_ context.Context, _ *infra.Deps, _, _ string, _ bool) error {
	return nil
}

func dbGetChatParticipants(_ context.Context, _ *infra.Deps, _ string) ([]string, error) {
	var chat Chat
	return chat.Participants, nil
}

func dbGetUnreadCountsPerChat(_ context.Context, _ *infra.Deps, _ string) ([]Chat, map[string]int64, error) {
	var chats []Chat

	countsMap := make(map[string]int64)
	if len(chats) == 0 {
		return chats, countsMap, nil
	}

	return chats, countsMap, nil
}

func dbSearchMessages(_ context.Context, _ *infra.Deps, _, _ string, _, _ int) ([]Message, error) {
	return nil, nil
}
