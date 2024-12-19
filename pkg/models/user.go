package models

import (
	"fmt"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	Id          uuid.UUID `json:"id"`
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
			zap.String("function", "GetUserWithId"))
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
			zap.String("funcion", "GetUserWith"))
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
