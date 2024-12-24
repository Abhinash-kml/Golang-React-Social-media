package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Dob         string    `json:"dob"`
	Created_at  string    `json:"created_at"`
	Modified_at string    `json:"modified_at"`
	Lastlogin   string    `json:"last_login"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	City        string    `json:"city"`
	BanLevel    byte      `json:"ban_level"`
}

func NewUser(name, email, password, dob string, ban_level byte) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Dob:      dob,
		BanLevel: ban_level,
	}
}

type Post struct {
	UserId     uuid.UUID `json:"userid"`
	Title      string    `json:"title"`
	Created_at string    `json:"created_at"`
	Body       string    `json:"body"`
	Likes      int       `json:"likes"`
	Comments   int       `json:"comments"`
	MediaUrl   string    `json:"media_url"`
	Hashtag    string    `json:"hashtag"`
}

func NewPost(uuid uuid.UUID, content, hashtag, media_url string, likes, comments int) *Post {
	return &Post{
		UserId:   uuid,
		Body:     content,
		Hashtag:  hashtag,
		Likes:    likes,
		Comments: comments,
		MediaUrl: media_url,
	}
}

type Message struct {
	SenderID   uuid.UUID `json:"senderid"`
	RecieverID uuid.UUID `json:"recieverid"`
	Body       string    `json:"body"`
	Status     int       `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewMessage(senderid, receiverid uuid.UUID, content string) *Message {
	return &Message{
		SenderID:   senderid,
		RecieverID: receiverid,
		Body:       content,
		Timestamp:  time.Now(),
	}
}
