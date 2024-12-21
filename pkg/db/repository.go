package db

type Repository interface {
	Connect()
	Disconnect()
}
