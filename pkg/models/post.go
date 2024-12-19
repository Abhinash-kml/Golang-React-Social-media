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
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in Inserting post. Error: ", err)
		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from insertion is not same as supplied uuid in InsertPost()")
	}

	return true, nil
}

func UpdatePost(uuId uuid.UUID, newContent, hashtag string) (bool, error) {
	row := db.Connection.QueryRow("UPDATE posts SET content = $1, hastag = $2 WHERE userid = $3 RETURNING userid;", newContent, hashtag, uuId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in updating post. Error: ", err)
		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from update in not same as supplied uuid in UpdatePost()")
	}

	return true, nil
}

func DeletePostWithId(uuId uuid.UUID) (bool, error) {
	row := db.Connection.QueryRow("DELETE FROM posts WHERE userid = $1 returning userid;", uuId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in deleting post. Error: ", err)
		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from delete is not same as supplied uuid in DeletePostWithId()")
	}

	return true, nil
}

func GetPost(uuId uuid.UUID) *Post {
	row := db.Connection.QueryRow("SELECT * FROM posts WHERE userid = $1;", uuId)
	post := &Post{}
	if err := row.Scan(post.UserID, post.Content, post.Hashtag); err != nil {
		fmt.Println("Error occured in getting post. Error: ", err)
		return nil
	}

	return post
}
