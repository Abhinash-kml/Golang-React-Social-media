package models

import (
	"errors"
	"fmt"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
)

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

func InsertPost(uuId uuid.UUID, content, hashtag string) (bool, error) {
	row := db.Connection.QueryRow("INSERT INTO posts(userid, content, hashtag) VALUES($1, $2, $3) returning userid;", uuId, content, hashtag)
	var returnedId uuid.UUID
	if err := row.Scan(returnedId); err != nil {
		fmt.Println("Error occured in Inserting post. Error: ", err)
		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from insertion is not same as inserted uuid in InsertPost()")
	}

	return true, nil
}

func UpdatePost() {

}

func DeletePost() {

}

func GetPost() {

}
