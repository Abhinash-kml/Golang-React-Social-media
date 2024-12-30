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

	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	if conn == nil {
		log.Fatal("Database connection failed.")
		return
	}

	d.primary = conn

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
		d.logger.Error("Error scanning row", zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUserWithName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}

	row := d.primary.QueryRowContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE name = $1;", name)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
		d.logger.Error("Error scanning row", zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUserWithEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	row := d.primary.QueryRowContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE email = $1;", email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &user.BanDuration, &user.AvatarUrl); err != nil {
		d.logger.Error("No rows in table", zap.Error(err))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUsersWithAttribute(ctx context.Context, attribute, value string) ([]*model.User, error) {
	users := []*model.User{}

	query := fmt.Sprintf("SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users WHERE %s = $1;", attribute)
	rows, err := d.primary.QueryContext(ctx, query, value)
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		var modifiedat, banduration sql.NullTime
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &modifiedat, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &banduration, &user.AvatarUrl); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		if modifiedat.Valid {
			user.Modified_at = modifiedat.Time
		}
		if banduration.Valid {
			user.BanDuration = banduration.Time
		}

		users = append(users, &user)
	}
	return users, nil
}

func (d *Postgres) UpdateUserWithId(ctx context.Context, userid uuid.UUID, name, email, country, state string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE userid = $5;", name, email, country, state, userid)
	if err != nil {
		d.logger.Error("Error updating row", zap.Error(err), zap.Any("uuid", userid))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) UpdateUserWithName(ctx context.Context, oldName, newName, email, country, state string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE name = $5;", newName, email, country, state, oldName)
	if err != nil {
		d.logger.Error("Error updating row", zap.Error(err), zap.Any("name", oldName))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users := []*model.User{}

	rows, err := d.primary.QueryContext(ctx, "SELECT userid, name, email, created_at, modified_at, last_login, country, state, city, ban_level, ban_duration, avatar_url FROM users;")
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banduration sql.NullTime
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Modified_at, &user.Lastlogin, &user.Country, &user.State, &user.City, &user.BanLevel, &banduration, &user.AvatarUrl); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		if banduration.Valid {
			user.BanDuration = banduration.Time
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
		d.logger.Error("Error scanning row", zap.Error(err))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned userid is not same as the supplied userid")
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithId(ctx context.Context, userid uuid.UUID) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM users WHERE userid = $1;", userid)
	if err != nil {
		d.logger.Error("Error deleting row", zap.Any("uuid", userid), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithName(ctx context.Context, name string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM users WHERE name = $1;", name)
	if err != nil {
		d.logger.Error("Error deleting row", zap.Any("name", name), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithEmail(ctx context.Context, email string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM users WHERE email = $1;", email)
	if err != nil {
		d.logger.Error("Error deleting row", zap.Any("email", email), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeleteUsersWithAttribute(ctx context.Context, attribute string) (bool, int, error) {
	return true, 0, nil
}

func (d *Postgres) DeleteAllUsers(ctx context.Context) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM users;")
	if err != nil {
		d.logger.Error("Error deleting all users", zap.Error(err))
		return false, 0, err
	}

	rowsEffected, _ := result.RowsAffected()
	if rowsEffected == 0 {
		return false, 0, nil
	}

	return true, int(rowsEffected), nil
}

func (d *Postgres) GetPasswordOfUserWithEmail(ctx context.Context, email string) (string, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT password FROM users WHERE email = $1;", email)
	var password string
	if err := row.Scan(&password); err != nil {
		d.logger.Error("Error scanning row", zap.Error(err))

		return "", nil
	}

	return password, nil
}

func (d *Postgres) InsertMediaWithId(ctx context.Context, postId uuid.UUID, url string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "INSERT INTO media(postid, url) VALUES($1, $2);", postId, url)
	if err != nil {
		d.logger.Error("Error inserting media", zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) UpdateMediaWithId(ctx context.Context, postId uuid.UUID, newUrl string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE media SET url = $1 WHERE postid = $2;", newUrl, postId)
	if err != nil {
		d.logger.Error("Cannot update media", zap.Any("uuid", postId), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *Postgres) DeleteMediaWithId(ctx context.Context, postId uuid.UUID) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM media WHERE postid = $1;", postId)
	if err != nil {
		d.logger.Error("Cannot update media", zap.Any("uuid", postId), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *Postgres) InsertMessageIntoConversation(ctx context.Context, message *model.Message) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "INSERT INTO messages(senderid, recieverid, body) VALUES($1, $2, $3);", message.SenderID, message.RecieverID, message.Body)
	if err != nil {
		d.logger.Error("Error inserting message into conversation", zap.Any("message", *message), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) GetAllMessagesOfConversation(ctx context.Context, senderId, receiverId uuid.UUID) ([]*model.Message, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT senderid, recieverid, body, status FROM messages WHERE senderid IN($1, $2) AND recieverid IN($1, $2);", senderId, receiverId)
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(&message.SenderID, &message.RecieverID, &message.Body, &message.Status); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (d *Postgres) GetAllMessagesInDB(ctx context.Context) ([]*model.Message, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM messages;")
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(&message.SenderID, &message.RecieverID, &message.Body, &message.Status, &message.Timestamp); err != nil {
			fmt.Println(err)
			d.logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (d *Postgres) UpdateMessageOfConversation(ctx context.Context, senderId, recieverId uuid.UUID, messageid int, newBody string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE messages SET body = $1 WHERE senderid = $2 AND recieverid = $3 AND id = $4;", newBody, senderId, recieverId, messageid)
	if err != nil {
		d.logger.Error("Error updating message of conversation", zap.Any("messageid", messageid), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeleteMessageOfConversation(ctx context.Context, senderId, receiverId uuid.UUID, messageId int) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM messages WHERE senderid = $1 AND receiverid = $2 AND id = $3;", senderId, receiverId, messageId)
	if err != nil {
		d.logger.Error("Error deleting message of conversation", zap.Any("messageid", messageId), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) InsertPost(ctx context.Context, userid uuid.UUID, title, body, mediaurl, hashtag string) (bool, error) {
	randomUuid := uuid.New()
	result, err := d.primary.ExecContext(ctx, "INSERT INTO posts(id, userid, title, body, media_url, hashtag) VALUES($1, $2, $3, $4, $5, $6);", randomUuid, userid, title, body, mediaurl, hashtag)
	if err != nil {
		d.logger.Error("Error inserting new post", zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) GetPostWithId(ctx context.Context, uuId uuid.UUID) (*model.Post, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT * FROM posts WHERE id = $1;", uuId)
	post := &model.Post{}
	var modifiedat sql.NullTime
	if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &modifiedat); err != nil {
		d.logger.Error("Error scanning row", zap.Error(err))

		return nil, err
	}

	if modifiedat.Valid {
		post.Modified_at = modifiedat.Time
	}

	return post, nil
}

func (d *Postgres) GetPostsOfUser(ctx context.Context, userid uuid.UUID) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts WHERE userid = $1;", userid)
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	posts := []*model.Post{}

	for rows.Next() {
		var mediaurl sql.NullString
		var hashtag sql.NullString
		var createdat sql.NullTime
		var modifiedat sql.NullTime
		var post model.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &mediaurl, &hashtag, &createdat, &modifiedat); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		if mediaurl.Valid {
			post.MediaUrl = mediaurl.String
		}
		if hashtag.Valid {
			post.Hashtag = hashtag.String
		}
		if createdat.Valid {
			post.Created_at = createdat.Time
		}
		if modifiedat.Valid {
			post.Modified_at = modifiedat.Time
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (d *Postgres) GetPostsOfHashtag(ctx context.Context, hashtag string) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts WHERE hashtag = $1;", hashtag)
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	posts := make([]*model.Post, 1)
	post := &model.Post{}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &post.Modified_at); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Postgres) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM posts ORDER BY title ASC;")
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	posts := []*model.Post{}

	for rows.Next() {
		var post model.Post
		var modifiedat sql.NullTime
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Body, &post.Likes, &post.Comments, &post.MediaUrl, &post.Hashtag, &post.Created_at, &modifiedat); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}
		if modifiedat.Valid {
			post.Modified_at = modifiedat.Time
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (d *Postgres) UpdatePostWithId(ctx context.Context, postId uuid.UUID, newTitle, newContent, hashtag string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE posts SET title = $1, body = $2, hashtag = $3, modified_at = $4 WHERE id = $5;", newTitle, newContent, hashtag, time.Now(), postId)
	if err != nil {
		d.logger.Error("Error updating post", zap.Any("postid", postId), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeletePostWithId(ctx context.Context, uuId uuid.UUID) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts WHERE id = $1;", uuId)
	if err != nil {
		d.logger.Error("Error deleting post", zap.Any("postid", uuId), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeletePostsOfUser(ctx context.Context, userId uuid.UUID) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts WHERE userid = $1;", userId)
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set", zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	if rowsEffected == 0 {
		return false, 0, nil
	}

	return true, int(rowsEffected), nil
}
func (d *Postgres) DeletePostsOfHashtag(ctx context.Context, hashtag string) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts WHERE hashtag = $1;", hashtag)
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set", zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	if rowsEffected == 0 {
		return false, 0, nil
	}

	return true, int(rowsEffected), nil
}
func (d *Postgres) DeleteAllPosts(ctx context.Context) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM posts;")
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Error("No rows in result set", zap.Error(err))

			return false, 0, err
		}
	}

	rowsEffected, _ := result.RowsAffected()
	if rowsEffected == 0 {
		return false, 0, nil
	}

	return true, int(rowsEffected), nil
}

func (d *Postgres) GetCommentWithId(ctx context.Context, commentid uuid.UUID) (*model.Comment, error) {
	row := d.primary.QueryRowContext(ctx, "SELECT * FROM comments WHERE id = $1;", commentid)
	comment := &model.Comment{}
	var modifiedat sql.NullTime
	if err := row.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.Created_at, &modifiedat); err != nil {
		d.logger.Error("Error scanning row", zap.Error(err))

		return nil, err
	}

	if modifiedat.Valid {
		comment.Modified_at = modifiedat.Time
	}

	return comment, nil
}

func (d *Postgres) GetCommentsOfPost(ctx context.Context, postid uuid.UUID) ([]*model.Comment, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM comments WHERE postid = $1;", postid)
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}

	for rows.Next() {
		var modifiedat sql.NullTime
		var comment model.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.Created_at, &modifiedat); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		if modifiedat.Valid {
			comment.Modified_at = modifiedat.Time
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (d *Postgres) DeleteCommentWithId(ctx context.Context, commentid uuid.UUID) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM comments WHERE id = $1;", commentid)
	if err != nil {
		d.logger.Error("Error updating message of conversation", zap.Any("commentid", commentid), zap.Error(err))
		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) DeleteCommentsOfPost(ctx context.Context, postid uuid.UUID) (bool, int, error) {
	result, err := d.primary.ExecContext(ctx, "DELETE FROM comments WHERE postid = $1;")
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Error("No rows to delete in result set", zap.Error(err))

			return false, 0, err
		}
	}

	deletedRows, _ := result.RowsAffected()
	if deletedRows == 0 {
		return false, 0, nil
	}

	return true, int(deletedRows), nil
}

func (d *Postgres) GetAllComments(ctx context.Context) ([]*model.Comment, error) {
	rows, err := d.primary.QueryContext(ctx, "SELECT * FROM comments;")
	if err != nil {
		d.logger.Error("Failed to execute sql query", zap.Error(err))

		rows.Close()
		return nil, err
	}
	defer rows.Close()
	comments := []*model.Comment{}

	for rows.Next() {
		var comment model.Comment
		var modifiedat sql.NullTime
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &comment.Created_at, &modifiedat); err != nil {
			d.logger.Error("Error scanning row", zap.Error(err))

			return nil, err
		}

		if modifiedat.Valid {
			comment.Modified_at = modifiedat.Time
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (d *Postgres) AddCommentToPostId(ctx context.Context, postid uuid.UUID, body string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "INSERT INTO comments(id, postid, body) VALUES($1, $2, $3);", uuid.New(), postid, body)
	if err != nil {
		d.logger.Error("Error inserting new comment", zap.Error(err))

		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Postgres) UpdateCommentWithId(ctx context.Context, commentid uuid.UUID, newBody string) (bool, error) {
	result, err := d.primary.ExecContext(ctx, "UPDATE comments SET body = $1, modified_at = $2 WHERE id = $3;", newBody, time.Now(), commentid)
	if err != nil {
		d.logger.Error("Error updating comment with id", zap.Error(err))

		return false, err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return false, nil
	}

	return true, nil
}
