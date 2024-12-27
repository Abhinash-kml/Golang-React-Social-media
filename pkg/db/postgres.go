package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	model "github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	/*"github.com/golang-migrate/migrate/v4"*/

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Postgres struct {
	primary *sql.DB
	replica *sql.DB
	logger  *zap.Logger
}

func (d *Postgres) Connect() {
	fmt.Println("Establising connection to postgres...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	d.primary, err = sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	if d.primary == nil {
		log.Fatal("Database connection failed.")
		return
	}

	fmt.Println("Connected to postgres.")

	REPLICA_DB_URL := os.Getenv("REPLICA_DATABASE_URL")
	if REPLICA_DB_URL == "" {
		d.replica = nil
		fmt.Println("The replication database is not setup, no replication will take place.")
	}

	d.CreateTables()
	d.logger, _ = zap.NewProduction()
}

func (d *Postgres) Disconnect() {
	if d.primary == nil {
		fmt.Println("Trying to close a connection which is already nil")
		return
	}

	d.primary.Close()
	d.logger.Sync()
	fmt.Println("Disconnected from postgres.")
}

func (d *Postgres) CreateTables() {
	// m, err := migrate.New(
	// 	"file://pkg/db/migrations",
	// 	"postgresql://postgres:Abx305@localhost:5432/SocialMedia?sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := m.Up(); err != nil {
	// 	log.Fatal(err)
	// }
}

func (d *Postgres) GetUserWithId(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	row := d.primary.QueryRowContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE userid = $1;", id)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetUserWithId"),
			zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUserWithName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}

	row := d.primary.QueryRowContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE name = $1;", name)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetUserWithName"),
			zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUserWithEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	row := d.primary.QueryRowContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE email = $1;", email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetUserWithEmail"),
			zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUsersWithAttribute(ctx context.Context, attribute, value string) ([]*model.User, error) {
	users := make([]*model.User, 1)
	var user model.User

	rows, err := d.primary.QueryContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE $1 = $2;", attribute, value)
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in table",
				zap.String("function", "GetUsersWithAttribute"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetUsersWithAttribute"),
				zap.Error(err))

			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

/*
* Todo: add this function into interface
 */
func (d *Postgres) UpdateUserWithId(ctx context.Context, userid uuid.UUID, name, email, country, state string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE userid = $5 RETURNING userid;", name, email, country, state, userid)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "UpdateUserWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != userid {
		return false, errors.New("returned userid is not same as supplied userid in UpdateUserWithId()")
	}

	return true, nil
}

func (d *Postgres) UpdateUserWithName(ctx context.Context, oldName, newName, email, country, state string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE name = $5 RETURNING userid;", newName, email, country, state, oldName)
	var returnedName string
	if err := row.Scan(&returnedName); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "UpdateUserWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedName != newName {
		return false, errors.New("returned newname is not same as supplied newname in UpdateUserWithId()")
	}

	return true, nil
}

func (d *Postgres) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users := make([]*model.User, 1)
	var user model.User

	rows, err := d.primary.QueryContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users;")
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in table",
				zap.String("function", "GetUsersWithAttribute"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetUsersWithAttribute"),
				zap.Error(err))

			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

func (d *Postgres) InsertUser(ctx context.Context, fullname, email, password, dob, country, state, city, avatar_url string) (bool, error) {
	uuId, _ := uuid.NewRandom()
	row := d.primary.QueryRowContext(ctx, `INSERT INTO users(userid, name, email, password, dob, country, state, city, avatar_url)
										   VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
										   ON CONFLICT(name)
										   DO NOTHING
										   RETURNING userid`,
		uuId, fullname, email, password, dob, country, state, city, avatar_url)

	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "InsertUser"),
			zap.Error(err))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned userid is not same as the supplied userid")
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithId(ctx context.Context, userid uuid.UUID) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM users WHERE userid = $1 RETURNING userid;", userid)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "DeleteUserWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != userid {
		return false, errors.New("returned userid is not same as supplied userid in DeleteUserWithId()")
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithName(ctx context.Context, name string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM users WHERE name = $1 RETURNING name;", name)
	var returnedName string
	if err := row.Scan(&returnedName); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "DeleteUserWithName"),
			zap.Error(err))

		return false, err
	}

	if returnedName != name {
		return false, errors.New("returned userid is not same as supplied userid in DeleteUserWithName()")
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithEmail(ctx context.Context, email string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM users WHERE email = $1 RETURNING email;", email)
	var returnedMail string
	if err := row.Scan(&returnedMail); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "DeleteUserWithEmail"),
			zap.Error(err))

		return false, err
	}

	if returnedMail != email {
		return false, errors.New("returned userid is not same as supplied userid in DeleteUserWithEmail()")
	}

	return true, nil
}

func (d *Postgres) DeleteUsersWithAttribute(ctx context.Context, attribute string) (bool, int, error) {
	return true, 0, nil
}

func (d *Postgres) DeleteAllUsers(ctx context.Context) (bool, int, error) {
	return true, 0, nil
}

func (d *Postgres) GetPasswordOfUserWithEmail(ctx context.Context, email string) (string, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT password FROM users WHERE email = $1;", email)
	var returendPassword string
	if err := row.Scan(&returendPassword); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("Function", "GetPasswordOfUserWithEmail"),
			zap.Error(err))

		return "", nil
	}

	return returendPassword, nil
}

func (d *Postgres) InsertMediaWithId(ctx context.Context, postId uuid.UUID, url string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "INSERT INTO media(postid, url) VALUES($1, $2) returning postid;", postId, url)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "InsertMediaWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in InsertMedia()")
	}

	return true, nil
}

func (d *Postgres) UpdateMediaWithId(ctx context.Context, postId uuid.UUID, newUrl string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "UPDATE media SET url = $1 WHERE postid = $2 returning postid;", newUrl, postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "UpdateMediaWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in UpdateMedia()")
	}

	return true, nil
}

func (d *Postgres) DeleteMediaWithId(ctx context.Context, postId uuid.UUID) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM media WHERE postid = $1 returning postid;", postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "DeleteMediaWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned postid is different from the supplied post id in DeleteMedia()")
	}

	return true, nil
}

func (d *Postgres) InsertMessageIntoConversation(ctx context.Context, message *model.Message) (bool, error) {
	var sender_id, reciever_id uuid.UUID
	row := d.primary.QueryRowContext(ctx, "INSERT INTO messages(senderid, recieverid, body) VALUES($1, $2, $3) RETURNING senderid, recieverid;", message.SenderID, message.RecieverID, message.Body)
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "InsertMessageIntoConversation"),
			zap.Error(err))

		return false, err
	}

	if sender_id == message.SenderID && reciever_id == message.RecieverID {
		return true, nil
	}

	return false, errors.New("something wrong happened in InsertMessageInDb()")
}

func (d *Postgres) GetAllMessagesOfConversation(ctx context.Context, senderId, receiverId uuid.UUID) ([]*model.Message, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT senderid, recieverid, body, status FROM messages WHERE senderid IN($1, $2) AND recieverid IN($1, $2);", senderId, receiverId)
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "GetAllMessagesOfConversation"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(&message.SenderID, &message.RecieverID, &message.Body, &message.Status); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetAllMessagesOfConversation"),
				zap.Error(err))

			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (d *Postgres) GetAllMessagesInDB(ctx context.Context) ([]*model.Message, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM messages;")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
	}
	defer rows.Close()

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(&message.SenderID, &message.RecieverID, &message.Body, &message.Status, &message.Timestamp); err != nil {
			fmt.Println(err)
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (d *Postgres) UpdateMessageOfConversation(ctx context.Context, senderId, recieverId uuid.UUID, newBody string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "UPDATE messages SET body = $1 WHERE senderid = $2 AND recieverid = $3 RETURNING senderid, receiverid;", newBody, senderId, recieverId)
	var sender_id, reciever_id uuid.UUID
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "UpdateMessageOfConversation"),
			zap.Error(err))

		return false, err
	}

	if sender_id != senderId || reciever_id != recieverId {
		return false, errors.New("something wrong happened in UpdateMessageWithId()")
	}

	return true, nil
}

func (d *Postgres) DeleteMessageOfConversation(ctx context.Context, senderId, receiverId uuid.UUID, messageId int) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM messages WHERE senderid = $1 AND receiverid = $2 AND id = $3 RETURNING senderid, receiverid;", senderId, receiverId, messageId)
	var sender_id, receiver_id uuid.UUID
	if err := row.Scan(&sender_id, &receiver_id); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", ""),
			zap.Error(err))

		return false, err
	}

	if sender_id != senderId || receiver_id != receiverId {
		return false, errors.New("something wrong happended in DeleteMessageOfConversation()")
	}

	return true, nil
}

func (d *Postgres) InsertPost(ctx context.Context, userid uuid.UUID, body, hashtag, title string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "INSERT INTO posts(userid, title, body, hashtag) VALUES($1, $2, $3, $4) returning userid;", userid, title, body, hashtag)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "InsertPost"),
			zap.Error(err))

		return false, err
	}

	if returnedId != userid {
		return false, errors.New("the returned uuid from insertion is not same as supplied uuid in InsertPost()")
	}

	return true, nil
}

func (d *Postgres) GetPostWithId(ctx context.Context, uuId uuid.UUID) (*model.Post, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT * FROM posts WHERE id = $1;", uuId)
	post := &model.Post{}

	if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &post.Modified_at); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetPostWithId"),
			zap.Error(err))

		return nil, err
	}

	return post, nil
}

func (d *Postgres) GetPostsOfUser(ctx context.Context, userid uuid.UUID) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts WHERE userid = $1;", userid)
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "GetPostsOfUser"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	posts := make([]*model.Post, 1)
	post := &model.Post{}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &post.Modified_at); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetPostsOfUser"),
				zap.Error(err))

			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Postgres) GetPostsOfHashtag(ctx context.Context, hashtag string) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts WHERE hashtag = $1;", hashtag)
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "GetPostsOfUser"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	posts := make([]*model.Post, 1)
	post := &model.Post{}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &post.Modified_at); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetPostsOfHashtag"),
				zap.Error(err))

			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Postgres) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts ORDER BY title ASC;")
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "GetPostsOfUser"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	posts := make([]*model.Post, 1)
	post := &model.Post{}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &post.Modified_at); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetAllPosts"),
				zap.Error(err))

			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Postgres) UpdatePostWithId(ctx context.Context, postId uuid.UUID, newTitle, newContent, hashtag string) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "UPDATE posts SET title = $1, body = $2, hastag = $3, modified_at = $4 WHERE id = $5 RETURNING id;", newTitle, newContent, hashtag, time.Now(), postId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "UpdatePostWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != postId {
		return false, errors.New("the returned uuid from update in not same as supplied uuid in UpdatePost()")
	}

	return true, nil
}

func (d *Postgres) DeletePostWithId(ctx context.Context, uuId uuid.UUID) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM posts WHERE id = $1 returning id;", uuId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "DeletePostWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from delete is not same as supplied uuid in DeletePostWithId()")
	}

	return true, nil
}

func (d *Postgres) DeletePostsOfUser(ctx context.Context, userId uuid.UUID) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts WHERE userid = $1;", userId)
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "DeletePostsOfUser"),
				zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	return true, int(rowsEffected), nil
}
func (d *Postgres) DeletePostsOfHashtag(ctx context.Context, hashtag string) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts WHERE hashtag = $1;", hashtag)
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "DeletePostsOfHashtag"),
				zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	return true, int(rowsEffected), nil
}
func (d *Postgres) DeleteAllPosts(ctx context.Context) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts;")
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "DeleteAllPosts"),
				zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	return true, int(rowsEffected), nil
}

func (d *Postgres) GetCommentWithId(ctx context.Context, commentid uuid.UUID) (*model.Comment, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT * FROM comments WHERE id = $1;", commentid)
	comment := &model.Comment{}
	if err := row.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.Created_at, &comment.Modified_at); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetCommentWithId"),
			zap.Error(err))

		return nil, err
	}

	return comment, nil
}

func (d *Postgres) GetCommentsOfPost(ctx context.Context, postid uuid.UUID) ([]*model.Comment, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM comments WHERE postid = $1;", postid)
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows in result set",
				zap.String("function", "GetCommentsOfPost"),
				zap.Error(err))

			rows.Close()
			return nil, err
		}
	}
	defer rows.Close()

	comments := make([]*model.Comment, 1)
	comment := &model.Comment{}

	for rows.Next() {
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.Created_at, &comment.Modified_at); err != nil {
			d.logger.Error("Error scanning row",
				zap.String("function", "GetCommentsOfPost"),
				zap.Error(err))

			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (d *Postgres) DeleteCommentWithId(ctx context.Context, commentid uuid.UUID) (bool, error) {
	row := d.primary.QueryRowContext(ctx, "DELETE FROM comments WHERE id = $1 RETURNING id;", commentid)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		d.logger.Error("Error scanning row",
			zap.String("function", "GetCommentWithId"),
			zap.Error(err))

		return false, err
	}

	if returnedId != commentid {
		return false, errors.New("the returned commentid doesnt match with supplied in DeleteCommentWithId()")
	}

	return true, nil
}

func (d *Postgres) DeleteCommentsOfPost(ctx context.Context, postid uuid.UUID) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM comments WHERE postid = $1;")
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows to delete in result set",
				zap.String("function", "DeleteCommentsOfPost"),
				zap.Error(err))

			return false, 0, err
		}
	}

	deletedRows, _ := result.RowsAffected()
	return true, int(deletedRows), nil

}
