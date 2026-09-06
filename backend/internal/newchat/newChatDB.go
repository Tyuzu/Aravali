package newchat

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	chatsCollection    = config.Collections.ChatsCollection
	messagesCollection = config.Collections.MessagesCollection
)

func getChatByID(ctx context.Context, app *infra.Deps, chatID string) (Chat, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, chatsCollection, map[string]any{"chatid": chatID}, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func getChatForUser(ctx context.Context, app *infra.Deps, chatID, userID string) (Chat, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, chatsCollection, map[string]any{
		"chatid": chatID,
		"users":  map[string]any{"$in": []string{userID}},
	}, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func getChatMessages(ctx context.Context, app *infra.Deps, chatID string) ([]Message, error) {
	var messages []Message
	opts := db.FindManyOptions{
		Sort: []bson.E{{Key: "createdAt", Value: 1}},
	}
	if err := app.DB.FindManyWithOptions(ctx, messagesCollection, map[string]any{"chatid": chatID}, opts, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func getRoomMessages(ctx context.Context, app *infra.Deps, room string) ([]Message, error) {
	var messages []Message
	opts := db.FindManyOptions{
		Sort:  []bson.E{{Key: "timestamp", Value: -1}},
		Limit: 20,
	}
	if err := app.DB.FindManyWithOptions(ctx, messagesCollection, map[string]any{"room": room}, opts, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func insertMessage(ctx context.Context, app *infra.Deps, msg Message) error {
	return app.DB.InsertOne(ctx, messagesCollection, msg)
}

func updateChatLastMessage(ctx context.Context, app *infra.Deps, chatID, userID string, timestamp time.Time, previewText string) error {
	if previewText == "" {
		return nil
	}
	update := map[string]any{
		"$set": map[string]any{
			"lastMessage": MessagePreview{
				Text:      previewText,
				UserID:    userID,
				Timestamp: timestamp,
			},
			"updatedAt": timestamp,
		},
	}
	_, err := app.DB.UpdateOne(ctx, chatsCollection, map[string]any{"chatid": chatID}, update)
	return err
}

func UpdatexMessage(userID string, id string, newContent string, app *infra.Deps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := map[string]any{
		"messageid": id,
		"senderid":  userID,
	}
	update := map[string]any{
		"$set": map[string]any{"content": newContent},
	}

	_, err := app.DB.UpdateOne(ctx, messagesCollection, filter, update)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return errors.New("message not found or unauthorized")
		}
		return err
	}
	return nil
}

func DeletexMessage(userID string, id string, app *infra.Deps) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := map[string]any{
		"messageid": id,
		"senderid":  userID,
	}

	_, err := app.DB.DeleteOne(ctx, messagesCollection, filter)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return errors.New("message not found or unauthorized")
		}
		return err
	}
	return nil
}

func findMessageRoom(id string, app *infra.Deps) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var msg Message
	err := app.DB.FindOne(ctx, messagesCollection, map[string]any{"messageid": id}, &msg)
	if err != nil {
		return "", err
	}
	return msg.Room, nil
}

func findMessageByID(ctx context.Context, app *infra.Deps, msgID string) (Message, error) {
	var msg Message
	if err := app.DB.FindOne(ctx, messagesCollection, map[string]any{"messageid": msgID}, &msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func updateMessageText(ctx context.Context, app *infra.Deps, msgID, text string) error {
	update := map[string]any{"$set": map[string]any{"text": text}}
	_, err := app.DB.UpdateOne(ctx, messagesCollection, map[string]any{"messageid": msgID}, update)
	return err
}

func deleteMessageByID(ctx context.Context, app *infra.Deps, msgID string) error {
	_, err := app.DB.DeleteOne(ctx, messagesCollection, map[string]any{"messageid": msgID})
	return err
}

func touchChatUpdatedAt(ctx context.Context, app *infra.Deps, chatID string) error {
	_, err := app.DB.UpdateOne(ctx, chatsCollection, map[string]any{"chatid": chatID}, map[string]any{
		"$set": map[string]any{"updatedAt": time.Now()},
	})
	return err
}

func findChatByUsers(ctx context.Context, app *infra.Deps, users []string) (Chat, error) {
	var chat Chat
	if err := app.DB.FindOne(ctx, chatsCollection, map[string]any{"users": users}, &chat); err != nil {
		return Chat{}, err
	}
	return chat, nil
}

func createChat(ctx context.Context, app *infra.Deps, chat Chat) error {
	return app.DB.InsertOne(ctx, chatsCollection, chat)
}

func getUserChats(ctx context.Context, app *infra.Deps, userID string) ([]Chat, error) {
	var chats []Chat
	opts := db.FindManyOptions{
		Sort:  []bson.E{{Key: "updatedAt", Value: -1}},
		Limit: 15,
	}
	if err := app.DB.FindManyWithOptions(ctx, chatsCollection, map[string]any{"users": map[string]any{"$in": []string{userID}}}, opts, &chats); err != nil {
		return nil, err
	}
	if chats == nil {
		return []Chat{}, nil
	}
	return chats, nil
}
