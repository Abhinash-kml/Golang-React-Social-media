package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"Id"`
	Fullname    string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Dob         string    `json:"dob"`
	Created_at  string    `json:"created_at"`
	Modified_at string    `json:"modified_at"`
	Lastlogin   string    `json:"lastlogin"`
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

func GetUserWithID(DB *sql.DB, id uuid.UUID) (*User, error) {
	if DB == nil {
		return nil, errors.New("tried calling GetUserWithID() with invalid pointer to Database")
	}

	user := &User{}

	if err := DB.QueryRow("SELECT * FROM users WHERE id = $1;", id).Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsersWithIDs(DB *sql.DB, ids []uuid.UUID, limit int) ([]*User, error) {
	if DB == nil {
		return nil, errors.New("tried calling GetUsersWithId() with invalid pointer to Database")
	}

	var users []*User

	for _, val := range ids {
		user, err := GetUserWithID(DB, val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserWith(DB *sql.DB, with string, condition string, value string) (*User, error) {
	if DB == nil {
		return nil, errors.New("tried calling GetUserWith() with invalid pointer to Database")
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE %s %s %s", with, condition, value)

	user := &User{}

	if err := DB.QueryRow(query).Scan(user.Id, user.Fullname, user.Email, user.Password, user.Dob, user.Created_at, user.Modified_at, user.Lastlogin); err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser() {
	// logger, _ := zap.NewProduction()
}

func DeleteUser() {

}

func InsertUser() {

}
