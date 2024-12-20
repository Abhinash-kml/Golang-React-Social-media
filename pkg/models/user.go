package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	Id          uuid.UUID `json:"uuid"`
	Fullname    string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Dob         string    `json:"dob"`
	Created_at  string    `json:"created_at"`
	Modified_at string    `json:"modified_at"`
	Lastlogin   string    `json:"lastlogin"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	City        string    `json:"city"`
}

func NewUser(name, email, password, dob, created_at, modified_at, last_login string) *User {
	return &User{
		Fullname:    name,
		Email:       email,
		Password:    password,
		Dob:         dob,
		Created_at:  created_at,
		Modified_at: modified_at,
		Lastlogin:   last_login,
	}
}

func GetUserWithID(logger *zap.Logger, id uuid.UUID) (*User, error) {
	user := &User{}

	row := db.Connection.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	if err := row.Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		logger.Error("Error scanning row",
			zap.String("function", "GetUserWithId"),
			zap.String("Error", err.Error()))

		return nil, err
	}

	return user, nil
}

func GetUsersWithIDs(logger *zap.Logger, ids []uuid.UUID, limit int) ([]*User, error) {
	var users []*User

	for _, val := range ids {
		user, err := GetUserWithID(logger, val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserWith(logger *zap.Logger, with string, condition string, value string) (*User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE %s %s %s", with, condition, value)

	user := &User{}

	row := db.Connection.QueryRow(query)
	if err := row.Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		logger.Error("Error scanning row",
			zap.String("Funcion", "GetUserWith"),
			zap.String("Error", err.Error()))

		return nil, err
	}

	return user, nil
}

func UpdateUserWithId(logger *zap.Logger, userid uuid.UUID, name, email, country, state string) (bool, error) {
	row := db.Connection.QueryRow("UPDATE users SET name = $1, email = $2, country = $3, state = $4 WHERE userid = $5 RETURNING userid;", name, email, country, state, userid)
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

func DeleteUserWithId(logger *zap.Logger, userid uuid.UUID) (bool, error) {
	row := db.Connection.QueryRow("DELETE FROM users WHERE userid = $1 RETURNING userid;", userid)
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

func InsertUser(logger *zap.Logger, fullname, email, password, dob, country, state, city string) (bool, error) {
	uuId, _ := uuid.NewRandom()
	createdAt := time.Now().String()
	modifiedAt := time.Now().String()
	lastLogin := time.Now().String()

	row := db.Connection.QueryRow(`INSERT INTO users(uuid, fullname, email, password, dob, country, state, city, created_at, modified_at, last_login)
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING user_id`,
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

func GetPasswordOfUserWithEmail(logger *zap.Logger, email string) (string, error) {
	row := db.Connection.QueryRow("SELECT password FROM users WHERE email = $1;", email)
	var returnedEmail string
	if err := row.Scan(&returnedEmail); err != nil {
		logger.Error("Error scanning row",
			zap.String("Function", "GetPasswordOfUserWithEmail"),
			zap.String("Error", err.Error()))

		return "", nil
	}

	return returnedEmail, nil
}

func InsertNewUserIntoDatabase(logger *zap.Logger, email, password string) error {
	ok, err := InsertUser(logger, "", email, password, "", "India", "Random", "Random")
	if err != nil {
		return err
	}

	return nil
}
