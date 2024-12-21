package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"uuid"`
	Fullname    string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Dob         string    `json:"dob"`
	Created_at  string    `json:"created_at"`
	Modified_at string    `json:"modified_at"`
	Lastlogin   string    `json:"lastlogin"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	City        string    `json:"city"`
}

func NewUser(name, email, password, dob, created_at, modified_at, last_login string) *User {
	return &User{
		Fullname:    name,
		Email:       email,
		Password:    password,
		Dob:         dob,
		Created_at:  created_at,
		Modified_at: modified_at,
		Lastlogin:   last_login,
	}
}

type Post struct {
	UserID  uuid.UUID `json:"userid"`
	Content string    `json:"content"`
	Hashtag string    `json:"hashtag"`
}

func NewPost(uuid uuid.UUID, content, hashtag string) *Post {
	return &Post{
		UserID:  uuid,
		Content: content,
		Hashtag: hashtag,
	}
}

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

type Media struct {
	PostID uuid.UUID `json:"postid"`
	Url    string    `json:"url"`
}

func NewMedia(postid uuid.UUID, url string) *Media {
	return &Media{
		PostID: postid,
		Url:    url,
	}
}
