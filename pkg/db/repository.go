package db

import (
	"context"

	model "github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"
)

type Repository interface {
	Connect(context.Context)
	Disconnect(context.Context)
	GetUserWithId(context.Context, string) (*model.User, error)
	GetUserWithName(context.Context, string) (*model.User, error)
	GetUserWithEmail(context.Context, string) (*model.User, error)
	GetUsersWithAttribute(context.Context, string) ([]*model.User, error)
	GetAllUsers(context.Context) ([]*model.User, error)
	DeleteUserWithId(context.Context) (bool, error)
	DeleteUserWithName(context.Context) (bool, error)
	DeleteUserWithEmail(context.Context) (bool, error)
	DeleteUsersWithAttribute(string) (bool, error)
	DeleteAllUsers(context.Context) (bool, error)
	GetPostWithId(context.Context) (*model.Post, error)
	GetPostsOfUser(context.Context) ([]*model.Post, error)
	GetPostsOfHashtag(context.Context) ([]*model.Post, error)
	GetAllPosts(context.Context) ([]*model.Post, error)
	DeletePostWithId(context.Context) (bool, error)
	DeletePostsOfUser(context.Context) (bool, error)
	DeletePostsOfHashtag(context.Context) (bool, error)
	DeleteAllPosts(context.Context) (bool, error)
	GetCommentWithId(context.Context) (*model.Comment, error)
	GetCommentsOfPost(context.Context) ([]*model.Comment, error)
	DeleteCommentWithId(context.Context) (bool, error)
	DeleteCommentsOfPost(context.Context) (bool, error)
	GetMessagesOfUsersWithId(context.Context, string, string) ([]*model.Message, error)
	DeleteMessagesOfUsersWithId(context.Context, string, string) (bool, error)
	GetMessagesOfUserWithId(context.Context, string, string) ([]*model.Message, error)
	DeleteMessagesOfUserWithId(context.Context, string, string) (bool, error)
}
