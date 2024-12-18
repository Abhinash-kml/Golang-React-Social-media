package models

import "github.com/google/uuid"

type Post struct {
	ID      int       `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
	Hashtag string    `json:"hashtag"`
}

func InsertPost() {

}

func UpdatePost() {

}

func DeletePost() {

}

func GetPost() {

}
