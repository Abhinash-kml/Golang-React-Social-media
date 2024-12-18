package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
)

type Message struct {
	SenderID   uuid.UUID `json:"sender_id"`
	RecieverID uuid.UUID `json:"reciever_id"`
	Content    string    `json:"content"`
}

func NewMessage(senderid, receiverid uuid.UUID, content string) *Message {
	return &Message{
		SenderID:   senderid,
		RecieverID: receiverid,
		Content:    content,
	}
}

func InsertMessageInDb(message *Message) (bool, error) {
	var sender_id, reciever_id uuid.UUID
	row := db.Connection.QueryRow("INSERT INTO messages(sender_id, reciever_id, content) VALUES($1, $2, $3) RETURNING sender_id, reciever_id;")
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		fmt.Println("Error scanning row in InsertMessageInDb(). Error: ", err)
		return false, err
	}

	if sender_id == message.SenderID && reciever_id == message.RecieverID {
		return true, nil
	}

	return false, errors.New("something wrong happened in InsertMessageInDb()")
}

func GetAllMessagesOfSenderAndReciever(sender_id, receiver_id uuid.UUID) []*Message {
	rows, err := db.Connection.Query("SELECT sender_id, reciever_id, content FROM messages WHERE sender_id = $1 AND reciever_id = $2;", sender_id, receiver_id)
	if err != nil {
		fmt.Println("Error querying rows in GetAllMessagesOfSenderAndReciever().")
		return nil
	}

	var (
		messages []*Message
		message  *Message = &Message{}
	)

	for rows.Next() {
		if err := rows.Scan(message.SenderID, message.RecieverID, message.Content); err != nil {
			fmt.Println("Error scanning rows in GetAllMessagesOfSenderAndReciever().")
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

func UpdateMessage(senderId, recieverId uuid.UUID, newContent string) (bool, error) {

}

func DeleteMessage(senderId, receiverId uuid.UUID) (bool, error) {

}
