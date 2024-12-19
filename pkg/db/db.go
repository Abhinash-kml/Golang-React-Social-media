package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Connection *sql.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		log.Fatal("Returned db is still invalid")
	} else {
		fmt.Println("Connected to postgres db")
	}

	Connection = db

	CreateTables()
}

func Disconnect() {
	if Connection == nil {
		log.Fatal("Trying to close a databse connection which is already closed or invalid")
	}

	Connection.Close()
	Connection = nil
}

func CreateTables() {

}
