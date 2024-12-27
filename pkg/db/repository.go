package db

import (
	"context"

	model "github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	Connect()
	Disconnect()
	GetUserWithId(ctx context.Context, uuid string) (*model.User, error)
	GetUserWithName(ctx context.Context, name string) (*model.User, error)
	GetUserWithEmail(ctx context.Context, email string) (*model.User, error)
	GetUsersWithAttribute(ctx context.Context, attibuteType string, value string) ([]*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	InsertUser(ctx context.Context, fullname, email, password, dob, country, state, city, avatar_url string) (bool, error)
	UpdateUserWithId(ctx context.Context, userid uuid.UUID, name, email, country, state string) (bool, error)
	UpdateUserWithName(ctx context.Context, oldName, newName, email, country, state string) (bool, error)
	DeleteUserWithId(ctx context.Context, userid uuid.UUID) (bool, error)
	DeleteUserWithName(ctx context.Context, name string) (bool, error)
	DeleteUserWithEmail(ctx context.Context, email string) (bool, error)
	DeleteUsersWithAttribute(ctx context.Context, attribute string) (bool, int, error)
	DeleteAllUsers(ctx context.Context) (bool, int, error)
	InsertPost(ctx context.Context, userid uuid.UUID, body, hashtag, title string) (bool, error)
	GetPostWithId(ctx context.Context, postid uuid.UUID) (*model.Post, error)
	GetPostsOfUser(ctx context.Context, userid uuid.UUID) ([]*model.Post, error)
	GetPostsOfHashtag(ctx context.Context, hashtag string) ([]*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	UpdatePostWithId(ctx context.Context, postid uuid.UUID, newtitle, newcontent, hashtag string) (bool, error)
	DeletePostWithId(ctx context.Context, postid uuid.UUID) (bool, error)
	DeletePostsOfUser(ctx context.Context, userid uuid.UUID) (bool, int, error)
	DeletePostsOfHashtag(ctx context.Context, hashtag string) (bool, int, error)
	DeleteAllPosts(ctx context.Context) (bool, int, error)
	GetCommentWithId(ctx context.Context, commentid uuid.UUID) (*model.Comment, error)
	GetCommentsOfPost(ctx context.Context, postid uuid.UUID) ([]*model.Comment, error)
	DeleteCommentWithId(ctx context.Context, commentid uuid.UUID) (bool, error)
	DeleteCommentsOfPost(ctx context.Context, postid uuid.UUID) (bool, int, error)
	InsertMessageIntoConversation(ctx context.Context, message *model.Message) (bool, error)
	UpdateMessageOfConversation(ctx context.Context, senderid uuid.UUID, recieverid uuid.UUID, newbody string) (bool, error)
	DeleteMessageOfConversation(ctx context.Context, senderid uuid.UUID, reciverid uuid.UUID, messageid int) (bool, error)
	GetAllMessagesOfConversation(ctx context.Context, senderid uuid.UUID, reciverid uuid.UUID) ([]*model.Message, error)
	GetAllMessagesInDB(ctx context.Context) ([]*model.Message, error)
	GetPasswordOfUserWithEmail(ctx context.Context, email string) (string, error)
}
