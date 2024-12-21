package db

import (
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
)

type Postgres struct {
	conn *sql.DB
}

func (d *Postgres) Connect() {
	fmt.Println("Establising connection to postgres...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	d.conn, err = sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	if d.conn == nil {
		log.Fatal("Database connection failed.")
		return
	}

	fmt.Println("Connected to postgres.")
	d.CreateTables()
}

func (d *Postgres) Disconnect() {
	if d.conn == nil {
		fmt.Println("Trying to close a connection which is already nil")
		return
	}

	d.conn.Close()
	fmt.Println("Disconnected from postgres.")
}

func (d *Postgres) CreateTables() {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS users(
		id SERIAL,
		userid uuid UNIQUE NOT NULL,
		name VARCHAR UNIQUE NOT NULL,
		email VARCHAR,
		password VARCHAR,
		dob	VARCHAR,
		country VARCHAR,
		state VARCHAR,
		city VARCHAR,
		created_at VARCHAR,
		modified_at VARCHAR,
		last_login VARCHAR)`)

	// userid, name, email, password, dob, country, state, city, created_at, modified_at, last_login

	if err != nil {
		fmt.Println(err)
	}
}

func (d *Postgres) GetUserWithID(logger *zap.Logger, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	row := d.conn.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	if err := row.Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "GetUserWithId"),
			zap.String("Error", err.Error()))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) GetUsersWithIDs(logger *zap.Logger, ids []uuid.UUID, limit int) ([]*model.User, error) {
	var users []*model.User

	for _, val := range ids {
		user, err := d.GetUserWithID(logger, val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (d *Postgres) GetUserWith(logger *zap.Logger, with string, condition string, value string) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE %s %s %s", with, condition, value)

	user := &model.User{}

	row := d.conn.QueryRow(query)
	if err := row.Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		logger.Error("Error scanning row",
			zap.String("Funcion", "GetUserWith"),
			zap.String("Error", err.Error()))

		return nil, err
	}

	return user, nil
}

func (d *Postgres) UpdateUserWithId(logger *zap.Logger, userid uuid.UUID, name, email, country, state string) (bool, error) {
	row := d.conn.QueryRow("UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE userid = $5 RETURNING userid;", name, email, country, state, userid)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("Function", "UpdateUserWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != userid {
		return false, errors.New("returned userid is not same as supplied userid in UpdateUserWithId()")
	}

	return true, nil
}

func (d *Postgres) DeleteUserWithId(logger *zap.Logger, userid uuid.UUID) (bool, error) {
	row := d.conn.QueryRow("DELETE FROM users WHERE userid = $1 RETURNING userid;", userid)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("Function", "DeleteUserWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != userid {
		return false, errors.New("returned userid is not same as supplied userid in DeleteUserWithId()")
	}

	return true, nil
}

func (d *Postgres) InsertUser(logger *zap.Logger, fullname, email, password, dob, country, state, city string) (bool, error) {
	uuId, _ := uuid.NewRandom()
	createdAt := time.Now().String()
	modifiedAt := time.Now().String()
	lastLogin := time.Now().String()

	row := d.conn.QueryRow(`INSERT INTO users(userid, name, email, password, dob, country, state, city, created_at, modified_at, last_login)
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
								ON CONFLICT(name)
								DO NOTHING
								RETURNING userid`,
		uuId, fullname, email, password, dob, country, state, city, createdAt, modifiedAt, lastLogin)

	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("Function", "InsertUser"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned userid is not same as the supplied userid")
	}

	return true, nil
}

func (d *Postgres) GetPasswordOfUserWithEmail(logger *zap.Logger, email string) (string, error) {
	row := d.conn.QueryRow("SELECT password FROM users WHERE email = $1;", email)
	var returnedEmail string
	if err := row.Scan(&returnedEmail); err != nil {
		logger.Error("Error scanning row",
			zap.String("Function", "GetPasswordOfUserWithEmail"),
			zap.String("Error", err.Error()))

		return "", nil
	}

	return returnedEmail, nil
}

func (d *Postgres) InsertNewUserIntoDatabase(logger *zap.Logger, name, email, password string) error {
	_, err := d.InsertUser(logger, name, email, password, "", "India", "Random", "Random")
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) InsertMediaWithId(logger *zap.Logger, postId uuid.UUID, url string) (bool, error) {
	row := d.conn.QueryRow("INSERT INTO media(postid, url) VALUES($1, $2) returning postid;", postId, url)
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

func (d *Postgres) UpdateMediaWithId(logger *zap.Logger, postId uuid.UUID, newUrl string) (bool, error) {
	row := d.conn.QueryRow("UPDATE media SET url = $1 WHERE postid = $2 returning postid;", newUrl, postId)
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

func (d *Postgres) DeleteMediaWithId(logger *zap.Logger, postId uuid.UUID) (bool, error) {
	row := d.conn.QueryRow("DELETE FROM media WHERE postid = $1 returning postid;", postId)
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

func (d *Postgres) GetMediaWithId(logger *zap.Logger, postId uuid.UUID) *model.Media {
	row := d.conn.QueryRow("SELECT * FROM media WHERE postid = $1;", postId)
	media := &model.Media{}
	if err := row.Scan(media.PostID, media.Url); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "GetMediaWithId"),
			zap.String("Error", err.Error()))

		return nil
	}

	return media
}

func (d *Postgres) InsertMessageInDb(logger *zap.Logger, message *model.Message) (bool, error) {
	var sender_id, reciever_id uuid.UUID
	row := d.conn.QueryRow("INSERT INTO messages(sender_id, reciever_id, content) VALUES($1, $2, $3) RETURNING sender_id, reciever_id;")
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "InsertMessageInDb"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id == message.SenderID && reciever_id == message.RecieverID {
		return true, nil
	}

	return false, errors.New("something wrong happened in InsertMessageInDb()")
}

func (d *Postgres) GetAllMessagesOfSenderAndReciever(logger *zap.Logger, sender_id, receiver_id uuid.UUID) []*model.Message {
	rows, err := d.conn.Query("SELECT sender_id, reciever_id, content FROM messages WHERE sender_id IN($1, $2) AND reciever_id IN($1, $2);", sender_id, receiver_id)
	if err != nil {
		logger.Error("Error scanning rows",
			zap.String("function", "GetAllMessagesOfSenderAndReciever"),
			zap.String("Error", err.Error()))
		return nil
	}

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(message.SenderID, message.RecieverID, message.Content); err != nil {
			logger.Error("Error scanning row",
				zap.String("function", "GetAllMessagesOfSenderAndReciever"),
				zap.String("Error", err.Error()))
			return nil
		}

		messages = append(messages, message)
	}

	return messages
}

func (d *Postgres) GetAllMessagesInDB() ([]*model.Message, error) {
	rows, err := d.conn.Query("SELECT * FROM messages;")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
	}

	var (
		messages []*model.Message
		message  *model.Message = &model.Message{}
	)

	for rows.Next() {
		if err := rows.Scan(message.SenderID, message.RecieverID, message.Content); err != nil {
			fmt.Println(err)
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (d *Postgres) UpdateMessageWithId(logger *zap.Logger, senderId, recieverId uuid.UUID, newContent string) (bool, error) {
	row := d.conn.QueryRow("UPDATE messages SET content = $1 WHERE sender_id = $2 AND reciever_id = $3 RETURNING sender_id, receiver_id;", newContent, senderId, recieverId)
	var sender_id, reciever_id uuid.UUID
	if err := row.Scan(&sender_id, &reciever_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "UpdateMessageWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id != senderId || reciever_id != recieverId {
		return false, errors.New("something wrong happened in UpdateMessageWithId()")
	}

	return true, nil
}

func (d *Postgres) DeleteMessage(logger *zap.Logger, senderId, receiverId uuid.UUID, content string) (bool, error) {
	row := d.conn.QueryRow("DELETE FROM messages WHERE sender_id = $1 AND receiver_id = $2 AND content = $3 RETURNING sender_id, receiver_id;", senderId, receiverId, content)
	var sender_id, receiver_id uuid.UUID
	if err := row.Scan(&sender_id, &receiver_id); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", ""),
			zap.String("Error", err.Error()))

		return false, err
	}

	if sender_id != senderId || receiver_id != receiverId {
		return false, errors.New("something wrong happended in DeleteMessage()")
	}

	return true, nil
}

func (d *Postgres) InsertPost(logger *zap.Logger, uuId uuid.UUID, content, hashtag string) (bool, error) {
	row := d.conn.QueryRow("INSERT INTO posts(userid, content, hashtag) VALUES($1, $2, $3) returning userid;", uuId, content, hashtag)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "InsertPost"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from insertion is not same as supplied uuid in InsertPost()")
	}

	return true, nil
}

func (d *Postgres) UpdatePostWithId(logger *zap.Logger, uuId uuid.UUID, newContent, hashtag string) (bool, error) {
	row := d.conn.QueryRow("UPDATE posts SET content = $1, hastag = $2 WHERE userid = $3 RETURNING userid;", newContent, hashtag, uuId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "UpdatePostWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from update in not same as supplied uuid in UpdatePost()")
	}

	return true, nil
}

func (d *Postgres) DeletePostWithId(logger *zap.Logger, uuId uuid.UUID) (bool, error) {
	row := d.conn.QueryRow("DELETE FROM posts WHERE userid = $1 returning userid;", uuId)
	var returnedId uuid.UUID
	if err := row.Scan(&returnedId); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "DeletePostWithId"),
			zap.String("Error", err.Error()))

		return false, err
	}

	if returnedId != uuId {
		return false, errors.New("the returned uuid from delete is not same as supplied uuid in DeletePostWithId()")
	}

	return true, nil
}

func (d *Postgres) GetPostWithId(logger *zap.Logger, uuId uuid.UUID) *model.Post {
	row := d.conn.QueryRow("SELECT * FROM posts WHERE userid = $1;", uuId)
	post := &model.Post{}
	if err := row.Scan(post.UserID, post.Content, post.Hashtag); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "GetPostWithId"),
			zap.String("Error", err.Error()))

		return nil
	}

	return post
}
