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
	if app == nil || app.SQLDB == nil {
		return chat, nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE chatid = $1 LIMIT 1`, chatID)
	if err := row.Scan(&chat.ChatID, &chat.Users, &chat.LastMessage, &chat.ReadStatus, &chat.Participants, &chat.CreatedAt, &chat.UpdatedAt, &chat.EntityType, &chat.EntityId, &chat.LastSeq); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func getChatForUser(ctx context.Context, app *infra.Deps, chatID, userID string) (Chat, error) {
	var chat Chat
	if app == nil || app.SQLDB == nil {
		return chat, nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE chatid = $1 AND $2 = ANY(participants) LIMIT 1`, chatID, userID)
	if err := row.Scan(&chat.ChatID, &chat.Users, &chat.LastMessage, &chat.ReadStatus, &chat.Participants, &chat.CreatedAt, &chat.UpdatedAt, &chat.EntityType, &chat.EntityId, &chat.LastSeq); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func getChatMessages(ctx context.Context, app *infra.Deps, chatID string) ([]Message, error) {
	var messages []Message
	if app == nil || app.SQLDB == nil || chatID == "" {
		return messages, nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT messageid, chatid, sender, text, fileurl, filetype, createdat, replyto, sendername, avatarurl, content, media, editedat, deleted, readby, status, room, senderid, timestamp FROM messages WHERE chatid = $1 ORDER BY createdat ASC`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.FileURL, &msg.FileType, &msg.CreatedAt, &msg.ReplyTo, &msg.SenderName, &msg.AvatarURL, &msg.Content, &msg.Media, &msg.EditedAt, &msg.Deleted, &msg.ReadBy, &msg.Status, &msg.Room, &msg.SenderID, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func getRoomMessages(ctx context.Context, app *infra.Deps, room string) ([]Message, error) {
	return getChatMessages(ctx, app, room)
}

func insertMessage(ctx context.Context, app *infra.Deps, msg Message) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO messages (messageid, chatid, sender, text, fileurl, filetype, createdat, replyto, sendername, avatarurl, content, media, editedat, deleted, readby, status, room, senderid, timestamp) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`, msg.MessageID, msg.ChatID, msg.UserID, msg.Text, msg.FileURL, msg.FileType, msg.CreatedAt, msg.ReplyTo, msg.SenderName, msg.AvatarURL, msg.Content, msg.Media, msg.EditedAt, msg.Deleted, msg.ReadBy, msg.Status, msg.Room, msg.SenderID, msg.Timestamp)
	return err
}

func updateChatLastMessage(ctx context.Context, app *infra.Deps, chatID, userID string, timestamp time.Time, previewText string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE chats SET lastmessage = $1, updatedat = NOW() WHERE chatid = $2 AND $3 = ANY(participants)`, previewText, chatID, userID)
	_ = timestamp
	return err
}

func UpdatexMessage(userID string, id string, newContent string, app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(context.Background(), `UPDATE messages SET content = $1, editedat = NOW() WHERE messageid = $2 AND senderid = $3`, newContent, id, userID)
	return err
}

func DeletexMessage(userID string, id string, app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(context.Background(), `DELETE FROM messages WHERE messageid = $1 AND senderid = $2`, id, userID)
	return err
}

func findMessageRoom(id string, app *infra.Deps) (string, error) {
	if app == nil || app.SQLDB == nil || id == "" {
		return "", nil
	}
	var room string
	err := app.SQLDB.QueryRow(context.Background(), `SELECT room FROM messages WHERE messageid = $1 LIMIT 1`, id).Scan(&room)
	return room, err
}

func findMessageByID(ctx context.Context, app *infra.Deps, msgID string) (Message, error) {
	var msg Message
	if app == nil || app.SQLDB == nil {
		return msg, nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT messageid, chatid, sender, text, fileurl, filetype, createdat, replyto, sendername, avatarurl, content, media, editedat, deleted, readby, status, room, senderid, timestamp FROM messages WHERE messageid = $1 LIMIT 1`, msgID)
	if err := row.Scan(&msg.MessageID, &msg.ChatID, &msg.UserID, &msg.Text, &msg.FileURL, &msg.FileType, &msg.CreatedAt, &msg.ReplyTo, &msg.SenderName, &msg.AvatarURL, &msg.Content, &msg.Media, &msg.EditedAt, &msg.Deleted, &msg.ReadBy, &msg.Status, &msg.Room, &msg.SenderID, &msg.Timestamp); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func updateMessageText(ctx context.Context, app *infra.Deps, msgID, text string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE messages SET content = $1, editedat = NOW() WHERE messageid = $2`, text, msgID)
	return err
}

func deleteMessageByID(ctx context.Context, app *infra.Deps, msgID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `DELETE FROM messages WHERE messageid = $1`, msgID)
	return err
}

func touchChatUpdatedAt(ctx context.Context, app *infra.Deps, chatID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE chats SET updatedat = NOW() WHERE chatid = $1`, chatID)
	return err
}

func findChatByUsers(ctx context.Context, app *infra.Deps, users []string) (Chat, error) {
	var chat Chat
	if app == nil || app.SQLDB == nil || len(users) == 0 {
		return chat, nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE $1 = ANY(participants) ORDER BY updatedat DESC LIMIT 1`, users[0])
	if err := row.Scan(&chat.ChatID, &chat.Users, &chat.LastMessage, &chat.ReadStatus, &chat.Participants, &chat.CreatedAt, &chat.UpdatedAt, &chat.EntityType, &chat.EntityId, &chat.LastSeq); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func createChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO chats (chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, chat.ChatID, chat.Users, chat.LastMessage, chat.ReadStatus, chat.Participants, chat.CreatedAt, chat.UpdatedAt, chat.EntityType, chat.EntityId, chat.LastSeq)
	return err
}

func getUserChats(ctx context.Context, app *infra.Deps, userID string) ([]Chat, error) {
	var chats []Chat
	if app == nil || app.SQLDB == nil || userID == "" {
		return chats, nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT chatid, users, lastmessage, readstatus, participants, createdat, updatedat, entitytype, entityid, lastseq FROM chats WHERE $1 = ANY(participants) ORDER BY updatedat DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ChatID, &chat.Users, &chat.LastMessage, &chat.ReadStatus, &chat.Participants, &chat.CreatedAt, &chat.UpdatedAt, &chat.EntityType, &chat.EntityId, &chat.LastSeq); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}
