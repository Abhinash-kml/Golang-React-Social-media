package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID    `json:"uuid,omitempty"`
	Name        string       `json:"name,omitempty"`
	Email       string       `json:"email,omitempty"`
	Password    string       `json:"password,omitempty"`
	Dob         string       `json:"dob,omitempty"`
	Created_at  string       `json:"created_at,omitempty"`
	Modified_at string       `json:"modified_at,omitempty"`
	Lastlogin   string       `json:"last_login,omitempty"`
	Country     string       `json:"country,omitempty"`
	State       string       `json:"state,omitempty"`
	City        string       `json:"city,omitempty"`
	AvatarUrl   string       `json:"avatar_url,omitempty"`
	BanLevel    string       `json:"ban_level,omitempty"` // (0 - no ban, 1 - text chat, 2 - voice chat, 3 - both voice and text, 4 - complete ip ban)
	BanDuration sql.NullTime `json:"ban_duration,omitempty"`
}

/*func NewUser(name, email, password, dob, ban_level string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Dob:      dob,
		BanLevel: ban_level,
	}
}*/

type Post struct {
	Id          int            `json:"id,omitempty"`
	UserId      uuid.UUID      `json:"userid,omitempty"`
	Title       string         `json:"title,omitempty"`
	Created_at  sql.NullTime   `json:"created_at,omitempty"`
	Modified_at sql.NullTime   `json:"modified_at,omitempty"`
	Body        string         `json:"body,omitempty"`
	Likes       int            `json:"likes,omitempty"`
	Comments    int            `json:"comments,omitempty"`
	MediaUrl    sql.NullString `json:"media_url,omitempty"`
	Hashtag     sql.NullString `json:"hashtag,omitempty"`
}

func NewPost(uuid uuid.UUID, content, hashtag, media_url string, likes, comments int) *Post {
	return &Post{
		UserId:   uuid,
		Body:     content,
		Hashtag:  sql.NullString{String: hashtag, Valid: true},
		Likes:    likes,
		Comments: comments,
		MediaUrl: sql.NullString{String: media_url, Valid: true},
	}
}

type Message struct {
	SenderID   uuid.UUID `json:"senderid,omitempty"`
	RecieverID uuid.UUID `json:"recieverid,omitempty"`
	Body       string    `json:"body,omitempty"`
	Status     int       `json:"status,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}

func NewMessage(senderid, receiverid uuid.UUID, content string) *Message {
	return &Message{
		SenderID:   senderid,
		RecieverID: receiverid,
		Body:       content,
		Timestamp:  time.Now(),
	}
}

type Comment struct {
	Id          uuid.UUID `json:"uuid,omitempty"`
	PostId      uuid.UUID `json:"postid,omitempty"`
	Body        string    `json:"body,omitempty"`
	Created_at  time.Time `json:"created_at,omitempty"`
	Modified_at time.Time `json:"modified_at,omitempty"`
}

func NewComment(postid uuid.UUID, body string, created_at time.Time) *Comment {
	return &Comment{
		PostId:     postid,
		Body:       body,
		Created_at: created_at,
	}
}
