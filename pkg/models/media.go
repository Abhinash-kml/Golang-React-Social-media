package models

import (
	"errors"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func InsertMediaWithId(logger *zap.Logger, postId uuid.UUID, url string) (bool, error) {
	row := db.Connection.QueryRow("INSERT INTO media(postid, url) VALUES($1, $2) returning postid;", postId, url)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "InsertMediaWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in InsertMedia()")
	}

	return true, nil
}

func UpdateMediaWithId(logger *zap.Logger, postId uuid.UUID, newUrl string) (bool, error) {
	row := db.Connection.QueryRow("UPDATE media SET url = $1 WHERE postid = $2 returning postid;", newUrl, postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "UpdateMediaWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in UpdateMedia()")
	}

	return true, nil
}

func DeleteMediaWithId(logger *zap.Logger, postId uuid.UUID) (bool, error) {
	row := db.Connection.QueryRow("DELETE FROM media WHERE postid = $1 returning postid;", postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "DeleteMediaWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in DeleteMedia()")
	}

	return true, nil
}

func GetMediaWithId(logger *zap.Logger, postId uuid.UUID) *Media {
	row := db.Connection.QueryRow("SELECT * FROM media WHERE postid = $1;", postId)
	media := &Media{}
	if err := row.Scan(media.PostID, media.Url); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "GetMediaWithId"),
			zap.String("Error", err.Error()))

		return nil
	}

	return media
}
