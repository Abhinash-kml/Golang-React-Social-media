package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connection *sql.DB

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
	}

	connection = db

	CreateTables()
}

func Disconnect() {
	if connection == nil {
		log.Fatal("Trying to close a databse connection which is already closed or invalid")
	}

	connection.Close()
	connection = nil
}

func CreateTables() {

}
