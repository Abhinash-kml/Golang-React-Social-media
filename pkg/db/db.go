package db

type Database interface {
	Connect()
	Disconnect()
	CreateTables()
}
