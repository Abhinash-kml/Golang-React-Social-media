package models

import (
	"errors"
	"fmt"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
)

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

func InsertMedia(postId uuid.UUID, url string) (bool, error) {
	row := db.Connection.QueryRow("INSERT INTO media(postid, url) VALUES($1, $2) returning postid;", postId, url)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in scanning row in InsertMedia(). Error: ", err)
		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in InsertMedia()")
	}

	return true, nil
}

func UpdateMedia(postId uuid.UUID, newUrl string) (bool, error) {
	row := db.Connection.QueryRow("UPDATE media SET url = $1 WHERE postid = $2 returning postid;", newUrl, postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in scanning row in UpdateMedia(). Error: ", err)
		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in UpdateMedia()")
	}

	return true, nil
}

func DeleteMedia(postId uuid.UUID) (bool, error) {
	row := db.Connection.QueryRow("DELETE FROM media WHERE postid = $1 returning postid;", postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		fmt.Println("Error occured in scanning row in DeleteMedia(). Error: ", err)
		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in DeleteMedia()")
	}

	return true, nil
}

func GetMedia(postId uuid.UUID) *Media {
	row := db.Connection.QueryRow("SELECT * FROM media WHERE postid = $1;", postId)
	media := &Media{}
	if err := row.Scan(media.PostID, media.Url); err != nil {
		fmt.Println("Error occured in scanning row in GetMedia(). Error: ", err)
		return nil
	}

	return media
}
