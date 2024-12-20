package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Message struct {
	SenderID   uuid.UUID `json:"senderid"`
	RecieverID uuid.UUID `json:"recieverid"`
	Content    string    `json:"content"`
}

func NewMessage(senderid, receiverid uuid.UUID, content string) *Message {
	return &Message{
		SenderID:   senderid,
		RecieverID: receiverid,
		Content:    content,
	}
}

func InsertMessageInDb(logger *zap.Logger, message *Message) (bool, error) {
	var sender_id, reciever_id uuid.UUID
	row := db.Connection.QueryRow("INSERT INTO messages(sender_id, reciever_id, content) VALUES($1, $2, $3) RETURNING sender_id, reciever_id;")
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "InsertMessageInDb"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id == message.SenderID && reciever_id == message.RecieverID {
		return true, nil
	}

	return false, errors.New("something wrong happened in InsertMessageInDb()")
}

func GetAllMessagesOfSenderAndReciever(logger *zap.Logger, sender_id, receiver_id uuid.UUID) []*Message {
	rows, err := db.Connection.Query("SELECT sender_id, reciever_id, content FROM messages WHERE sender_id IN($1, $2) AND reciever_id IN($1, $2);", sender_id, receiver_id)
	if err != nil {
		logger.Error("Error scanning rows",
			zap.String("function", "GetAllMessagesOfSenderAndReciever"),
			zap.String("Error", err.Error()))
		return nil
	}

	var (
		messages []*Message
		message  *Message = &Message{}
	)

	for rows.Next() {
		if err := rows.Scan(message.SenderID, message.RecieverID, message.Content); err != nil {
			logger.Error("Error scanning row",
				zap.String("function", "GetAllMessagesOfSenderAndReciever"),
				zap.String("Error", err.Error()))
			return nil
		}

		messages = append(messages, message)
	}

	return messages
}

func GetAllMessagesInDB() ([]*Message, error) {
	rows, err := db.Connection.Query("SELECT * FROM messages;")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
	}

	var (
		messages []*Message
		message  *Message = &Message{}
	)

	for rows.Next() {
		if err := rows.Scan(message.SenderID, message.RecieverID, message.Content); err != nil {
			fmt.Println(err)
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func UpdateMessageWithId(logger *zap.Logger, senderId, recieverId uuid.UUID, newContent string) (bool, error) {
	row := db.Connection.QueryRow("UPDATE messages SET content = $1 WHERE sender_id = $2 AND reciever_id = $3 RETURNING sender_id, receiver_id;", newContent, senderId, recieverId)
	var sender_id, reciever_id uuid.UUID
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "UpdateMessageWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id != senderId || reciever_id != recieverId {
		return false, errors.New("something wrong happened in UpdateMessageWithId()")
	}

	return true, nil
}

func DeleteMessage(logger *zap.Logger, senderId, receiverId uuid.UUID, content string) (bool, error) {
	row := db.Connection.QueryRow("DELETE FROM messages WHERE sender_id = $1 AND receiver_id = $2 AND content = $3 RETURNING sender_id, receiver_id;", senderId, receiverId, content)
	var sender_id, receiver_id uuid.UUID
	if err := row.Scan(&sender_id, &receiver_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", ""),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id != senderId || receiver_id != receiverId {
		return false, errors.New("something wrong happended in DeleteMessage()")
	}

	return true, nil
}
