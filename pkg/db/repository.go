package db

import model "github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"

type Repository interface {
	GetUserWithId() (*model.User, error)
	GetUserWithName() (*model.User, error)
	GetUserWithEmail() (*model.User, error)
	GetUsersWithAttribute(string) ([]*model.User, error)
	GetAllUsers() ([]*model.User, error)
	DeleteUserWithId() (bool, error)
	DeleteUserWithName() (bool, error)
	DeleteUserWithEmail() (bool, error)
	DeleteUsersWithAttribute(bool, error)
	DeleteAllUsers() (bool, error)
	GetPostWithId() (*model.Post, error)
	GetPostsOfUser() ([]*model.Post, error)
	GetPostsOfHashtag() ([]*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	DeletePostWithId() (bool, error)
	DeletePostsOfUser() (bool, error)
	DeletePostsOfHashtag() (bool, error)
	DeleteAllPosts() (bool, error)
	GetCommentWithId() (*model.Comment, error)
	GetCommentsOfPost() ([]*model.Comment, error)
	DeleteCommentWithId() (bool, error)
	DeleteCommentsOfPost() (bool, error)
	GetMessagesOfUsersWithId(string, string) ([]*model.Message, error)
	DeleteMessagesOfUsersWithId(string, string) (bool, error)
	GetMessagesOfUserWithId(string, string) ([]*model.Message, error)
	DeleteMessagesOfUserWithId(string, string) (bool, error)
}
