package db

import (
	"context"

	model "github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	Connect(context.Context)
	Disconnect(context.Context)
	GetUserWithId(context.Context, string) (*model.User, error)
	GetUserWithName(context.Context, string) (*model.User, error)
	GetUserWithEmail(context.Context, string) (*model.User, error)
	GetUsersWithAttribute(context.Context, string) ([]*model.User, error)
	GetAllUsers(context.Context) ([]*model.User, error)
	InsertUser(context.Context, string, string, string, string, string, string, string, string) (bool, error)
	DeleteUserWithId(context.Context) (bool, error)
	DeleteUserWithName(context.Context) (bool, error)
	DeleteUserWithEmail(context.Context) (bool, error)
	DeleteUsersWithAttribute(string) (bool, error)
	DeleteAllUsers(context.Context) (bool, error)
	InsertPost(context.Context, uuid.UUID, string, string, string) (bool, error)
	GetPostWithId(context.Context) (*model.Post, error)
	GetPostsOfUser(context.Context) ([]*model.Post, error)
	GetPostsOfHashtag(context.Context) ([]*model.Post, error)
	GetAllPosts(context.Context) ([]*model.Post, error)
	UpdatePostWithId(context.Context, uuid.UUID, string, string) (bool, error)
	DeletePostWithId(context.Context) (bool, error)
	DeletePostsOfUser(context.Context, uuid.UUID) (bool, int, error)
	DeletePostsOfHashtag(context.Context, string) (bool, int, error)
	DeleteAllPosts(context.Context) (bool, int, error)
	GetCommentWithId(context.Context) (*model.Comment, error)
	GetCommentsOfPost(context.Context) ([]*model.Comment, error)
	DeleteCommentWithId(context.Context) (bool, error)
	DeleteCommentsOfPost(context.Context) (bool, error)
	InsertMessageIntoConversation(context.Context, uuid.UUID, uuid.UUID, *model.Message) (bool, error)
	GetAllMessagesOfConversation(context.Context, string, string) ([]*model.Message, error)
	UpdateMessageOfConversation(context.Context, uuid.UUID, uuid.UUID, int)
	DeleteMessagesOfUsersWithId(context.Context, string, string) (bool, error)
	GetMessagesOfUserWithId(context.Context, string, string) ([]*model.Message, error)
	DeleteMessagesOfUserWithId(context.Context, string, string) (bool, error)
	GetAllMessagesInDB(context.Context) ([]*model.Message, error)
}
